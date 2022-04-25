package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/emicklei/proto"
	"github.com/emicklei/proto-contrib/pkg/protofmt"
)

func main() {
	b := &protoBuilder{
		registry: make(map[string]*protoReference),
	}
	// quick
	var last *proto.Proto
	for _, each := range os.Args[1:] {
		last = b.loadProto(each)
	}
	b.processComposed()
	fmt.Println(formatted(last))
}

type protoBuilder struct {
	pkg string
	// key = package.message-name
	registry map[string]*protoReference
}

type protoReference struct {
	pkg          string
	message      *proto.Message
	composeSpecs []composeSpec
}

func (b *protoBuilder) loadProto(absFilename string) *proto.Proto {
	reader, err := os.Open(absFilename)
	check(err)
	defer reader.Close()
	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	check(err)
	proto.Walk(definition,
		proto.WithPackage(b.handlePackage),
		proto.WithMessage(b.handleMessage))
	return definition
}

func (b *protoBuilder) handlePackage(p *proto.Package) {
	b.pkg = p.Name
}

func (b *protoBuilder) handleMessage(m *proto.Message) {
	specs := []composeSpec{}
	if m.Comment != nil {
		for _, each := range m.Comment.Lines {
			if strings.Contains(each, "@compose") {
				spec := newComposeSpec(each)
				specs = append(specs, spec)
			}
		}
	}

	key := fmt.Sprintf("%s.%s", b.pkg, m.Name)
	b.registry[key] = &protoReference{
		pkg:          b.pkg,
		message:      m,
		composeSpecs: specs,
	}
}

func (b *protoBuilder) processComposed() {
	for _, v := range b.registry {
		if len(v.composeSpecs) > 0 {
			// flush all existing
			v.message.Elements = []proto.Visitee{}
			// add according to spec
			for _, each := range v.composeSpecs {
				elem := b.fieldAt(each)
				v.message.Elements = append(v.message.Elements, elem)
			}
			// renumber all
			v.message.Accept(new(renumber))
		}
	}
}

func (b *protoBuilder) fieldAt(spec composeSpec) proto.Visitee {
	msg, ok := b.registry[spec.registryKey]
	if !ok {
		check(fmt.Errorf("not found:[%s]", spec.registryKey))
	}
	f := FieldOfMessage(msg.message, spec.fieldName)
	if f == nil {
		check(fmt.Errorf("not found:[%s]", spec.fieldName))
	}
	return f
}

func FieldOfMessage(m *proto.Message, fieldName string) proto.Visitee {
	for _, each := range m.Elements {
		if f, ok := each.(*proto.NormalField); ok {
			if f.Name == fieldName {
				return f
			}
		}
	}
	return nil
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
