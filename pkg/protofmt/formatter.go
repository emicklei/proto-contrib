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
	"fmt"
	"io"

	"github.com/emicklei/proto"
)

// Formatter visits a Proto and writes formatted source.
type Formatter struct {
	w               io.Writer
	indentSeparator string
	indentLevel     int
	lastStmt        string
	lastLevel       int
}

// NewFormatter returns a new Formatter. Only the indentation separator is configurable.
func NewFormatter(writer io.Writer, indentSeparator string) *Formatter {
	return &Formatter{w: writer, indentSeparator: indentSeparator}
}

// Format visits all proto elements and writes formatted source.
func (f *Formatter) Format(p *proto.Proto) {
	for _, each := range p.Elements {
		each.Accept(f)
	}
}

// VisitComment formats a Comment and writes a newline.
func (f *Formatter) VisitComment(c *proto.Comment) {
	f.printComment(c)
	f.nl()
}

// VisitEnum formats a Enum.
func (f *Formatter) VisitEnum(e *proto.Enum) {
	f.begin("enum", e)
	fmt.Fprintf(f.w, "enum %s {", e.Name)
	if len(e.Elements) > 0 {
		f.nl()
		f.level(1)
		f.printAsGroups(e.Elements)
		f.indent(-1)
	}
	io.WriteString(f.w, "}\n")
	f.end("enum")
}

// VisitEnumField formats a EnumField.
func (f *Formatter) VisitEnumField(e *proto.EnumField) {
	f.printAsGroups([]proto.Visitee{e})
}

// VisitImport formats a Import.
func (f *Formatter) VisitImport(i *proto.Import) {
	f.printAsGroups([]proto.Visitee{i})
}

// VisitMessage formats a Message.
func (f *Formatter) VisitMessage(m *proto.Message) {
	f.begin("message", m)
	if m.IsExtend {
		fmt.Fprintf(f.w, "extend ")
	} else {
		fmt.Fprintf(f.w, "message ")
	}
	fmt.Fprintf(f.w, "%s {", m.Name)
	if len(m.Elements) > 0 {
		f.nl()
		f.level(1)
		f.printAsGroups(m.Elements)
		f.indent(-1)
	}
	io.WriteString(f.w, "}\n")
	f.end("message")
}

// VisitOption formats a Option.
func (f *Formatter) VisitOption(o *proto.Option) {
	f.begin("option", o)
	fmt.Fprintf(f.w, "option %s = ", o.Name)
	if o.AggregatedConstants != nil {
		fmt.Fprintf(f.w, "{\n")
		f.level(1)
		for _, each := range o.AggregatedConstants {
			f.indent(0)
			fmt.Fprintf(f.w, "%s: %s\n", each.Name, each.Literal.SourceRepresentation())
		}
		f.indent(-1)
		fmt.Fprintf(f.w, "}")
	} else {
		// TODO printAs groups with fixed length
		fmt.Fprintf(f.w, o.Constant.SourceRepresentation())
	}
	fmt.Fprintf(f.w, ";")
	if o.InlineComment != nil {
		fmt.Fprintf(f.w, " //%s", o.InlineComment.Message())
	}
	f.nl()
}

// VisitPackage formats a Package.
func (f *Formatter) VisitPackage(p *proto.Package) {
	f.nl()
	f.printAsGroups([]proto.Visitee{p})
}

// VisitService formats a Service.
func (f *Formatter) VisitService(s *proto.Service) {
	f.begin("service", s)
	fmt.Fprintf(f.w, "service %s {", s.Name)
	if len(s.Elements) > 0 {
		f.nl()
		f.level(1)
		f.printAsGroups(s.Elements)
		f.indent(-1)
	}
	io.WriteString(f.w, "}\n")
	f.end("service")
}

// VisitSyntax formats a Syntax.
func (f *Formatter) VisitSyntax(s *proto.Syntax) {
	f.begin("syntax", s)
	fmt.Fprintf(f.w, "syntax = %q", s.Value)
	f.endWithComment(s.InlineComment)
}

// VisitOneof formats a Oneof.
func (f *Formatter) VisitOneof(o *proto.Oneof) {
	f.begin("oneof", o)
	fmt.Fprintf(f.w, "oneof %s {", o.Name)
	if len(o.Elements) > 0 {
		f.nl()
		f.level(1)
		f.printAsGroups(o.Elements)
		f.indent(-1)
	}
	io.WriteString(f.w, "}\n")
	f.end("oneof")
}

// VisitOneofField formats a OneofField.
func (f *Formatter) VisitOneofField(o *proto.OneOfField) {
	f.printAsGroups([]proto.Visitee{o})
}

// VisitReserved formats a Reserved.
func (f *Formatter) VisitReserved(r *proto.Reserved) {
	f.begin("reserved", r)
	io.WriteString(f.w, "reserved ")
	if len(r.Ranges) > 0 {
		for i, each := range r.Ranges {
			if i > 0 {
				io.WriteString(f.w, ", ")
			}
			fmt.Fprintf(f.w, "%s", each.SourceRepresentation())
		}
	} else {
		for i, each := range r.FieldNames {
			if i > 0 {
				io.WriteString(f.w, ", ")
			}
			fmt.Fprintf(f.w, "%q", each)
		}
	}
	f.endWithComment(r.InlineComment)
}

// VisitRPC formats a RPC.
func (f *Formatter) VisitRPC(r *proto.RPC) {
	f.printAsGroups([]proto.Visitee{r})
}

// VisitMapField formats a MapField.
func (f *Formatter) VisitMapField(m *proto.MapField) {
	f.printAsGroups([]proto.Visitee{m})
}

// VisitNormalField formats a NormalField.
func (f *Formatter) VisitNormalField(f1 *proto.NormalField) {
	f.printAsGroups([]proto.Visitee{f1})
}

// VisitGroup formats a proto2 Group.
func (f *Formatter) VisitGroup(g *proto.Group) {
	f.begin("group", g)
	if g.Optional {
		io.WriteString(f.w, "optional ")
	}
	fmt.Fprintf(f.w, "group %s = %d {", g.Name, g.Sequence)
	if len(g.Elements) > 0 {
		f.nl()
		f.level(1)
		f.printAsGroups(g.Elements)
		f.indent(-1)
	}
	io.WriteString(f.w, "}\n")
	f.end("group")
}

// VisitExtensions formats a proto2 Extensions.
func (f *Formatter) VisitExtensions(e *proto.Extensions) {
	f.begin("extensions", e)
	io.WriteString(f.w, "extensions ")
	for i, each := range e.Ranges {
		if i > 0 {
			io.WriteString(f.w, ", ")
		}
		fmt.Fprintf(f.w, "%s", each.SourceRepresentation())
	}
	f.endWithComment(e.InlineComment)
}
