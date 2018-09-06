package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/emicklei/proto"
	"github.com/emicklei/proto-contrib/pkg/proto2avro"
)

func main() {
	if len(os.Args) <= 2 {
		log.Fatalln("proto2avro [filename].proto [message]")
	}
	input, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	parser := proto.NewParser(input)
	definition, _ := parser.Parse()
	builder := proto2avro.NewBuilder()
	proto.Walk(definition, proto.WithMessage(func(m *proto.Message) {
		builder.AddMessage(m)
	}))
	record, ok := builder.Build(os.Args[2])
	if !ok {
		log.Fatalln("no definition found for", os.Args[2])
	}
	e := json.NewEncoder(os.Stdout)
	e.SetIndent("", "\t")
	e.Encode(record)
}
