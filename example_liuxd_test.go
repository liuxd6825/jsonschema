package jsonschema_test

import (
	"bytes"
	"github.com/liuxd6825/jsonschema/v6"
	"log"
	"testing"
	"time"
)

func Test_DateType(t *testing.T) {
	compiler := jsonschema.NewCompiler()
	doc := `{
      "type": "object",
      "properties": {
        "commandId": {
          "type": "string"
        },
        "data": {
          "type": "object",
          "properties": {
			"name": {
			  "type": "string",
			  "minLength": 10
			},
			"birthday":{
			  "type": "date",
              "notNull": true
			}
 		  },
		  "required": ["name"]
		}
	  },
  	  "required": []
	}
	`
	reader, err := jsonschema.UnmarshalJSON(bytes.NewReader([]byte(doc)))
	if err != nil {
		log.Fatalf("UnmarshalJSON失败: %v", err)
	}

	schemaFile := "schema.json"
	if err := compiler.AddResource(schemaFile, reader); err != nil {
		log.Fatal("加载文件：", err)
	}
	sch, err := compiler.Compile(schemaFile)
	if err != nil {
		log.Fatalf("编译失败: %v", err)
	}
	inst := map[string]any{
		"commandId": "commandId",
		"data": map[string]any{
			"name":     "1111111111111",
			"birthday": time.Now(),
		},
	}
	err = sch.Validate(inst)
	if err != nil {
		t.Error(err)
	}
}

func Test_DateType2(t *testing.T) {
	compiler := jsonschema.NewCompiler()
	doc := `{
      "type": "object",
      "properties": {
		"birthday":{
		  "type": "date",
		  "notNull": true
		}
	  },
  	  "required": []
	}
	`
	reader, err := jsonschema.UnmarshalJSON(bytes.NewReader([]byte(doc)))
	if err != nil {
		log.Fatalf("UnmarshalJSON失败: %v", err)
	}

	schemaFile := "schema.json"
	if err := compiler.AddResource(schemaFile, reader); err != nil {
		log.Fatal("加载文件：", err)
	}
	sch, err := compiler.Compile(schemaFile)
	if err != nil {
		log.Fatalf("编译失败: %v", err)
	}
	inst := map[string]any{
		"birthday": "time.Now()",
	}
	err = sch.Validate(inst)
	if err != nil {
		t.Error(err)
	}
}

func Test_OneOf(t *testing.T) {
	compiler := jsonschema.NewCompiler()
	doc := `
	{
      "type": "object",
      "properties": {
		"birthday":{
		  "type": ["date","null"]
		}
	  }
	}
	`
	reader, err := jsonschema.UnmarshalJSON(bytes.NewReader([]byte(doc)))
	if err != nil {
		log.Fatalf("UnmarshalJSON失败: %v", err)
	}

	schemaFile := "schema.json"
	if err := compiler.AddResource(schemaFile, reader); err != nil {
		log.Fatal("加载文件：", err)
	}
	sch, err := compiler.Compile(schemaFile)
	if err != nil {
		log.Fatalf("编译失败: %v", err)
	}
	inst := map[string]any{
		"birthday": nil,
	}
	err = sch.Validate(inst)
	if err != nil {
		t.Error(err)
	}
}
