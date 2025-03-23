package extensions

import (
	"bytes"
	"github.com/liuxd6825/jsonschema/v6"
	"os"
	"testing"
)

func TestSchema_DBExtentsion(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path + "/testdata/db.json")
	if err != nil {
		t.Fatal(err)
	}

	reader, err := jsonschema.UnmarshalJSON(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	schemaFile := "schema.json"
	compiler := jsonschema.NewCompiler()
	compiler.AssertVocabs()
	compiler.RegisterVocabulary(NewDBVocabulary())

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
