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
	Name    string      `json:"name"`
	Doc     string      `json:"doc,omitempty"`
	Type    interface{} `json:"type"` // string or array of type
	Default string      `json:"default,omitempty"`
	Order   string      `json:"order,omitempty"` // ascending,descending,ignore
	Aliases []string    `json:"aliases,omitempty"`
}

type FieldComplexType struct {
	Type  string      `json:"type"`
	Items interface{} `json:"items"`
}

type Builder struct {
	records  map[string]Record
	messages map[string]*proto.Message
}

func NewBuilder() Builder {
	return Builder{
		records:  map[string]Record{},
		messages: map[string]*proto.Message{},
	}
}

func (b Builder) AddMessage(m *proto.Message) {
	b.messages[m.Name] = m
}

func (b Builder) Build(name string) (Record, bool) {
	if r, ok := b.records[name]; ok {
		// already build
		return r, true
	}
	// new record
	m := b.messages[name] // TODO handle missing
	r := Record{}
	r.Type = "record"
	r.Namespace = "main"
	if m.Comment != nil {
		r.Doc = m.Comment.Message()
	}
	r.Name = m.Name
	for _, each := range m.Elements {
		if normalField, ok := each.(*proto.NormalField); ok {
			r.Fields = append(r.Fields, b.toRecordField(normalField))
		}
	}
	b.records[m.Name] = r
	return r, true
}

func (b Builder) toRecordField(f *proto.NormalField) Field {
	if f.Repeated {
		var itemType interface{}
		itemType, ok := toRecordFieldType(f.Type)
		if !ok {
			itemType, ok = b.Build(f.Type)
		}
		if !ok {
			return Field{Name: f.Name, Type: "ERROR: missing " + f.Type}
		}
		return Field{
			Name: f.Name,
			Type: FieldComplexType{
				Type:  "array",
				Items: itemType,
			},
		}
	}
	rf := Field{
		Name: f.Name,
	}
	rf.Type, _ = toRecordFieldType(f.Type)
	if f.Comment != nil {
		rf.Doc = f.Comment.Message()
	}
	return rf
}

func toOptionalField(typeName string) []string {
	return []string{"null", typeName}
}

// https://developers.google.com/protocol-buffers/docs/proto#scalar
func toRecordFieldType(typeName string) (string, bool) {
	switch typeName {
	case "bool":
		return "boolean", true
	case "uint16", "sint32", "int32", "fixed32", "sfixed32":
		return "int", true
	case "int64", "uint64", "sint64", "fixed64", "sfixed64":
		return "long", true
	case "double", "float", "string", "bytes":
		return typeName, true
	default:
		return typeName, false
	}
}
