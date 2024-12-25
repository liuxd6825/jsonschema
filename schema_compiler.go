package jsonschema

import (
	"golang.org/x/text/message"
	"log"
	"strings"
)

func OrderKeysVocab() *Vocabulary {
	//url := "http://example.com/meta/unique-keys"
	url := "https://json-schema.org/draft/2020-12/vocab/order"
	schema, err := UnmarshalJSON(strings.NewReader(`{
		"properties": {
			"order": { "type": "number" },
			"name": { "type": "string" }
		}
	}`))
	if err != nil {
		log.Fatal(err)
	}

	c := NewCompiler()
	if err := c.AddResource(url, schema); err != nil {
		log.Fatal(err)
	}
	sch, err := c.Compile(url)
	if err != nil {
		log.Fatal(err)
	}

	return &Vocabulary{
		URL:     url,
		Schema:  sch,
		Compile: compileOrder,
	}
}

func compileOrder(ctx *CompilerContext, obj map[string]any) (SchemaExt, error) {
	v, ok := obj["order"]
	if !ok {
		return nil, nil
	}
	s, ok := v.(string)
	if !ok {
		return nil, nil
	}
	return &orderKeys{pname: s}, nil
}

type orderKeys struct {
	pname string
}

func (s *orderKeys) Validate(ctx *ValidatorContext, v any) {
	arr, ok := v.([]any)
	if !ok {
		return
	}
	var keys []any
	for _, item := range arr {
		obj, ok := item.(map[string]any)
		if !ok {
			continue
		}
		key, ok := obj[s.pname]
		if ok {
			keys = append(keys, key)
		}
	}

	i, j, err := ctx.Duplicates(keys)
	if err != nil {
		ctx.AddErr(err)
		return
	}
	if i != -1 {
		ctx.AddError(&OrderKeys{Key: s.pname, Duplicates: []int{i, j}})
	}
}

// ErrorKind --

type OrderKeys struct {
	Key        string
	Duplicates []int
}

func (*OrderKeys) KeywordPath() []string {
	return []string{"orderKeys"}
}

func (k *OrderKeys) LocalizedString(p *message.Printer) string {
	return p.Sprintf("order at %d and %d have same %s", k.Duplicates[0], k.Duplicates[1], k.Key)
}
