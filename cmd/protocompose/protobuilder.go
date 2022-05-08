package main

import (
	"fmt"
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
				if each.inlineFields {
					// TODO check existing fields to handle future fields
					v.message.Elements = append(v.message.Elements, b.copiedsFieldsAt(each)...)
				} else if each.embedMessage {
					f := &proto.NormalField{
						Field: &proto.Field{
							Comment: b.messageAt(each.registryKey).Comment,
							Name:    fieldNameFromMessage(each.registryKey),
							Type:    each.registryKey,
						},
					}
					v.message.Elements = append(v.message.Elements, f)
				} else {
					elem := b.copiedFieldAt(each)
					v.message.Elements = append(v.message.Elements, elem)
				}
			}
			// renumber all
			v.message.Accept(new(renumber))
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

func (b *protoBuilder) copiedsFieldsAt(spec composeSpec) (list []proto.Visitee) {
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
