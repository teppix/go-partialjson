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
