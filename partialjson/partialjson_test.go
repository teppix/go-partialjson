package partialjson

import (
	"bytes"
	"encoding/json"
	"testing"
)

// checkJSON is a helper to test json marshalling results
func checkJSON(t *testing.T, expected, result interface{}) {
	jsonExpected, err := json.Marshal(expected)
	if err != nil {
		t.Error("Json marshal error ", err)
		return
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		t.Error("Json marshal error ", err)
		return
	}

	if !bytes.Equal(jsonExpected, jsonResult) {
		t.Error(
			"Expected:", string(jsonExpected),
			"Got:", string(jsonResult),
		)
	}
}

func TestSimple(t *testing.T) {
	// empty struct
	type struct1 struct {
	}

	// full marshal
	checkJSON(t,
		&struct1{},
		Begin(&struct1{}),
	)

	// partial marshal
	checkJSON(t,
		&struct1{},
		Begin(&struct1{}).Partial(),
	)

	// struct with simple values
	type struct2 struct {
		A int
		B string
		C float64
		D bool
	}

	// test full encoding, empty values
	checkJSON(t,
		&struct2{},
		Begin(&struct2{}),
	)

	// test full encoding, initialized values
	checkJSON(t,
		&struct2{99, "foobar", 12.56, true},
		Begin(&struct2{}).
			Set("A", 99).
			Set("B", "foobar").
			Set("C", 12.56).
			Set("D", true),
	)

	// test partial encoding, all values set
	checkJSON(t,
		&struct2{99, "foobar", 12.5, true},
		Begin(&struct2{}).
			Set("A", 99).
			Set("B", "foobar").
			Set("C", 12.5).
			Set("D", true).
			Partial(),
	)

}

func TestPartial(t *testing.T) {
	type struct1a struct {
		A string
	}

	type struct1b struct {
		A string
		B string
	}

	// test partial encoding, only field A - using Set()
	checkJSON(t,
		&struct1a{"hello"},
		Begin(&struct1b{}).
			Set("A", "hello").
			Partial(),
	)

	// test partial encoding, only field A - using Use()
	checkJSON(t,
		&struct1a{"hello"},
		Begin(&struct1b{A: "hello"}).
			Use("A").
			Partial(),
	)
}

func TestComplexValues(t *testing.T) {
	type struct1 struct {
		A string
	}

	type struct2 struct {
		A struct1
		B []struct1
	}

	// test embedding lists and structs
	checkJSON(t,
		&struct2{
			struct1{"foo"},
			[]struct1{{"bar"}, {"baz"}},
		},

		Begin(&struct2{}).
			Partial().
			Set("A", struct1{"foo"}).
			Set("B", []struct1{{"bar"}, {"baz"}}))
}

func TestMultipleAssignments(t *testing.T) {
	type struct1 struct {
		A string
	}

	// test calling Set() twice for same key
	checkJSON(t,
		&struct1{"foo"},
		Begin(&struct1{}).
			Set("A", "bar").
			Set("A", "foo").
			Partial(),
	)

	// test calling Use() twice for same key
	checkJSON(t,
		&struct1{"foo"},
		Begin(&struct1{"foo"}).
			Use("A").
			Use("A").
			Partial(),
	)
}
