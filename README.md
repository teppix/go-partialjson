# PartialJSON

Selective marshalling of structs to JSON in Go.

## Overview

PartialJSON makes it possible to marshal structs to JSON, while selecting which fields to include in the output.

The type check occurs at runtime - the idea is to fail early instead of producing incorrect json data.

## Notes about type safety

When using `Set()` to assign values, the library will _panic_ if an incorrect type is encountered at the moment.

To ensure type safety at compile time, assign values to the struct as usual, and flag them for inclusion by calling `Use()`

## Usage

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/teppix/go-partialjson/partialjson"
)

type MyStruct struct {
	A int
	B string
	C float64
}

func main() {
	var data []byte
	var val *partialjson.Builder

	// Marshal only assigned values
	// Will cause a panic if wrong type is supplied to Set()
	val = partialjson.Begin(&MyStruct{}).
		Set("A", 12).
		Set("B", "testing").
		Partial()

	data, _ = json.Marshal(val)
	fmt.Println(string(data))

	// output: {"A":12,"B":"testing"}

	// Marshal complete struct
	// Will cause a panic if wrong type is supplied to Set()
	val = partialjson.Begin(&MyStruct{}).
		Set("A", 12).
		Set("B", "testing")

	data, _ = json.Marshal(val)
	fmt.Println(string(data))

	// output: {"A":12,"B":"testing","C":0}

	// Marshal predefined struct fields
	// Should be type safe and panic-free
	val = partialjson.Begin(
		&MyStruct{
			B: "testing",
		}).
		Use("B").
		Partial()

	data, _ = json.Marshal(val)
	fmt.Println(string(data))

	// output: {"B":"testing"}
}
```
