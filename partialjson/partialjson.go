// PartialJSON - selective marshalling of structs to JSON.
package partialjson

// TODO - allow embedding Builder (partial) as struct field
// TODO - check for duplicate fields in Set() and Use()

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

// Set marks a field for inclusion and updates its value.
func (b *Builder) Set(name string, val interface{}) *Builder {
	container := reflect.Indirect(reflect.ValueOf(b.base))
	field := container.FieldByName(name)

	if !field.IsValid() {
		b.badFields = append(b.badFields, name)
		return b
	}

	field.Set(reflect.ValueOf(val))
	b.keys = append(b.keys, name)
	return b
}

// Use marks a field for inclusion without changing its value.
func (b *Builder) Use(name string) *Builder {
	b.keys = append(b.keys, name)
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
