package protodecode

import (
	"testing"
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
		FieldInt32S: []int32{1, 2, 3, 4},
	}
	result := encodeDecode(m, t)
	list := result["field_int32s"].([]interface{})
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
		FieldInt64S: []int64{1, 2, 3, 4},
	}
	result := encodeDecode(m, t)
	list := result["field_int64s"].([]interface{})
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
		fail(t, got, want)
	}
}

func TestRepeatedFloat(t *testing.T) {
	m := &Test{
		FieldFloats: []float32{3.14, 0.234},
	}
	result := encodeDecode(m, t)
	list := result["field_floats"].([]interface{})
	if got, want := len(list), 2; got != want {
		fail(t, got, want)
	}
	if got, want := list[0].(float32), float32(3.14); got != want {
		fail(t, got, want)
	}
}

func TestString(t *testing.T) {
	m := &Test{
		FieldString: "hello",
	}
	result := encodeDecode(m, t)
	if got, want := result["field_string"], "hello"; got != want {
		fail(t, got, want)
	}
}

func TestRepeatedString(t *testing.T) {
	m := &Test{
		FieldStrings: []string{"hello", "world"},
	}
	result := encodeDecode(m, t)
	list := result["field_strings"].([]interface{})
	if got, want := len(list), 2; got != want {
		fail(t, got, want)
	}
	if got, want := list[1].(string), "world"; got != want {
		fail(t, got, want)
	}
}

func TestBool(t *testing.T) {
	// false value is not written TODO
	m := &Test{
		FieldBool: true,
	}
	result := encodeDecode(m, t)
	if got, want := result["field_bool"], true; got != want {
		fail(t, got, want)
	}
}

func TestRepeatedBool(t *testing.T) {
	m := &Test{
		FieldBools: []bool{true, false},
	}
	result := encodeDecode(m, t)
	list := result["field_bools"].([]interface{})
	if got, want := len(list), 2; got != want {
		fail(t, got, want)
	}
	if got, want := list[1].(bool), false; got != want {
		fail(t, got, want)
	}
}

func TestFoo(t *testing.T) {
	m := &Test{
		FieldFoo: &Foo{Foo: "foo"},
	}
	result := encodeDecode(m, t)
	foo := result["field_foo"].(map[string]interface{})
	if got, want := foo != nil, true; got != want {
		fail(t, got, want)
	}
	if got, want := foo["foo"], "foo"; got != want {
		fail(t, got, want)
	}
}

func TestRepeatedFoo(t *testing.T) {
	m := &Test{
		FieldFoos: []*Foo{&Foo{Foo: "foo1"}, &Foo{Foo: "foo2"}},
	}
	result := encodeDecode(m, t)
	foos := result["field_foos"].([]interface{})
	if got, want := len(foos), 2; got != want {
		fail(t, got, want)
	}
	foo := foos[0].(map[string]interface{})
	if got, want := foo["foo"], "foo1"; got != want {
		fail(t, got, want)
	}
}

func TestMapStringInt32(t *testing.T) {
	m := &Test{
		FieldMapStringInt32: map[string]int32{
			"hello": 1,
			"world": 2,
		},
	}
	result := encodeDecode(m, t)
	field := result["field_map_string_int32"].(map[interface{}]interface{})
	if got, want := field["hello"], uint64(1); got != want {
		fail(t, got, want)
	}
}

func TestMapInt64Foo(t *testing.T) {
	m := &Test{
		FieldMapInt64_Foo: map[int64]*Foo{
			1: &Foo{Foo: "foo1"},
			2: &Foo{Foo: "foo2"},
		},
	}
	result := encodeDecode(m, t)
	//dump(result)
	field := result["field_map_int64_Foo"].(map[string]interface{})
	foo1 := field["1 (uint64)"].(map[string]interface{})
	if got, want := foo1["foo"], "foo1"; got != want {
		fail(t, got, want)
	}
}
