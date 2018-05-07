package proto2avro

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/emicklei/proto"
)

func TestMessageToRecord1(t *testing.T) {
	src := `// Wonderful example
			message Test {
			// that's all
			string content = 1;
			repeated int32 bits = 2;
		}`

	parser := proto.NewParser(strings.NewReader(src))
	definition, _ := parser.Parse()
	proto.Walk(definition, proto.WithMessage(func(m *proto.Message) {
		r := MessageToRecord(m)
		json.NewEncoder(os.Stdout).Encode(r)
	}))
}
