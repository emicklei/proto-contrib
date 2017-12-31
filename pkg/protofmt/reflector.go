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

import "github.com/emicklei/proto"

// reflector is a Visitor that can tell the short type name of a Visitee.
type reflector struct {
	name string
}

// sole instance of reflector
var namer = new(reflector)

func (r *reflector) VisitMessage(m *proto.Message)         { r.name = "Message" }
func (r *reflector) VisitService(v *proto.Service)         { r.name = "Service" }
func (r *reflector) VisitSyntax(s *proto.Syntax)           { r.name = "Syntax" }
func (r *reflector) VisitPackage(p *proto.Package)         { r.name = "Package" }
func (r *reflector) VisitOption(o *proto.Option)           { r.name = "Option" }
func (r *reflector) VisitImport(i *proto.Import)           { r.name = "Import" }
func (r *reflector) VisitNormalField(i *proto.NormalField) { r.name = "NormalField" }
func (r *reflector) VisitEnumField(i *proto.EnumField)     { r.name = "EnumField" }
func (r *reflector) VisitEnum(e *proto.Enum)               { r.name = "Enum" }
func (r *reflector) VisitComment(e *proto.Comment)         { r.name = "Comment" }
func (r *reflector) VisitOneof(o *proto.Oneof)             { r.name = "Oneof" }
func (r *reflector) VisitOneofField(o *proto.OneOfField)   { r.name = "OneOfField" }
func (r *reflector) VisitReserved(rs *proto.Reserved)      { r.name = "Reserved" }
func (r *reflector) VisitRPC(rpc *proto.RPC)               { r.name = "RPC" }
func (r *reflector) VisitMapField(f *proto.MapField)       { r.name = "MapField" }
func (r *reflector) VisitGroup(g *proto.Group)             { r.name = "Group" }
func (r *reflector) VisitExtensions(e *proto.Extensions)   { r.name = "Extensions" }

// nameOfVisitee returns the short type name of a Visitee.
func nameOfVisitee(e proto.Visitee) string {
	e.Accept(namer)
	return namer.name
}
