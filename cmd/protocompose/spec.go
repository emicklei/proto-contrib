package main

import "strings"

type composeSpec struct {
	registryKey string
	fieldName   string
}

// pre: commentLine ends with [package].[field]
func newComposeSpec(commentLine string) composeSpec {
	lineParts := strings.Split(commentLine, " ")
	composePath := lineParts[len(lineParts)-1]
	parts := strings.Split(composePath, ".")
	fullType := strings.Join(parts[0:len(parts)-1], ".")
	fieldName := parts[len(parts)-1]
	return composeSpec{
		registryKey: fullType,
		fieldName:   fieldName,
	}
}
