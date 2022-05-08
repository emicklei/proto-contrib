package main

import "strings"

type composeSpec struct {
	registryKey  string
	fieldName    string
	inlineFields bool
	embedMessage bool
}

// pre: commentLine ends with:
// 1 [package].[type].[field]
// 2 ..[package].[type]
// 3 #[package].[type]
func newComposeSpec(commentLine string) composeSpec {
	lineParts := strings.Split(commentLine, " ")
	composePath := lineParts[len(lineParts)-1]
	if strings.HasPrefix(composePath, "..") {
		// inline fields
		fullType := composePath[2:]
		return composeSpec{
			registryKey:  fullType,
			inlineFields: true,
			embedMessage: false,
		}
	}
	if strings.HasPrefix(composePath, "#") {
		// embed message
		fullType := composePath[1:]
		return composeSpec{
			registryKey:  fullType,
			embedMessage: true,
		}
	}
	// normal field
	parts := strings.Split(composePath, ".")
	fullType := strings.Join(parts[0:len(parts)-1], ".")
	fieldName := parts[len(parts)-1]
	return composeSpec{
		registryKey:  fullType,
		fieldName:    fieldName,
		inlineFields: false,
		embedMessage: false,
	}
}

func (c composeSpec) forFieldName(name string) composeSpec {
	return composeSpec{
		registryKey:  c.registryKey,
		fieldName:    name,
		inlineFields: false,
		embedMessage: false,
	}
}
