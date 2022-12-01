package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"os"

	"github.com/emicklei/proto"
	"github.com/emicklei/proto-contrib/pkg/protofmt"
)

var (
	oInclude = flag.String("i", "", "proto file to include") // currently just one
	oProcess = flag.String("p", "", "proto file to process") // currently just one
	oOutput  = flag.String("o", "output.proto", "proto file to generate")
)

func main() {
	flag.Parse()
	b := newProtoBuilder()
	b.loadProto(*oInclude)
	toProcess := b.loadProto(*oProcess)
	b.processComposed()
	log.Println("writing", *oOutput)
	out, err := os.Create(*oOutput)
	check(err)
	defer out.Close()
	io.WriteString(out, formatted(toProcess))
}

func FieldOfMessage(m *proto.Message, fieldName string) proto.Visitee {
	for _, each := range m.Elements {
		// TODO other types
		if f, ok := each.(*proto.NormalField); ok {
			if f.Name == fieldName {
				return f
			}
		}
		if f, ok := each.(*proto.MapField); ok {
			if f.Name == fieldName {
				return f
			}
		}
	}
	return nil
}

func FieldNamesOfMessage(m *proto.Message) (list []string) {
	for _, each := range m.Elements {
		if f, ok := each.(*proto.NormalField); ok {
			// TODO other types
			list = append(list, f.Name)
		}
	}
	return
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func formatted(v proto.Visitee) string {
	b := new(bytes.Buffer)
	f := protofmt.NewFormatter(b, "  ") // 2 spaces
	v.Accept(f)
	return b.String()
}
