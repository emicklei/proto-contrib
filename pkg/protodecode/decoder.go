package protodecode

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"

	pp "github.com/emicklei/proto"
	pb "github.com/golang/protobuf/proto"
)

var EOM = errors.New("end of message reached")

type decoder struct {
	d       *Definitions
	m       *pp.Message
	p       string
	b       *pb.Buffer
	r       map[string]interface{}
	verbose bool
}

func (d *decoder) Decode(pkg, t string) (map[string]interface{}, error) {
	m := d.d.Message(pkg, t)
	d.p = pkg
	if m == nil {
		return nil, fmt.Errorf("no definition found for %s %s", pkg, t)
	}
	d.m = m
	for {
		op, err := d.b.DecodeVarint()
		if err != nil {
			break
		}
		tag := op >> 3
		wire := op & 7
		if err := d.decodeTag(tag, wire); err != nil {
			return d.r, err
		}
	}
	return d.r, EOM
}

// NewDecoder is the thing
func NewDecoder(d *Definitions, b *pb.Buffer) *decoder {
	return &decoder{
		d: d,
		b: b,
		r: map[string]interface{}{},
	}
}

func (d *decoder) decodeTag(tag, wire uint64) error {
	for _, each := range d.m.Elements {
		if f, ok := each.(*pp.NormalField); ok {
			if f.Sequence == int(tag) {
				return d.decodeNormalField(f, wire)
			}
		}
		if f, ok := each.(*pp.MapField); ok {
			if f.Sequence == int(tag) {
				return d.decodeMapField(f, wire)
			}
		}
		if f, ok := each.(*pp.OneOfField); ok {
			if f.Sequence == int(tag) {
				return d.decodeOneOfField(f, wire)
			}
		}
	}
	return nil
}

func (d *decoder) decodeNormalField(f *pp.NormalField, wire uint64) error {
	if "string" == f.Type {
		return d.handleString(f.Name, f.Repeated)
	}
	if "int64" == f.Type {
		return d.handleInt64(f.Name, f.Repeated)
	}
	if "int32" == f.Type {
		return d.handleInt32(f.Name, f.Repeated)
	}
	if "float" == f.Type {
		return d.handleFloat(f.Name, f.Repeated)
	}
	if "bool" == f.Type {
		return d.handleBool(f.Name, f.Repeated)
	}
	if m := d.d.Message(d.p, f.Type); m == nil {
		return fmt.Errorf("unknown type:%s", f.Type)
	}
	nextData, err := d.b.DecodeRawBytes(true)
	if err != nil {
		if io.ErrUnexpectedEOF == err {
			return err
		}
		return fmt.Errorf("unable or read raw bytes of message of type:%s", f.Type)
	}
	sub := NewDecoder(d.d, pb.NewBuffer(nextData))
	if f.Repeated {
		for {
			if d.verbose {
				log.Println("BEGIN repeated", f.Name, ":", f.Type)
			}
			if _, err := sub.Decode(d.p, f.Type); err != nil { // what package?
				if io.ErrUnexpectedEOF == err {
					if d.verbose {
						log.Println("END", f.Name, ":", f.Type)
					}
					// TODO?????
					d.add(f.Name, sub.r, RepeatedField, !MapField)
					break
				}
				if EOM == err {
					// TODO?????
					d.add(f.Name, sub.r, RepeatedField, !MapField)
					break
				}
				return fmt.Errorf("unable to decode repeated message of type:%v error:%v", f.Type, err)
			}
			d.add(f.Name, sub.r, RepeatedField, !MapField)
		}
	} else {
		// single
		if d.verbose {
			log.Println("BEGIN single", f.Name, ":", f.Type)
		}
		if _, err := sub.Decode(d.p, f.Type); err != nil { // what package
			if io.ErrUnexpectedEOF == err || EOM == err {
				if d.verbose {
					log.Println("END", f.Name, ":", f.Type)
				}
				d.add(f.Name, sub.r, !RepeatedField, !MapField)
			} else {
				return fmt.Errorf("unable to decode single message of type:%v error:%v", f.Type, err)
			}
		}
	}
	return nil
}

const (
	RepeatedField = true
	MapField      = true
)

func (d *decoder) add(key string, value interface{}, repeated bool, isMap bool) {
	if d.verbose {
		log.Printf("[%s] add [%s=%v] repeated:%v map:%v\n", d.m.Name, key, value, repeated, isMap)
	}
	if repeated {
		if val, ok := d.r[key]; ok {
			maps := val.([]interface{})
			maps = append(maps, value)
			d.r[key] = maps
		} else {
			d.r[key] = []interface{}{value}
		}
	} else if isMap {
		if val, ok := d.r[key]; ok {
			// map exists
			outMap := val.(map[string]interface{}) // TODO key can be any type
			inMap := value.(map[string]interface{})
			for k, v := range inMap {
				outMap[k] = v
			}
			// needed?
			d.r[key] = outMap
		} else {
			// map did not exist
			outMap := map[string]interface{}{}
			inMap := value.(map[string]interface{})
			for k, v := range inMap {
				outMap[k] = v
			}
			d.r[key] = outMap
		}
	} else {
		d.r[key] = value
	}
}

func (d *decoder) decodeOneOfField(f *pp.OneOfField, wire uint64) error {
	// TODO
	log.Println("WARN:unhandled oneof field")
	return nil
}

