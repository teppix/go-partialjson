// PartialJSON - selective marshalling of structs to JSON.
package partialjson

// TODO - allow embedding Builder (partial) as struct field

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

type Builder struct {
	keys      []string
	base      interface{}
	partial   bool
	badFields []string
}

func (b *Builder) Partial() *Builder {
	b.partial = true
	return b
}

// IsSet returns true if key is set
func (b *Builder) IsSet(name string) bool {
	return b.indexOf(name) != -1
}

// indexOf returns index of key.
// -1 is returned if there is no match.
func (b *Builder) indexOf(name string) int {
	for i, key := range b.keys {
		if key == name {
			return i
		}
	}
	return -1
}

// Remove removes the key from the builder.
//
// The actual field value of the underlying struct
// remains unchanged.
func (b *Builder) Remove(name string) *Builder {
	index := b.indexOf(name)
	if index == -1 {
		return b
	}

	b.keys = append(b.keys[:index], b.keys[index+1:]...)
	return b
}

// Set marks a field for inclusion and updates its value.
// - Each key is only added once
// - Every call using same key replaces the previous value
func (b *Builder) Set(name string, val interface{}) *Builder {
	container := reflect.Indirect(reflect.ValueOf(b.base))
	field := container.FieldByName(name)

	if !field.IsValid() {
		b.badFields = append(b.badFields, name)
		return b
	}

	field.Set(reflect.ValueOf(val))
	b.Use(name)
	return b
}

// Use marks a field for inclusion without changing its value.
// Each key is only added once
func (b *Builder) Use(name string) *Builder {
	if !b.IsSet(name) {
		b.keys = append(b.keys, name)
	}
	return b
}

func (b *Builder) MarshalJSON() ([]byte, error) {
	var err error
	var tmp []byte
	container := reflect.Indirect(reflect.ValueOf(b.base))

	if len(b.badFields) > 0 {
		return nil, fmt.Errorf("Invalid struct fields: %v", b.badFields)
	}

	if b.partial {
		var buf bytes.Buffer
		first := true

		buf.WriteByte('{')
		for _, key := range b.keys {
			if first {
				first = false
			} else {
				buf.WriteByte(',')
			}

			// write key
			tmp, err = json.Marshal(key)
			if err != nil {
				return nil, err
			}
			buf.Write(tmp)

			// separator
			buf.WriteByte(':')

			// write value
			field := container.FieldByName(key)
			tmp, err = json.Marshal(field.Interface())
			if err != nil {
				return nil, err
			}
			buf.Write(tmp)
		}
		buf.WriteByte('}')
		return buf.Bytes(), nil
	} else {
		return json.Marshal(b.base)
	}
}

// Begin is the starting point for creating jsonbuilders.
// base is an instance of the struct you want to encode.
func Begin(base interface{}) *Builder {
	return &Builder{
		keys:      []string{},
		base:      base,
		badFields: []string{},
	}
}
