package jsonschema

import (
	"encoding/json"
	"github.com/invopop/jsonschema"
)

func ReflectJson(v any) ([]byte, error) {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties:  false, // 禁止额外字段
		RequiredFromJSONSchemaTags: true,  // 从 JSON Schema 标签解析 required
	}
	schema := reflector.Reflect(v)
	return json.MarshalIndent(schema, "", "  ")
}
