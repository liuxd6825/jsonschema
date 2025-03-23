package extensions

import (
	"encoding/json"
	"github.com/liuxd6825/jsonschema/v6"
	"log"
	"strings"
)

type DBOrder string

const (
	DBOrder_None DBOrder = ""
	DBOrder_Asc  DBOrder = "asc"
	DBOrder_Desc DBOrder = "desc"
)

type DBExtension struct {
	FieldName  string
	Size       int64
	PrimaryKey bool
	NotNull    bool
	Unique     bool
	Updatable  bool
	Readable   bool
	Order      DBOrder
}

func (d *DBExtension) Validate(ctx *jsonschema.ValidatorContext, v any) {

}

func NewDBVocabulary() *jsonschema.Vocabulary {
	url := "http://example.com/meta/db"
	schema, err := jsonschema.UnmarshalJSON(strings.NewReader(`{
		"properties" : {
			"db": {
				"type": "object",
				"patternProperties": {
					"fieldName": {
						"type": "string"
					},
					"size": {
						"type": "integer"
					},
					"primaryKey":{
						"type": "boolean"
					},
					"notNull": {
						"type": "boolean"
					},
					"unique": {
						"type": "boolean"
					},
					"updatable": {
						"type": "boolean",
						"default": true
					},
					"readable": {
						"type": "boolean",
						"default": true
					},
					"order": {
						"type": "string"
					}
				}
			}
		}
	}`))
	if err != nil {
		log.Fatal(err)
	}

	c := jsonschema.NewCompiler()
	if err := c.AddResource(url, schema); err != nil {
		log.Fatal(err)
	}
	sch, err := c.Compile(url)
	if err != nil {
		log.Fatal(err)
	}

	return &jsonschema.Vocabulary{
		URL:    url,
		Schema: sch,
		Subschemas: []jsonschema.SchemaPath{
			{jsonschema.Prop("db"), jsonschema.AllProp{}, jsonschema.AllProp{}},
		},
		Compile: compileDiscriminator,
	}
}

func compileDiscriminator(ctx *jsonschema.CompilerContext, obj map[string]any) (jsonschema.SchemaExt, error) {
	v, ok := obj["db"]
	if !ok {
		return nil, nil
	}
	name := obj["name"].(string)

	d, ok := v.(map[string]any)
	if !ok {
		return nil, nil
	}
	var err error
	db := &DBExtension{
		FieldName: name,
		Size:      100,
	}
	for key, value := range d {
		switch key {
		case "primaryKey":
			if val, ok := value.(bool); ok {
				db.PrimaryKey = val
			}
		case "size":
			if val, ok := value.(json.Number); ok {
				db.Size, err = val.Int64()
			}
		case "fieldName":
			if val, ok := value.(string); ok {
				db.FieldName = val
			}
		case "notNull":
			if val, ok := value.(bool); ok {
				db.NotNull = val
			}
		case "unique":
			if val, ok := value.(bool); ok {
				db.Unique = val
			}
		case "updatable":
			if val, ok := value.(bool); ok {
				db.Updatable = val
			}
		case "readable":
			if val, ok := value.(bool); ok {
				db.Readable = val
			}
		case "order":
			if val, ok := value.(string); ok {
				db.Order = DBOrder(val)
			}
		}
	}
	return db, err
}
