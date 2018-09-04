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

func TestRepeatedInt32(t *testing.T) {
	m := &Test{
		FieldsInt32: []int32{1, 2, 3, 4},
	}
	result := encodeDecode(m, t)
	list := result["fields_int32"].([]interface{})
	if got, want := len(list), 4; got != want {
		t.Errorf("got %v (%T) %v (%T)", got, got, want, want)
	}
	if got, want := list[0].(int32), int32(1); got != want {
		t.Errorf("got %v (%T) %v (%T)", got, got, want, want)
	}
}

func TestInt64(t *testing.T) {
	m := &Test{
		FieldInt64: 42,
	}
	result := encodeDecode(m, t)
	if got, want := result["field_int64"], uint64(42); got != want {
		t.Errorf("got %v (%T) %v (%T)", got, got, want, want)
	}
}

func TestRepeatedInt64(t *testing.T) {
	m := &Test{
		FieldsInt64: []int64{1, 2, 3, 4},
	}
	result := encodeDecode(m, t)
	list := result["fields_int64"].([]interface{})
	if got, want := len(list), 4; got != want {
		t.Errorf("got %v (%T) %v (%T)", got, got, want, want)
	}
	if got, want := list[1].(uint64), uint64(2); got != want {
		t.Errorf("got %v (%T) %v (%T)", got, got, want, want)
	}
}

func TestFloat(t *testing.T) {
	m := &Test{
		FieldFloat: 3.14,
	}
	result := encodeDecode(m, t)
	if got, want := result["field_float"], float32(3.14); got != want {
		t.Errorf("got %v (%T) %v (%T)", got, got, want, want)
	}
}

func TestRepeatedFloat(t *testing.T) {
	m := &Test{
		FieldsFloat: []float32{3.14, 0.234},
	}
	result := encodeDecode(m, t)
	list := result["fields_float"].([]interface{})
	if got, want := len(list), 2; got != want {
		t.Errorf("got %v (%T) %v (%T)", got, got, want, want)
	}
	if got, want := list[0].(float32), float32(3.14); got != want {
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

func TestRepeatedString(t *testing.T) {
	m := &Test{
		FieldsString: []string{"hello", "world"},
	}
	result := encodeDecode(m, t)
	t.Logf("%#v", result)
	list := result["fields_string"].([]interface{})
	if got, want := len(list), 2; got != want {
		t.Errorf("got %v (%T) %v (%T)", got, got, want, want)
	}
	if got, want := list[1].(string), "world"; got != want {
		t.Errorf("got %v (%T) %v (%T)", got, got, want, want)
	}
}

func TestBool(t *testing.T) {
	// false value is not written TODO
	m := &Test{
		FieldBool: true,
	}
	result := encodeDecode(m, t)
	if got, want := result["field_bool"], true; got != want {
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
	dec.verbose = true
	result, err := dec.Decode("protodecode", "Test")
	if err != nil && err != EOM {
		t.Fatal(err)
	}
	return result
}
