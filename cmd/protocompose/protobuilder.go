package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/emicklei/proto"
)

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

func newProtoBuilder() *protoBuilder {
	return &protoBuilder{
		registry: make(map[string]*protoReference),
	}
}

func (b *protoBuilder) loadProto(absFilename string) *proto.Proto {
	log.Println("loading", absFilename)
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
	log.Println("handling", m.Name)
	specs := []composeSpec{}
	if m.Comment != nil {
		for _, each := range m.Comment.Lines {
			if strings.Contains(each, "@compose") {
				spec := newComposeSpec(each)
				log.Println("... compose with", spec.fieldName)
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
			// flush all
			v.message.Elements = []proto.Visitee{}

			// add according to spec
			for _, each := range v.composeSpecs {
				elem := b.copiedFieldAt(each)
				v.message.Elements = append(v.message.Elements, elem)
			}
			// renumber all
			v.message.Accept(newRenumber())
		}
	}
}

func (b *protoBuilder) messageAt(key string) *proto.Message {
	msg, ok := b.registry[key]
	if !ok {
		check(fmt.Errorf("message not found:[%s]", key))
	}
	return msg.message
}

func (b *protoBuilder) copiedFieldAt(spec composeSpec) proto.Visitee {
	msg, ok := b.registry[spec.registryKey]
	if !ok {
		check(fmt.Errorf("message not found:[%s]", spec.registryKey))
	}
	f := FieldOfMessage(msg.message, spec.fieldName)
	if f == nil {
		check(fmt.Errorf("field not found:[%s]", spec.fieldName))
	}
	copier := new(fieldcopier)
	f.Accept(copier)
	return copier.copy
}

func (b *protoBuilder) copiedFieldsAt(spec composeSpec) (list []proto.Visitee) {
	msg, ok := b.registry[spec.registryKey]
	if !ok {
		check(fmt.Errorf("message not found:[%s]", spec.registryKey))
	}
	for _, each := range FieldNamesOfMessage(msg.message) {
		list = append(list, b.copiedFieldAt(spec.forFieldName(each)))
	}
	return
}

// somepackage.v2.FileReference
func fieldNameFromMessage(fullType string) string {
	name := fullType[strings.LastIndex(fullType, ".")+1:]
	return strings.ToLower(name)
}
