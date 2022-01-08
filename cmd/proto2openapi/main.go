package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/emicklei/proto"
	"github.com/go-openapi/spec"
)

func main() {
	flag.Parse()
	filename := flag.Args()[0]
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()
	parser := proto.NewParser(file)
	parser.Filename(filename)
	def, err := parser.Parse()
	if err != nil {
		log.Fatalln(err)
	}
	b := newBuilder()
	proto.Walk(def,
		proto.WithService(b.handleService),
		proto.WithMessage(b.handleMessage))

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "\t")
	enc.Encode(b.swagger)
}

type builder struct {
	swagger spec.Swagger
}

func newBuilder() *builder {
	return &builder{swagger: spec.Swagger{}}
}

func (b *builder) handleService(s *proto.Service) {
	fmt.Println(s.Name)
}

func (b *builder) handleMessage(m *proto.Message) {
	fmt.Println(m.Name)
}