// https://developers.google.com/protocol-buffers/docs/proto3#maps
func (d *decoder) decodeMapField(f *pp.MapField, wire uint64) error {
	// create temporary proto Message such that we can use another decoder to do all the work
	entryMessage := new(pp.Message)
	entryMessage.Name = d.m.Name + "." + f.Name + ".Entry"
	entryMessage.Elements = []pp.Visitee{
		&pp.NormalField{
			Field: &pp.Field{
				Name:     "key",
				Type:     f.KeyType,
				Sequence: 1,
			}},
		&pp.NormalField{
			Field: &pp.Field{Name: "value",
				Type:     f.Type,
				Sequence: 2,
			},
		}}
	containerMessage := new(pp.Message)
	containerMessage.Name = d.m.Name + "." + f.Name + ".List"
	containerMessage.Elements = append(entryMessage.Elements,
		&pp.NormalField{
			Repeated: true,
			Field: &pp.Field{
				Name:     f.Name,
				Type:     entryMessage.Name,
				Sequence: 1}})
	nextData, err := d.b.DecodeRawBytes(true)
	if err != nil {
		if io.ErrUnexpectedEOF == err {
			return err
		}
		return fmt.Errorf("unable or read raw bytes of map of type:%s", f.Type)
	}
	// extend the specs TODO check for existence
	d.d.AddMessage(d.p, containerMessage.Name, containerMessage)
	d.d.AddMessage(d.p, entryMessage.Name, entryMessage)

	sub := NewDecoder(d.d, pb.NewBuffer(nextData))
	sub.Decode(d.p, containerMessage.Name) // what package?

	// one of the repeated
	result := map[string]interface{}{}
	result[sub.r["key"].(string)] = sub.r["value"]
	d.add(f.Name, result, !RepeatedField, MapField)

	return nil
}

func (d *decoder) handleInt64(n string, repeated bool) error {
	if repeated {
		data, err := d.b.DecodeRawBytes(true)
		if err != nil {
			return fmt.Errorf("cannot decode repeated int64 raw bytes:%v", err)
		}
		buf := pb.NewBuffer(data)
		for {
			x, err := buf.DecodeVarint()
			if err == io.ErrUnexpectedEOF {
				break
			}
			d.add(n, x, repeated, !MapField)
		}
		return nil
	}
	// non-repeated
	x, err := d.b.DecodeVarint()
	if err != nil {
		if io.ErrUnexpectedEOF == err {
			return err
		}
		return fmt.Errorf("cannot decode %s:int64:%v", n, err)
	}
	d.add(n, x, !RepeatedField, !MapField)
	return nil
}

func (d *decoder) handleInt32(n string, repeated bool) error {
	if repeated {
		data, err := d.b.DecodeRawBytes(true)
		if err != nil {
			return fmt.Errorf("cannot decode repeated int32 raw bytes:%v", err)
		}
		buf := pb.NewBuffer(data)
		for {
			x, err := buf.DecodeVarint()
			if err == io.ErrUnexpectedEOF {
				break
			}
			d.add(n, int32(x), repeated, !MapField)
		}
		return nil
	}
	// non-repeated
	x, err := d.b.DecodeVarint()
	if err != nil {
		if io.ErrUnexpectedEOF == err {
			return err
		}
		return fmt.Errorf("cannot decode %s:int32:%v", n, err)
	}
	d.add(n, x, !RepeatedField, !MapField)
	return nil
}

func (d *decoder) handleFloat(n string, repeated bool) error {
	if repeated {
		data, err := d.b.DecodeRawBytes(true)
		if err != nil {
			return fmt.Errorf("cannot decode repeated float32 raw bytes:%v", err)
		}
		buf := pb.NewBuffer(data)
		for {
			x, err := buf.DecodeFixed32()
			if err == io.ErrUnexpectedEOF {
				break
			}
			d.add(n, math.Float32frombits(uint32(x)), repeated, !MapField)
		}
		return nil
	}
	// non-repeated
	x, err := d.b.DecodeFixed32()
	if err != nil {
		if io.ErrUnexpectedEOF == err {
			return err
		}
		return fmt.Errorf("cannot decode %s:float:%v", n, err)
	}
	d.add(n, math.Float32frombits(uint32(x)), !RepeatedField, !MapField)
	return nil
}

func (d *decoder) handleString(n string, repeated bool) error {
	// non-repeated and repeated
	sb, err := d.b.DecodeStringBytes()
	if err != nil {
		if io.ErrUnexpectedEOF == err {
			return err
		}
		return fmt.Errorf("cannot decode %s:string:%v", n, err)
	}
	d.add(n, string(sb), repeated, !MapField)
	return nil
}

func (d *decoder) handleBool(n string, repeated bool) error {
	if repeated {
		data, err := d.b.DecodeRawBytes(true)
		if err != nil {
			return fmt.Errorf("cannot decode repeated bool raw bytes:%v", err)
		}
		buf := pb.NewBuffer(data)
		for {
			x, err := buf.DecodeVarint()
			if err == io.ErrUnexpectedEOF {
				break
			}
			d.add(n, x == 1, repeated, !MapField)
		}
		return nil
	}
	// non-repeated
	x, err := d.b.DecodeVarint()
	if err != nil {
		if io.ErrUnexpectedEOF == err {
			return err
		}
		return fmt.Errorf("cannot decode %s:bool:%v", n, err)
	}
	d.add(n, x == 1, repeated, !MapField)
	return nil
}
