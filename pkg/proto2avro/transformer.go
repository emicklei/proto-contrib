package proto2avro

import (
	"github.com/emicklei/proto"
)

// Record models an Avro record type
// https://avro.apache.org/docs/1.8.1/spec.html
type Record struct {
	Type      string   `json:"type"`
	Name      string   `json:"name"`
	Namespace string   `json:"namespace,omitempty"`
	Doc       string   `json:"doc,omitempty"`
	Aliases   []string `json:"aliases,omitempty"`
	Fields    []Field  `json:"fields"`
}

// Field models an Avro record field
// https://avro.apache.org/docs/1.8.1/spec.html
type Field struct {
	Name    string   `json:"name"`
	Doc     string   `json:"doc,omitempty"`
	Type    []string `json:"type,omitempty"`
	Default string   `json:"default,omitempty"`
	Order   string   `json:"order,omitempty"` // ascending,descending,ignore
	Aliases []string `json:"aliases,omitempty"`
}

func MessageToRecord(m *proto.Message) Record {
	r := Record{}
	r.Type = "record"
	r.Namespace = "main"
	if m.Comment != nil {
		r.Doc = m.Comment.Message()
	}
	r.Name = m.Name
	for _, each := range m.Elements {
		if normalField, ok := each.(*proto.NormalField); ok {
			r.Fields = append(r.Fields, toRecordField(normalField))
		}
	}
	return r
}

func toRecordField(f *proto.NormalField) Field {
	if f.Repeated {
		return Field{
			Name: f.Name,
			Type: toOptionalField("array"),
		}
	}
	rf := Field{
		Name: f.Name,
		Type: toOptionalField(toRecordFieldType(f.Type)),
	}
	if f.Comment != nil {
		rf.Doc = f.Comment.Message()
	}
	return rf
}

func toOptionalField(typeName string) []string {
	return []string{"null", typeName}
}

// https://developers.google.com/protocol-buffers/docs/proto#scalar
func toRecordFieldType(typeName string) string {
	switch typeName {
	case "bool":
		return "boolean"
	case "uint16", "sint32", "int32", "fixed32", "sfixed32":
		return "int"
	case "uint64", "sint64", "fixed64", "sfixed64":
		return "int"
	case "double", "float", "string", "bytes":
		return typeName
	default:
		return "object"
	}
}
