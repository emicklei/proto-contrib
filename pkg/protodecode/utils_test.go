package protodecode

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	proto "github.com/golang/protobuf/proto"
)

func fail(t *testing.T, got, want interface{}) {
	t.Helper()
	t.Fatalf("got %v (%T) want %v (%T)", got, got, want, want)
}

func dump(what interface{}) {
	fmt.Println(what)
	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "\t")
	if err := e.Encode(what); err != nil {
		log.Println(err)
	}
}

func encodeDecode(m *Test, t *testing.T) map[string]interface{} {
	t.Helper()
	data, err := proto.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	defs := NewDefinitions()
	defs.AddFromFile("test.proto")

	dec := NewDecoder(defs, proto.NewBuffer(data))
	// dec.verbose = true
	result, err := dec.Decode("protodecode", "Test")
	if err != nil && err != ErrEndOfMessage {
		t.Fatal(err)
	}
	return result
}
