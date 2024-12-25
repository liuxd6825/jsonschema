package jsonschema

import (
	"bytes"
	"os"
	"testing"
)

func TestSchema_GetSortProperties(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path + "/testdata/order.json")
	if err != nil {
		t.Fatal(err)
	}

	reader, err := UnmarshalJSON(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	schemaFile := "schema.json"
	compiler := NewCompiler()
	if err := compiler.AddResource(schemaFile, reader); err != nil {
		panic(err)
	}

	sch, err := compiler.Compile(schemaFile)
	if err != nil {
		panic(err)
	}

	props := sch.GetSortProperties()
	t.Log(props)
}
