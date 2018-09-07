package protodecode

import (
	"bytes"
	"fmt"
	"io/ioutil"

	pp "github.com/emicklei/proto"
)

type Definitions struct {
	messages      map[string]*pp.Message
	enums         map[string]*pp.Enum
	filenamesRead []string
}

func NewDefinitions() *Definitions {
	return &Definitions{
		messages:      map[string]*pp.Message{},
		enums:         map[string]*pp.Enum{},
		filenamesRead: []string{},
	}
}

// Read the proto definition from a filename.
// Recursively add all imports
func (d *Definitions) AddFromFile(filename string) error {
	for _, each := range d.filenamesRead {
		if each == filename {
			return nil
		}
	}
	d.filenamesRead = append(d.filenamesRead, filename)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	parser := pp.NewParser(bytes.NewReader(data))
	def, err := parser.Parse()
	if err != nil {
		return err
	}
	pkg := packageOf(def)
	pp.Walk(def, pp.WithMessage(func(each *pp.Message) {
		d.AddMessage(pkg, each.Name, each)
	}))
	pp.Walk(def, pp.WithEnum(func(each *pp.Enum) {
		d.AddEnum(pkg, each.Name, each)
	}))
	return nil
}

func (d *Definitions) Message(pkg string, name string) (m *pp.Message, ok bool) {
	key := fmt.Sprintf("%s.%s", pkg, name)
	m, ok = d.messages[key]
	return
}

func (d *Definitions) Enum(pkg string, name string) (e *pp.Enum, ok bool) {
	key := fmt.Sprintf("%s.%s", pkg, name)
	e, ok = d.enums[key]
	return
}

// AddEnum adds the Enum
func (d *Definitions) AddEnum(pkg string, name string, enu *pp.Enum) {
	key := fmt.Sprintf("%s.%s", pkg, name)
	d.enums[key] = enu
}

func (d *Definitions) AddMessage(pkg string, name string, message *pp.Message) {
	key := fmt.Sprintf("%s.%s", pkg, name)
	d.messages[key] = message
}

func packageOf(def *pp.Proto) string {
	for _, each := range def.Elements {
		if p, ok := each.(*pp.Package); ok {
			return p.Name
		}
	}
	return ""
}
