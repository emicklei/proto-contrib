package main

import (
	"github.com/emicklei/proto"
)

type renumber struct {
	reading       bool
	nr            int
	fieldNumbers  []int
	fieldMappings map[string]int
	proto.NoopVisitor
}

func newRenumber() *renumber {
	return &renumber{nr: 1}
}

func (r *renumber) VisitMapField(f *proto.MapField) {
	f.Sequence = r.nr
	r.nr++
}

func (r *renumber) VisitNormalField(f *proto.NormalField) {
	f.Sequence = r.nr
	r.nr++
}

func (r *renumber) VisitMessage(m *proto.Message) {
	for _, each := range m.Elements {
		each.Accept(r)
	}
}

type fieldcopier struct {
	proto.NoopVisitor
	copy proto.Visitee
}

func (c *fieldcopier) VisitNormalField(f *proto.NormalField) {
	field := *f.Field
	c.copy = &proto.NormalField{
		Repeated: f.Repeated,
		Optional: f.Optional,
		Field:    &field,
	}
}

func (c *fieldcopier) VisitMapField(f *proto.MapField) {
	field := *f.Field
	c.copy = &proto.MapField{
		Field:   &field,
		KeyType: f.KeyType,
	}
}
