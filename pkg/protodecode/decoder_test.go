package protodecode

import (
	"testing"

	proto "github.com/golang/protobuf/proto"
)

func TestInt32(t *testing.T) {
	m := &Test{
		FieldInt32: 42,
	}
	result := encodeDecode(m, t)
	if got, want := result["field_int32"], uint64(42); got != want {
		t.Errorf("got %v (%T) %v (%T)", got, got, want, want)
	}
}

func TestString(t *testing.T) {
	m := &Test{
		FieldString: "hello",
	}
	result := encodeDecode(m, t)
	if got, want := result["field_string"], "hello"; got != want {
		t.Errorf("got %v (%T) %v (%T)", got, got, want, want)
	}
}

func encodeDecode(m *Test, t *testing.T) map[string]interface{} {
	data, err := proto.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	defs := NewDefinitions()
	defs.AddFromFile("test.proto")

	dec := NewDecoder(defs, proto.NewBuffer(data))
	result, err := dec.Decode("protodecode", "Test")
	if err != nil && err != EOM {
		t.Fatal(err)
	}
	return result
}
