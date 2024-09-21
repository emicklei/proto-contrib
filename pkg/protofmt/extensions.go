package protofmt

import "github.com/emicklei/proto"

// keyValuePair returns key = value or "value"
func keyValuePair(o *proto.Option, embedded bool) (cols []aligned) {
	equals := alignedEquals
	name := o.Name
	if len(o.Constant.OrderedMap) > 0 {
		cols = append(cols, leftAligned(name), equals)
		cols = append(cols, columnsPrintablesFromMap(o.Constant.OrderedMap)...)
		return
	}
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

func columnsPrintablesFromMap(m proto.LiteralMap) (cols []aligned) {
	cols = append(cols, leftAligned("{"), alignedSpace)
	for i, each := range m {
		// TODO only works for simple constants
		if i > 0 {
			cols = append(cols, alignedSpace)
		}
		cols = append(cols, leftAligned(each.Name), alignedColon, leftAligned(each.SourceRepresentation()))
	}
	cols = append(cols, alignedSpace, leftAligned("}"))
	return
}
