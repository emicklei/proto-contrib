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
	"io"
	"strconv"

	"github.com/emicklei/proto"
)

func columns(v proto.Visitee) []aligned {
	return asColumnsPrintable(v).columns()
}

// columnsPrintable is for elements that can be printed in aligned columns.
type columnsPrintable interface {
	columns() (cols []aligned)
}

type columnsPrinter struct {
	cols    []aligned
	visitee proto.Visitee
}

func asColumnsPrintable(v proto.Visitee) *columnsPrinter {
	p := new(columnsPrinter)
	p.visitee = v
	return p
}

// columns is part of columnsPrintable
func (p *columnsPrinter) columns() []aligned {
	p.visitee.Accept(p)
	return p.cols
}

func (p *columnsPrinter) VisitMessage(m *proto.Message) {}
func (p *columnsPrinter) VisitService(v *proto.Service) {}
func (p *columnsPrinter) VisitSyntax(s *proto.Syntax)   {}
func (p *columnsPrinter) VisitPackage(pkg *proto.Package) {
	p.cols = append(p.cols, notAligned("package "), notAligned(pkg.Name), alignedSemicolon)
	if pkg.InlineComment != nil {
		p.cols = append(p.cols, notAligned(" //"), notAligned(pkg.InlineComment.Message()))
	}
}
func (p *columnsPrinter) VisitOption(o *proto.Option) {
	if !o.IsEmbedded {
		p.cols = append(p.cols, leftAligned("option "))
	} else {
		p.cols = append(p.cols, leftAligned(" ["))
	}
	p.cols = append(p.cols, keyValuePair(o, o.IsEmbedded)...)
	if o.IsEmbedded {
		p.cols = append(p.cols, leftAligned("]"))
	}
	if !o.IsEmbedded {
		p.cols = append(p.cols, alignedSemicolon)
		if o.InlineComment != nil {
			p.cols = append(p.cols, notAligned(" //"), notAligned(o.InlineComment.Message()))
		}
	}
}
func (p *columnsPrinter) VisitImport(i *proto.Import) {
	p.cols = append(p.cols, leftAligned("import"), alignedSpace)
	if len(i.Kind) > 0 {
		p.cols = append(p.cols, leftAligned(i.Kind), alignedSpace)
	}
	p.cols = append(p.cols, notAligned(fmt.Sprintf("%q", i.Filename)), alignedSemicolon)
	if i.InlineComment != nil {
		p.cols = append(p.cols, notAligned(" //"), notAligned(i.InlineComment.Message()))
	}
}

// VisitNormalField
// [|repeated][|optional][space][name][equals][sequence][|option]
func (p *columnsPrinter) VisitNormalField(f *proto.NormalField) {
	if f.Repeated {
		p.cols = append(p.cols, leftAligned("repeated "))
	} else if f.Optional {
		p.cols = append(p.cols, leftAligned("optional "))
	} else if f.Required {
		p.cols = append(p.cols, leftAligned("required "))
	} else {
		p.cols = append(p.cols, alignedEmpty)
	}
	p.cols = append(p.cols, leftAligned(f.Type), alignedSpace, leftAligned(f.Name), alignedEquals, rightAligned(strconv.Itoa(f.Sequence)))
	if len(f.Options) > 0 {
		p.cols = append(p.cols, leftAligned(" ["))
		for i, each := range f.Options {
			if i > 0 {
				p.cols = append(p.cols, alignedComma)
			}
			p.cols = append(p.cols, keyValuePair(each, true)...)
		}
		p.cols = append(p.cols, leftAligned("]"))
	}
	p.cols = append(p.cols, alignedSemicolon)
	if f.InlineComment != nil {
		p.cols = append(p.cols, alignedInlinePrefix(f.InlineComment), notAligned(f.InlineComment.Message()))
	}
}

