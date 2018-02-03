package proto2gql

import (
	"bytes"
	"github.com/emicklei/proto"
	"io"
	"strings"
)

var BUILTINS = map[string]string{
	"double":   "Float",
	"float":    "Float",
	"int32":    "Int",
	"int64":    "Int",
	"uint32":   "Int",
	"uint64":   "Int",
	"sint32":   "Int",
	"sint64":   "Int",
	"fixed32":  "Int",
	"fixed64":  "Int",
	"sfixed32": "Int",
	"sfixed64": "Int",
	"bool":     "Boolean",
	"string":   "String",
	"bytes":    "[String]",
}

type (
	Visitor struct {
		scope    *Scope
		buff     *bytes.Buffer
		children []*Visitor
		filter   Filter
	}
)

func NewVisitor(converter *Converter, filter Filter) *Visitor {
	return &Visitor{
		buff:     new(bytes.Buffer),
		children: make([]*Visitor, 0, 5),
		scope:    NewScope(converter),
		filter:   filter,
	}
}

func (v *Visitor) Fork(name string) *Visitor {
	child := &Visitor{
		buff:     new(bytes.Buffer),
		children: make([]*Visitor, 0, 5),
		scope:    v.scope.Fork(name),
		filter:   v.filter,
	}

	v.children = append(v.children, child)

	return child
}

func (v *Visitor) Flush(out io.Writer) {
	out.Write(v.buff.Bytes())

	v.buff.Reset()

	for _, child := range v.children {
		child.Flush(out)
	}
}

func (v *Visitor) VisitMessage(m *proto.Message) {
	// we add it to be able to resolve it in fields
	v.scope.AddLocalType(m.Name)

	if v.canTransformMessage(m) == false {
		return
	}

	v.buff.WriteString("\n")

	v.buff.WriteString("type " + v.scope.converter.NewTypeName(v.scope, m.Name) + " {\n")

	fields := make([]*proto.NormalField, 0, len(m.Elements))

	for _, element := range m.Elements {

		field, ok := element.(*proto.NormalField)

		// it's not a nested message/enum
		if ok == true {
			// we put it in array in order to process nested messages first
			// in case they exist and have them in a scope
			fields = append(fields, field)
		} else {
			// if so, create a nested visitor
			// we need to track a parent's convertedName
			// in order to generate a unique convertedName for nested ones
			// we create another visitor in order to unfold nested types since GraphQL does not support nested types
			element.Accept(v.Fork(m.Name))
		}
	}

	// now, having all nested messages in a scope, we can transform fields
	for _, field := range fields {
		field.Accept(v)
	}

	v.buff.WriteString("}\n")

}
func (v *Visitor) VisitService(s *proto.Service) {}
func (v *Visitor) VisitSyntax(s *proto.Syntax)   {}
func (v *Visitor) VisitPackage(p *proto.Package) {
	v.scope.SetPackageName(p.Name)
}
func (v *Visitor) VisitOption(o *proto.Option) {}
func (v *Visitor) VisitImport(i *proto.Import) {
	v.scope.AddImportedType(i.Filename)
}
func (v *Visitor) VisitNormalField(field *proto.NormalField) {
	if v.canTransformMessageField(field) == false {
		return
	}

	v.buff.WriteString("    " + field.Name + ":")

	typeName := v.scope.ResolveConvertedTypeName(field.Type)

	if field.Repeated == false {
		v.buff.WriteString(" " + typeName)
	} else {
		v.buff.WriteString(" [" + typeName + "]")
	}

	if field.Required == true {
		v.buff.WriteString("!")
	}

	v.buff.WriteString("\n")
}
func (v *Visitor) VisitEnumField(i *proto.EnumField) {
	v.buff.WriteString("    " + i.Name + "\n")
}
func (v *Visitor) VisitEnum(e *proto.Enum) {
	// we add it to be able to resolve it in fields
	v.scope.AddLocalType(e.Name)

	if v.canTransformEnum(e) == false {
		return
	}

	v.buff.WriteString("\n")

	v.buff.WriteString("enum " + v.scope.converter.NewTypeName(v.scope, e.Name) + " {\n")

	for _, element := range e.Elements {
		element.Accept(v)
	}

	v.buff.WriteString("}\n")
}
func (v *Visitor) VisitComment(e *proto.Comment)       {}
func (v *Visitor) VisitOneof(o *proto.Oneof)           {}
func (v *Visitor) VisitOneofField(o *proto.OneOfField) {}
func (v *Visitor) VisitReserved(r *proto.Reserved)     {}
func (v *Visitor) VisitRPC(r *proto.RPC)               {}
func (v *Visitor) VisitMapField(f *proto.MapField)     {}

// proto2
func (v *Visitor) VisitGroup(g *proto.Group)           {}
func (v *Visitor) VisitExtensions(e *proto.Extensions) {}

func (v *Visitor) canTransformMessage(m *proto.Message) bool {
	return v.filter(v.scope.converter.OriginalFullTypeName(v.scope, m.Name))
}

func (v *Visitor) canTransformMessageField(m *proto.NormalField) bool {
	// ignore builtins
	_, builtin := BUILTINS[m.Type]

	if builtin == true {
		return true
	}

	if strings.Contains(m.Type, ".") {
		return v.filter(m.Type)
	}

	return v.filter(v.scope.ResolveFullTypeName(m.Type))
}

func (v *Visitor) canTransformEnum(e *proto.Enum) bool {
	return v.filter(v.scope.converter.OriginalFullTypeName(v.scope, e.Name))
}
