package protodecode

import (
	"bytes"
	"fmt"
	"io/ioutil"

	pp "github.com/emicklei/proto"
)

type Definitions struct {
	specs         map[string]*pp.Message
	filenamesRead []string
}

func NewDefinitions() *Definitions {
	return &Definitions{
		specs:         map[string]*pp.Message{},
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
	return nil
}

func (d *Definitions) Message(pkg string, name string) *pp.Message {
	key := fmt.Sprintf("%s.%s", pkg, name)
	return d.specs[key]
}

func (d *Definitions) AddMessage(pkg string, name string, message *pp.Message) {
	key := fmt.Sprintf("%s.%s", pkg, name)
	d.specs[key] = message
}

func packageOf(def *pp.Proto) string {
	for _, each := range def.Elements {
		if p, ok := each.(*pp.Package); ok {
			return p.Name
		}
	}
	return ""
}
