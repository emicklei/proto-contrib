package protofmt

type inlineComment struct {
	line       string
	extraSlash bool
}

func (i inlineComment) columns() (list []aligned) {
	if len(i.line) == 0 {
		return append(list, notAligned(""))
	}
	prefix := "//"
	if i.extraSlash {
		prefix = "///"
	}
	return append(list, notAligned(prefix+i.line))
}
