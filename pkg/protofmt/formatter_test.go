// Copyright (c) 2017 Ernest Micklei
//
// MIT License
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package protofmt

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/emicklei/proto"
)

func newParserOn(def string) *proto.Parser {
	return proto.NewParser(strings.NewReader(def))
}

func TestPrintListOfColumns(t *testing.T) {
	e0 := new(proto.EnumField)
	e0.Name = "A"
	e0.Integer = 1
	op0 := new(proto.Option)
	op0.IsEmbedded = true
	op0.Name = "a"
	op0.Constant = proto.Literal{Source: "1234"}
	e0.ValueOption = op0

	e1 := new(proto.EnumField)
	e1.Name = "ABC"
	e1.Integer = 12
	op1 := new(proto.Option)
	op1.IsEmbedded = true
	op1.Name = "ab"
	op1.Constant = proto.Literal{Source: "1234"}
	e1.ValueOption = op1

	list := []columnsPrintable{asColumnsPrintable(e0), asColumnsPrintable(e1)}
	b := new(bytes.Buffer)
	f := NewFormatter(b)
	f.printListOfColumns(list)
	formatted := `A   =  1 [a  = 1234];
ABC = 12 [ab = 1234];
`
	if got, want := b.String(), formatted; got != want {
		t.Errorf("got [%v] want [%v]", got, want)
	}
}

func TestFormatCStyleComment(t *testing.T) {
	src := `/*
 * Hello
 * World
 */
`
	def, _ := proto.NewParser(strings.NewReader(src)).Parse()
	b := new(bytes.Buffer)
	f := NewFormatter(b, WithIndentSeparator(" "))
	f.Format(def)
	if got, want := formatted(def.Elements[0]), src; got != want {
		println(diff(got, want))
		t.Fail()
	}
}

func TestFormatExtendMessage(t *testing.T) {
	src := `// extend
extend google.protobuf.MessageOptions {
  
  // my_option
  optional string my_option = 51234; // mynumber
  
  // other
  string field      = 12;
  string no_comment = 13;
}

`
	p := newParserOn(src)
	pp, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	m, ok := pp.Elements[0].(*proto.Message)
	if !ok {
		t.Fatal("message expected")
	}
	if got, want := formatted(m), src; got != want {
		fmt.Println(diff(got, want))
		fmt.Println("<" + got + ">")
		fmt.Println("<" + want + ">")
		t.Fail()
	}
}

func TestFormatAggregatedOptionSyntax(t *testing.T) {
	// TODO format not that nice
	src := `service AggregateOption {
  rpc Find (Finder) returns (stream Result) {
    option (google.api.http) = {
      post: "/v1/finders/1"
      body: "*"
    };
  
  }
}
`
	p := newParserOn(src)
	def, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	svc := def.Elements[0].(*proto.Service)
	if got, want := len(svc.Elements), 1; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	if got, want := formatted(def.Elements[0]), src; got != want {
		fmt.Println(diff(got, want))
		fmt.Println("--- got")
		fmt.Println(got)
		fmt.Println("--- want")
		fmt.Println(want)
		t.Fail()
	}
}

func TestFormatCommentSample(t *testing.T) {
	proto := `
/*
 begin
*/

// comment 1
// comment 2
syntax = "proto"; // inline 1

// comment 3
// comment 4
package test; // inline 2

// comment 5
// comment 6
message Test {
    // comment 7
    // comment 8
    int64 i = 1; // inline 3
}
	
/// triple
`
	p := newParserOn(proto)
	def, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(def.Elements), 5; got != want {
		t.Fatalf("got [%v] want [%v]", got, want)
	}
	b := new(bytes.Buffer)
	f := NewFormatter(b, WithIndentSeparator("  ")) // 2 spaces
	f.Format(def)
}

// https://github.com/emicklei/proto-contrib/issues/8
func TestOptionWithStructureAndTwoFields(t *testing.T) {
	src := `service X {
  rpc Hello (google.protobuf.Empty) returns (google.protobuf.Empty) {
    option simple = "easy";

    option (google.api.http) = {
      get: "/hello"
      additional_bindings: {
        get: "/hello/world"
      }
    };

  }
}`
	p := newParserOn(src)
	def, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := formatted(def.Elements[0]), src; got != want {
		fmt.Println(diff(got, want))
		fmt.Println("--- got")
		fmt.Println(got)
		fmt.Println("--- want")
		fmt.Println(want)
		//t.Fail()
	}
}

func TestOptionTrue(t *testing.T) {
	src := `option alive = true;
`
	p := newParserOn(src)
	def, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := formatted(def.Elements[0]), src; got != want {
		fmt.Println(diff(got, want))
		t.Fail()
	}
}

func TestFormatMaps(t *testing.T) {
	src := `message A {
  bool                done  = 1;
  map <string,string> smap1 = 2;
  map <string,string> smap2 = 3;
  
  // comment
  map <string,string> smap3 = 4;
}

`
	p := newParserOn(src)
	pp, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	m, ok := pp.Elements[0].(*proto.Message)
	if !ok {
		t.Fatal("message expected")
	}
	if got, want := formatted(m), src; got != want {
		fmt.Println(diff(got, want))
		fmt.Println("<" + got + ">")
		fmt.Println("<" + want + ">")
		t.Fail()
	}
}
