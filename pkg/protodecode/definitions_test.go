package protodecode

import "testing"

func TestAddFromFile(t *testing.T) {
	d := NewDefinitions()
	d.AddFromFile("test.proto")
	if got, want := len(d.filenamesRead), 1; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	m, ok := d.Message("protodecode", "Test")
	if !ok {
		t.Fail()
	}
	if got, want := m.Name, "Test"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
