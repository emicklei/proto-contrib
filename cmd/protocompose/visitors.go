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
	return &renumber{
		fieldMappings: map[string]int{},
	}
}

func (r *renumber) VisitNormalField(f *proto.NormalField) {
	if r.reading {
		r.fieldMappings[f.Name] = f.Sequence
		r.fieldNumbers = append(r.fieldNumbers, f.Sequence)
		return
	}
	// assigning sequence nr
	if seq, ok := r.fieldMappings[f.Name]; ok {
		f.Sequence = seq
		return
	}
	// find a free sequence nr
	for s := 1; s <= len(r.fieldNumbers); s++ {
		taken := false
		for _, each := range r.fieldNumbers {
			if each == s {
				taken = true
			}
		}
		if !taken {
			// it is now taken
			r.fieldNumbers = append(r.fieldNumbers, s)
			// assign it
			f.Sequence = s
			return
		}
	}

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