func (p *columnsPrinter) VisitEnumField(f *proto.EnumField) {
	p.cols = append(p.cols, leftAligned(f.Name), alignedEquals, rightAligned(strconv.Itoa(f.Integer)))
	if f.ValueOption != nil {
		p.cols = append(p.cols, columns(f.ValueOption)...)
	}
	p.cols = append(p.cols, alignedSemicolon)
	if f.InlineComment != nil {
		p.cols = append(p.cols, alignedInlinePrefix(f.InlineComment), notAligned(f.InlineComment.Message()))
	}
}
func (p *columnsPrinter) VisitEnum(e *proto.Enum)       {}
func (p *columnsPrinter) VisitComment(e *proto.Comment) {}
func (p *columnsPrinter) VisitOneof(o *proto.Oneof)     {}
func (p *columnsPrinter) VisitOneofField(o *proto.OneOfField) {
	p.cols = append(p.cols,
		leftAligned(o.Type),
		alignedSpace,
		leftAligned(o.Name),
		alignedEquals,
		rightAligned(strconv.Itoa(o.Sequence)))
	if len(o.Options) > 0 {
		p.cols = append(p.cols, leftAligned(" ["))
		for i, each := range o.Options {
			if i > 0 {
				p.cols = append(p.cols, alignedComma)
			}
			p.cols = append(p.cols, keyValuePair(each, true)...)
		}
		p.cols = append(p.cols, leftAligned("]"))
	}
	p.cols = append(p.cols, alignedSemicolon)
	if o.InlineComment != nil {
		p.cols = append(p.cols, notAligned(" //"), notAligned(o.InlineComment.Message()))
	}
}
func (p *columnsPrinter) VisitReserved(rs *proto.Reserved) {}
func (p *columnsPrinter) VisitRPC(r *proto.RPC) {
	p.cols = append(p.cols,
		leftAligned("rpc "),
		leftAligned(r.Name),
		leftAligned(" ("))
	if r.StreamsRequest {
		p.cols = append(p.cols, leftAligned("stream "))
	} else {
		p.cols = append(p.cols, alignedEmpty)
	}
	p.cols = append(p.cols,
		leftAligned(r.RequestType),
		leftAligned(") "),
		leftAligned("returns"),
		leftAligned(" ("))
	if r.StreamsReturns {
		p.cols = append(p.cols, leftAligned("stream "))
	} else {
		p.cols = append(p.cols, alignedEmpty)
	}
	p.cols = append(p.cols,
		leftAligned(r.ReturnsType),
		leftAligned(")"))
	if len(r.Elements) > 0 {
		buf := new(bytes.Buffer)
		io.WriteString(buf, " {\n")
		f := NewFormatter(buf, WithIndentSeparator("  ")) // TODO get separator, now 2 spaces
		f.level(1)
		for _, each := range r.Elements {
			each.Accept(f)
			io.WriteString(buf, "\n")
		}
		f.indent(-1)
		io.WriteString(buf, "}")
		p.cols = append(p.cols, notAligned(buf.String()))
	} else {
		p.cols = append(p.cols, alignedSemicolon)
	}
	if r.InlineComment != nil {
		p.cols = append(p.cols, notAligned(" //"), notAligned(r.InlineComment.Message()))
	}
}
func (p *columnsPrinter) VisitMapField(f *proto.MapField) {
	p.cols = append(p.cols,
		alignedEmpty, // no repeated no optional
		leftAligned(fmt.Sprintf("map <%s,%s>", f.KeyType, f.Type)),
		alignedSpace,
		leftAligned(f.Name),
		alignedEquals,
		rightAligned(strconv.Itoa(f.Sequence)))
	if len(f.Options) > 0 {
		p.cols = append(p.cols, leftAligned(" ["))
		for i, each := range f.Options {
			if i > 0 {
				p.cols = append(p.cols, alignedComma)
			}
			p.cols = append(p.cols, keyValuePair(each, true)...)
		}
		p.cols = append(p.cols, leftAligned("]"))
	}
	p.cols = append(p.cols, alignedSemicolon)
	if f.InlineComment != nil {
		p.cols = append(p.cols, alignedInlinePrefix(f.InlineComment), notAligned(f.InlineComment.Message()))
	}
}
func (p *columnsPrinter) VisitGroup(g *proto.Group)           {}
func (p *columnsPrinter) VisitExtensions(e *proto.Extensions) {}
