package protofmt

import (
	"fmt"
	"testing"

	"github.com/emicklei/proto"
)

func TestOpenWithMap(t *testing.T) {
	o := new(proto.Option)
	o.Name = "(google.api.http)"
	o.Constant = proto.Literal{}

	get := new(proto.NamedLiteral)
	get.Name = "get"
	get.Literal = new(proto.Literal)
	get.Source = "/hello/world"
	get.IsString = true
	get.PrintsColon = true

	add := new(proto.NamedLiteral)
	add.Name = "additional_bindings"
	add.PrintsColon = true
	add.Literal = new(proto.Literal)

	o.Constant.OrderedMap = append(o.Constant.OrderedMap, get)
	o.Constant.OrderedMap = append(o.Constant.OrderedMap, add)

	got := formatted(o)
	fmt.Println(got)
}
