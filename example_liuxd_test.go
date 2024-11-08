package jsonschema_test

import (
	"github.com/liuxd6825/jsonschema/v6"
	"log"
	"testing"
)

func Test_DateType(t *testing.T) {
	doc := `
		"type": "object",
		"properties": {
			"birthday": {
				"type": "string",
				"format": "date-time",
			}
		}
	`
	schemaFile := "./testdata/examples/dateType.json"
	c := jsonschema.NewCompiler()
	if err := c.AddResource("schema.json", doc); err != nil {
		log.Fatal(err)
	}

	sch, err := c.Compile(schemaFile)
	if err != nil {
		log.Fatalf("failed to compile schema: %v", err)
	}

	inst := map[string]any{
		"name": "11111",
	}

	err = sch.Validate(inst)
	if err != nil {
		t.Error(err)
	}

}
