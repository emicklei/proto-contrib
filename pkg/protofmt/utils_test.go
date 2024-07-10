package protofmt

import (
	"bytes"

	"github.com/emicklei/proto"
)

// / testing utils
func formatted(v proto.Visitee) string {
	b := new(bytes.Buffer)
	f := NewFormatter(b, WithIndentSeparator("  ")) // 2 spaces
	v.Accept(f)
	return b.String()
}

func diff(left, right string) string {
	b := new(bytes.Buffer)
	w := func(char rune) {
		if '\n' == char {
			b.WriteString("(n)")
		} else if '\t' == char {
			b.WriteString("(t)")
		} else if ' ' == char {
			b.WriteString("(.)")
		} else {
			b.WriteRune(char)
		}
	}
	b.WriteString("got:\n")
	for _, char := range left {
		w(char)
	}
	if len(left) == 0 {
		b.WriteString("(empty)")
	}
	b.WriteString("\n")
	for _, char := range right {
		w(char)
	}
	b.WriteString("\n:wanted\n")
	return b.String()
}
