package proto2gql

import (
	"path"
	"path/filepath"
	"strings"
)

type (
	Type struct {
		originalPackageName  string
		originalName         string
		convertedPackageName string
		convertedName        string
	}

	Scope struct {
		converter            *Converter
		originalPackageName  string
		convertedPackageName string
		path                 []string
		types                map[string]*Type
		imports              map[string]*Type
		children             map[string]*Scope
	}
)

func NewScope(converter *Converter) *Scope {
	return &Scope{
		converter: converter,
		path:      make([]string, 0, 5),
		types:     make(map[string]*Type),
		imports:   make(map[string]*Type),
		children:  make(map[string]*Scope),
	}
}

func (s *Scope) Fork(name string) *Scope {
	p := make([]string, len(s.path))

	copy(p, s.path)

	childScope := &Scope{
		converter:            s.converter,
		originalPackageName:  s.originalPackageName,
		convertedPackageName: s.convertedPackageName,
		types:                s.types, // share types collection
		path:                 append(p, name),
		children:             make(map[string]*Scope),
	}

	s.children[name] = childScope

	return childScope
}

func (s *Scope) SetPackageName(name string) {
	s.originalPackageName = name
	s.convertedPackageName = s.converter.PackageName(strings.Split(name, "."))
}

func (s *Scope) AddLocalType(name string) {
	typeName := s.converter.OriginalTypeName(s, name)

	_, ok := s.types[typeName]

	if ok == false {
		s.types[typeName] = &Type{
			originalPackageName:  s.originalPackageName,
			originalName:         name,
			convertedPackageName: s.convertedPackageName,
			convertedName:        s.converter.NewTypeName(s, name),
		}
	}
}

func (s *Scope) AddImportedType(filename string) {
	dir := path.Dir(filename)
	separator := string(filepath.Separator)
	name := strings.ToUpper(strings.Replace(path.Base(filename), ".proto", "", -1))
	ref := strings.Replace(dir, separator, ".", -1)

	var originalPackageName string
	var convertedPackageName string

	if dir == "." {
		originalPackageName = s.originalPackageName
		convertedPackageName = s.convertedPackageName
	} else {
		originalPackageName = strings.Replace(dir, separator, ".", -1)
		convertedPackageName = s.converter.PackageName(strings.Split(dir, separator))
	}

	s.imports[ref] = &Type{
		originalPackageName:  originalPackageName,
		originalName:         name,
		convertedPackageName: convertedPackageName,
		convertedName:        name,
	}
}

func (s *Scope) ResolveConvertedTypeName(ref string) string {
	builtin, ok := BUILTINS[ref]

	if ok == true {
		return builtin
	}

	// try to find one in a global scope
	scoped, ok := s.types[ref]

	if ok == true {
		return scoped.convertedName
	}

	// try to find one among nested types
	nested, ok := s.types[s.converter.OriginalTypeName(s, ref)]

	if ok == true {
		return nested.convertedName
	}

	var foundInChildren string

	for _, childScope := range s.children {
		res := childScope.ResolveConvertedTypeName(ref)

		if res != ref {
			foundInChildren = res
			break
		}
	}

	if foundInChildren != "" {
		return foundInChildren
	}

	// if we are still here, probably it's an imported type

	// if type does not contain "." it means it's from the same package
	if strings.Contains(ref, ".") == false {
		// from the same package
		imported, ok := s.imports["."]

		if ok == true {
			return imported.convertedPackageName + ref
		}
	} else {
		// if it has "." it means it's from other package
		parts := strings.Split(ref, ".")

		var tail string

		for idx, segment := range parts {
			if tail == "" {
				tail = segment
			} else {
				tail += "." + segment
			}

			imported, ok := s.imports[tail]

			if ok == true {
				return imported.convertedPackageName + strings.Join(parts[idx+1:], "")
			}
		}
	}

	return ref
}

func (s *Scope) ResolveFullTypeName(ref string) string {
	builtin, ok := BUILTINS[ref]

	if ok == true {
		return builtin
	}

	// try to find one in a global scope
	scoped, ok := s.types[ref]

	if ok == true {
		return s.converter.OriginalFullTypeName(s, scoped.originalName)
	}

	// try to find one among nested types
	nested, ok := s.types[s.converter.OriginalTypeName(s, ref)]

	if ok == true {
		return s.converter.OriginalFullTypeName(s, nested.originalName)
	}

	var foundInChildren string

	for _, childScope := range s.children {
		res := childScope.ResolveFullTypeName(ref)

		if res != ref {
			foundInChildren = res
			break
		}
	}

	if foundInChildren != "" {
		return foundInChildren
	}

	// if we are still here, probably it's an imported type

	// if type does not contain "." it means it's from the same package
	if strings.Contains(ref, ".") == false {
		// from the same package
		imported, ok := s.imports["."]

		if ok == true {
			return imported.originalPackageName + "." + imported.originalName
		}
	}

	return ref
}
