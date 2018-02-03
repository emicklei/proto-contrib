package main

import (
	"flag"
	"github.com/emicklei/proto-contrib/pkg/proto2gql"
	"github.com/emicklei/proto-contrib/pkg/proto2gql/writers"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type (
	StringMap map[string]string
)

func (s *StringMap) String() string {
	pairs := make([]string, 0, len(*s))

	for key, value := range *s {
		pairs = append(pairs, key+"="+value)
	}

	return strings.Join(pairs, ",")
}

func (s *StringMap) Set(value string) error {
	parts := strings.Split(value, ",")

	for _, part := range parts {
		pair := strings.Split(part, "=")

		(*s)[strings.TrimSpace(pair[0])] = strings.TrimSpace(pair[1])
	}

	return nil
}

var (
	stdOut bool

	txtOut string

	goOut string

	jsOut string

	resolveImports StringMap

	packageAliases StringMap

	filter string

	filterN string

	noPrefix bool
)

func main() {
	resolveImports = make(StringMap)
	packageAliases = make(StringMap)

	flag.BoolVar(&stdOut, "std_out", false, "Writes transformed files to stdout")
	flag.StringVar(&txtOut, "txt_out", "", "Writes transformed files to .graphql file")
	flag.StringVar(&goOut, "go_out", "", "Writes transformed files to .go file")
	flag.StringVar(&jsOut, "js_out", "", "Writes transformed files to .js file")
	flag.Var(&resolveImports, "resolve_import", "Resolves given external packages")
	flag.Var(&packageAliases, "package_alias", "Renames packages using given aliases")
	flag.StringVar(&filter, "filter", "", "Regexp to filter out matched custom types")
	flag.StringVar(&filterN, "filterN", "", "Regexp to filter out not matched custom types")
	flag.BoolVar(&noPrefix, "no_prefix", false, "Disables package prefix for type names")

	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(0)
	}

	var transformer *proto2gql.Transformer
	ws := make([]io.Writer, 0, 5)

	if stdOut == true {
		ws = append(ws, os.Stdout)
	}

	if txtOut != "" {
		writer, err := createTextWriter(txtOut)

		if err != nil {
			gracefullyTerminate(err, ws)
		}

		ws = append(ws, writer)
	}

	if goOut != "" {
		writer, err := createGoWriter(goOut)

		if err != nil {
			gracefullyTerminate(err, ws)
		}

		ws = append(ws, writer)
	}

	if jsOut != "" {
		writer, err := createJsWriter(jsOut)

		if err != nil {
			gracefullyTerminate(err, ws)
		}

		ws = append(ws, writer)
	}

	if len(ws) == 0 {
		log.Println("output not defined")
		os.Exit(0)
	}

	transformer = proto2gql.NewTransformer(
		io.MultiWriter(ws...),
		withResolvingImports(resolveImports),
		withPackageAliases(packageAliases),
		withNoPrefix(noPrefix),
		withFilter(filter, filterN),
	)

	for _, filename := range flag.Args() {
		if err := readAndTransform(filename, transformer); err != nil {
			log.Fatalln("failed to transform file: " + err.Error())
		}
	}

	if err := saveWriters(ws); err != nil {
		log.Fatalln("failed to save output: " + err.Error())
	}
}

func withResolvingImports(imports StringMap) func(transformer *proto2gql.Transformer) {
	return func(t *proto2gql.Transformer) {
		for key, url := range imports {
			t.Import(key, url)
		}
	}
}

func withPackageAliases(aliases StringMap) func(transformer *proto2gql.Transformer) {
	return func(t *proto2gql.Transformer) {
		for pkg, alias := range aliases {
			t.SetPackageAlias(pkg, alias)
		}
	}
}

func withNoPrefix(noPrefix bool) func(transformer *proto2gql.Transformer) {
	return func(t *proto2gql.Transformer) {
		t.DisablePrefix(noPrefix)
	}
}

func withFilter(positive, negative string) func(transformer *proto2gql.Transformer) {
	return func(t *proto2gql.Transformer) {
		if positive == "" && negative == "" {
			return
		}

		chain := make([]func(typeName string) bool, 0, 2)

		if positive != "" {
			rPos, err := regexp.Compile(positive)

			if err != nil {
				log.Fatalln("invalid regular expression: " + err.Error())
			}

			chain = append(chain, func(typeName string) bool {
				// filter out matched types
				return rPos.Match([]byte(typeName)) == false
			})
		}

		if negative != "" {
			rNeg, err := regexp.Compile(negative)

			if err != nil {
				panic("invalid regular expression: " + err.Error())
			}

			chain = append(chain, func(typeName string) bool {
				// filter out not matched types
				return rNeg.Match([]byte(typeName))
			})
		}

		t.SetFilter(func(typeName string) bool {
			res := true

			for _, r := range chain {
				res = r(typeName)

				if res == false {
					break
				}
			}

			return res
		})
	}
}

func ensureExtension(filename, expectedExt string) string {
	ext := filepath.Ext(filename)

	if ext == "" {
		return filename + expectedExt
	}

	if ext != expectedExt {
		return strings.Replace(filename, ext, expectedExt, -1)
	}

	return filename
}

func createTextWriter(filename string) (io.Writer, error) {
	return writers.NewFileWriter(ensureExtension(filename, ".graphql"), "", "")
}

func createGoWriter(filename string) (io.Writer, error) {
	filename = ensureExtension(filename, ".go")
	abs, err := filepath.Abs(filename)

	if err != nil {
		return nil, err
	}

	name := strings.Replace(filepath.Base(abs), ".go", "", -1)

	openTag := "package " + filepath.Base(filepath.Dir(abs)) + "\n \n"
	openTag += "var " + strings.Title(name) + " = `\n"

	return writers.NewFileWriter(filename, openTag, "\n`")
}

func createJsWriter(filename string) (io.Writer, error) {
	openTag := "module.exports = `\n"

	return writers.NewFileWriter(ensureExtension(filename, ".js"), openTag, "\n`")
}

func saveWriters(ws []io.Writer) error {
	var err error

	for _, writer := range ws {
		fw, ok := writer.(*writers.FileWriter)

		if ok == true {
			if err == nil {
				err = fw.Save()

				if err != nil {
					break
				}
			} else {
				fw.Remove()
			}
		}
	}

	return err
}

func gracefullyTerminate(err error, ws []io.Writer) {
	for _, writer := range ws {
		fw, ok := writer.(*writers.FileWriter)

		if ok == true {
			fw.Remove()
		}
	}

	log.Fatalln("error occurred: " + err.Error())
}

func readAndTransform(filename string, transformer *proto2gql.Transformer) error {
	// open for read
	file, err := os.Open(filename)

	if err != nil {
		return err
	}

	defer file.Close()

	transformer.SetFilename(filename)
	if err := transformer.Transform(file); err != nil {
		return err
	}

	return nil
}
