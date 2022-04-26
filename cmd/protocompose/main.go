package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/emicklei/proto"
	"github.com/emicklei/proto-contrib/pkg/protofmt"
)

func main() {
	b := &protoBuilder{
		registry: make(map[string]*protoReference),
	}
	// quick
	var last *proto.Proto
	for _, each := range os.Args[1:] {
		last = b.loadProto(each)
	}
	b.processComposed()
	fmt.Println(formatted(last))
}

func FieldOfMessage(m *proto.Message, fieldName string) proto.Visitee {
	for _, each := range m.Elements {
		if f, ok := each.(*proto.NormalField); ok {
			if f.Name == fieldName {
				return f
			}
		}
	}
	return nil
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func formatted(v proto.Visitee) string {
	b := new(bytes.Buffer)
	f := protofmt.NewFormatter(b, "  ") // 2 spaces
	v.Accept(f)
	return b.String()
}
