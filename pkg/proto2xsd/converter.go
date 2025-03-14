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

package proto2xsd

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/emicklei/proto"
)

// Convert reads a Proto definition from input and writes XSD Type definition and XSD Root elements to output.
func Convert(filename, xsdNamespace string, input io.Reader, output io.Writer) error {
	parser := proto.NewParser(input)
	parser.Filename(filename)
	def, err := parser.Parse()
	if err != nil {
		return err
	}
	types, err := BuildXSDTypes(def)
	if err != nil {
		return err
	}
	fmt.Fprint(output, xml.Header)
	elements, err := buildXSDElements(def)
	if err != nil {
		return err
	}
	schema := BuildXSDSchema(xsdNamespace)
	schema.Types = types
	schema.Elements = elements
	data, err := xml.MarshalIndent(schema, "", "\t")
	if err != nil {
		return err
	}
	output.Write(data)

	return nil
}
