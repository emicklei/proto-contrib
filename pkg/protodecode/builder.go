package protodecode

type WireEvent struct {
	Key        string
	Value      interface{}
	IsRepeated bool
	IsMap      bool
}

type Handler interface {
	Handle(WireEvent)
}

type MapBuilder struct {
	m map[string]interface{}
}

// NewMapBuilder return a builder for constructing a generic map string->interface{}
func NewMapBuilder() MapBuilder {
	return MapBuilder{
		m: map[string]interface{}{},
	}
}

func (b MapBuilder) Build() map[string]interface{} {
	return b.m
}

// Handle implements Handler
func (b MapBuilder) Handle(e WireEvent) {
	if e.IsRepeated {
		if val, ok := b.m[e.Key]; ok {
			maps := val.([]interface{})
			maps = append(maps, e.Value)
			b.m[e.Key] = maps
		} else {
			b.m[e.Key] = []interface{}{e.Value}
		}
	} else if e.IsMap {
		if val, ok := b.m[e.Key]; ok {
			// map exists
			outMap := val.(map[string]interface{}) // TODO key can be any type
			inMap := e.Value.(map[string]interface{})
			for k, v := range inMap {
				outMap[k] = v
			}
			// needed?
			b.m[e.Key] = outMap
		} else {
			// map did not exist
			outMap := map[string]interface{}{}
			inMap := e.Value.(map[string]interface{})
			for k, v := range inMap {
				outMap[k] = v
			}
			b.m[e.Key] = outMap
		}
	} else {
		b.m[e.Key] = e.Value
	}
}
