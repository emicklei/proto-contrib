package protofmt

import "github.com/emicklei/proto"

// keyValuePair returns key = value or "value"
func keyValuePair(o *proto.Option, embedded bool) (cols []aligned) {
	equals := alignedEquals
	name := o.Name
	if embedded {
		return append(cols, leftAligned(name), equals, leftAligned(o.Constant.SourceRepresentation())) // numbers right, strings left? TODO
	}
	return append(cols, rightAligned(name), equals, rightAligned(o.Constant.SourceRepresentation()))
}

func alignedInlinePrefix(c *proto.Comment) aligned {
	prefix := " //"
	if c.ExtraSlash {
		prefix = " ///"
	}
	return notAligned(prefix)
}

func columnsPrintables(c *proto.Comment) (list []columnsPrintable) {
	for _, each := range c.Lines {
		list = append(list, inlineComment{each, c.ExtraSlash})
	}
	return
}

func typeAssertColumnsPrintable(v proto.Visitee) (columnsPrintable, bool) {
	return asColumnsPrintable(v), len(asColumnsPrintable(v).columns()) > 0
}
