package main

import (
	"github.com/emicklei/proto"
)

type renumber struct {
	nr int
	proto.NoopVisitor
}

func (r *renumber) VisitNormalField(f *proto.NormalField) {
	r.nr += 1
	f.Sequence = r.nr
}

func (r *renumber) VisitMessage(m *proto.Message) {
	for _, each := range m.Elements {
		each.Accept(r)
	}
}

type fieldcopier struct {
	proto.NoopVisitor
	copy *proto.NormalField
}

func (c *fieldcopier) VisitNormalField(f *proto.NormalField) {
	field := *f.Field
	c.copy = &proto.NormalField{
		Field: &field,
	}
}
