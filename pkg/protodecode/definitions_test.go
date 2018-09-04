package protodecode

import "testing"

func TestAddFromFile(t *testing.T) {
	d := NewDefinitions()
	d.AddFromFile("test.proto")
	if got, want := len(d.filenamesRead), 1; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := d.Message("protodecode", "Test").Name, "Test"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
