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
	get.QuoteRune = rune('"')
	get.Source = "/hello"
	get.IsString = true
	get.PrintsColon = true

	get2 := new(proto.NamedLiteral)
	get2.Name = "get"
	get2.Literal = new(proto.Literal)
	get2.Source = "/hello/world"
	get2.QuoteRune = rune('"')
	get2.IsString = true
	get2.PrintsColon = true

	add := new(proto.NamedLiteral)
	add.Name = "additional_bindings"
	add.PrintsColon = true
	add.Literal = new(proto.Literal)
	add.Literal.OrderedMap = append(o.Constant.OrderedMap, get2)

	o.Constant.OrderedMap = append(o.Constant.OrderedMap, get)
	o.Constant.OrderedMap = append(o.Constant.OrderedMap, add)

	got := formatted(o)
	fmt.Println(got)

	want := `option (google.api.http) = {
  get: "/hello"
  additional_bindings: {
    get: "/hello/world"
  }
};
`
	if got != want {
		fmt.Println(diff(got, want))
		fmt.Println("--- got")
		fmt.Println(got)
		fmt.Println("--- want")
		fmt.Println(want)
		t.Fail()
	}
}
