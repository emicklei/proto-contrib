package main

import "strings"

type composeSpec struct {
	registryKey string
	fieldName   string
}

// pre: commentLine ends with:
// 1 [package].[type].[field]
// 2 ...[package].[type]
// 3 #[package].[type]
func newComposeSpec(commentLine string) composeSpec {
	trimmed := strings.TrimSpace(commentLine)
	lineParts := strings.Split(trimmed, " ")
	composePath := lineParts[len(lineParts)-1]
	// normal field
	parts := strings.Split(composePath, ".")
	fullType := strings.Join(parts[0:len(parts)-1], ".")
	fieldName := parts[len(parts)-1]
	return composeSpec{
		registryKey: fullType,
		fieldName:   fieldName,
	}
}

func (c composeSpec) forFieldName(name string) composeSpec {
	return composeSpec{
		registryKey: c.registryKey,
		fieldName:   name,
	}
}
