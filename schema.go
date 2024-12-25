package jsonschema

import (
	"encoding/json"
	"fmt"
	"math/big"
	"time"
)

// Schema is the regpresentation of a compiled
// jsonschema.
type Schema struct {
	up                urlPtr             `json:"up"`
	resource          *Schema            `json:"resource,omitempty"`
	dynamicAnchors    map[string]*Schema `json:"dynamicAnchors,omitempty"`
	allPropsEvaluated bool               `json:"allPropsEvaluated,omitempty"`
	allItemsEvaluated bool               `json:"allItemsEvaluated,omitempty"`
	numItemsEvaluated int                `json:"numItemsEvaluated,omitempty"`

	DraftVersion int    `json:"draftVersion" json:"draftVersion,omitempty"`
	Location     string `json:"location" json:"location,omitempty"`

	// type agnostic --
	Bool            *bool       `json:"bool,omitempty"` // boolean schema
	ID              string      `json:"ID,omitempty"`
	Ref             *Schema     `json:"ref,omitempty"`
	Anchor          string      `json:"anchor,omitempty"`
	RecursiveRef    *Schema     `json:"recursiveRef,omitempty"`
	RecursiveAnchor bool        `json:"recursiveAnchor,omitempty"`
	DynamicRef      *DynamicRef `json:"dynamicRef,omitempty"`
	DynamicAnchor   string      `json:"dynamicAnchor,omitempty"` // "" if not specified
	Types           *Types      `json:"types,omitempty"`
	Enum            *Enum       `json:"enum,omitempty"`
	Const           *any        `json:"const,omitempty"`
	Not             *Schema     `json:"not,omitempty"`
	AllOf           []*Schema   `json:"allOf,omitempty"`
	AnyOf           []*Schema   `json:"anyOf,omitempty"`
	OneOf           []*Schema   `json:"oneOf,omitempty"`
	If              *Schema     `json:"if,omitempty"`
	Then            *Schema     `json:"then,omitempty"`
	Else            *Schema     `json:"else,omitempty"`
	Format          *Format     `json:"format,omitempty"`

	// object --
	MaxProperties         *int                `json:"maxProperties,omitempty"`
	MinProperties         *int                `json:"minProperties,omitempty"`
	Required              []string            `json:"required,omitempty"`
	PropertyNames         *Schema             `json:"propertyNames,omitempty"`
	Properties            map[string]*Schema  `json:"properties,omitempty"`
	PatternProperties     map[Regexp]*Schema  `json:"patternProperties,omitempty"`
	AdditionalProperties  any                 `json:"additionalProperties,omitempty"` // nil or bool or *Schema
	Dependencies          map[string]any      `json:"dependencies,omitempty"`         // value is []string or *Schema
	DependentRequired     map[string][]string `json:"dependentRequired,omitempty"`
	DependentSchemas      map[string]*Schema  `json:"dependentSchemas,omitempty"`
	UnevaluatedProperties *Schema             `json:"unevaluatedProperties,omitempty"`

	// array --
	MinItems         *int      `json:"minItems,omitempty"`
	MaxItems         *int      `json:"maxItems,omitempty"`
	UniqueItems      bool      `json:"uniqueItems,omitempty"`
	Contains         *Schema   `json:"contains,omitempty"`
	MinContains      *int      `json:"minContains,omitempty"`
	MaxContains      *int      `json:"maxContains,omitempty"`
	Items            any       `json:"items,omitempty"`           // nil or []*Schema or *Schema
	AdditionalItems  any       `json:"additionalItems,omitempty"` // nil or bool or *Schema
	PrefixItems      []*Schema `json:"prefixItems,omitempty"`
	Items2020        *Schema   `json:"items2020,omitempty"`
	UnevaluatedItems *Schema   `json:"unevaluatedItems,omitempty"`

	// string --
	MinLength        *int       `json:"minLength,omitempty"`
	MaxLength        *int       `json:"maxLength,omitempty"`
	Pattern          Regexp     `json:"pattern,omitempty"`
	ContentEncoding  *Decoder   `json:"contentEncoding,omitempty"`
	ContentMediaType *MediaType `json:"contentMediaType,omitempty"`
	ContentSchema    *Schema    `json:"contentSchema,omitempty"`

	// number --
	Maximum          *big.Rat `json:"maximum,omitempty"`
	Minimum          *big.Rat `json:"minimum,omitempty"`
	ExclusiveMaximum *big.Rat `json:"exclusiveMaximum,omitempty"`
	ExclusiveMinimum *big.Rat `json:"exclusiveMinimum,omitempty"`
	MultipleOf       *big.Rat `json:"multipleOf,omitempty"`

	Extensions []SchemaExt `json:"extensions,omitempty"`

	// annotations --
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Default     *any   `json:"default,omitempty"`
	Comment     string `json:"comment,omitempty"`
	ReadOnly    bool   `json:"readOnly,omitempty"`
	WriteOnly   bool   `json:"writeOnly,omitempty"`
	Examples    []any  `json:"examples,omitempty"`
	Deprecated  bool   `json:"deprecated,omitempty"`

	// liuxd extend field
	Order *int   `json:"order,omitempty"` // 排列顺序
	Name  string `json:"name,omitempty"`  // 名称
}

// --

type jsonType int

const (
	invalidType jsonType = 0
	nullType    jsonType = 1 << iota
	booleanType
	numberType
	integerType
	stringType
	arrayType
	objectType
	dateType
	dateTimeType
)

type JsonType = jsonType

const (
	JsonType_InvalidType  = invalidType
	JsonType_NullType     = nullType
	JsonType_BooleanType  = booleanType
	JsonType_NumberType   = numberType
	JsonType_IntegerType  = integerType
	JsonType_StringType   = stringType
	JsonType_ArrayType    = arrayType
	JsonType_ObjectType   = objectType
	JsonType_DateType     = dateType
	JsonType_DateTimeType = dateTimeType
)

type ISchemaType interface {
	GetSchemaType() string
}

func typeOf(v any) jsonType {
	switch v.(type) {
	case nil:
		return nullType
	case bool:
		return booleanType
	case json.Number, float32, float64, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return numberType
	case string:
		return stringType
	case []any:
		return arrayType
	case map[string]any:
		return objectType
	case time.Time:
		return dateTimeType
	case *time.Time:
		return dateTimeType
	default:
		if schemaType, ok := v.(ISchemaType); ok {
			sType := schemaType.GetSchemaType()
			return typeFromString(sType)
		}
		if _, ok := v.(*time.Time); ok {
			return dateTimeType
		}
		return invalidType
	}
}

func typeFromString(s string) jsonType {
	switch s {
	case "date":
		return dateType
	case "datetime":
		return dateTimeType
	case "null":
		return nullType
	case "boolean":
		return booleanType
	case "number":
		return numberType
	case "integer":
		return integerType
	case "string":
		return stringType
	case "array":
		return arrayType
	case "object":
		return objectType
	}
	return invalidType
}

func (jt jsonType) String() string {
	switch jt {
	case nullType:
		return "null"
	case booleanType:
		return "boolean"
	case numberType:
		return "number"
	case integerType:
		return "integer"
	case stringType:
		return "string"
	case arrayType:
		return "array"
	case objectType:
		return "object"
	case dateType:
		return "date"
	case dateTimeType:
		return "datetime"
	}
	return ""
}

// --

// Types encapsulates list of json value types.
type Types int

func newTypes(v any) *Types {
	var types Types
	switch v := v.(type) {
	case string:
		types.add(typeFromString(v))
	case []any:
		for _, item := range v {
			if s, ok := item.(string); ok {
				types.add(typeFromString(s))
			}
		}
	}
	if types.IsEmpty() {
		return nil
	}
	return &types
}

func (tt Types) IsEmpty() bool {
	return tt == 0
}

func (tt *Types) add(t jsonType) {
	*tt = Types(int(*tt) | int(t))
}

func (tt Types) contains(t jsonType) bool {
	val := int(tt)&int(t) != 0
	return val
}

func (tt Types) Contains(t jsonType) bool {
	return tt.contains(t)
}

func (tt Types) ToStrings() []string {
	types := []jsonType{
		nullType, booleanType, numberType, integerType,
		stringType, arrayType, objectType, dateType, dateTimeType,
	}
	var arr []string
	for _, t := range types {
		if tt.contains(t) {
			arr = append(arr, t.String())
		}
	}
	return arr
}

func (tt Types) String() string {
	return fmt.Sprintf("%v", tt.ToStrings())
}

// --

type Enum struct {
	Values []any
	types  Types
}

func newEnum(arr []any) *Enum {
	var types Types
	for _, item := range arr {
		types.add(typeOf(item))
	}
	return &Enum{arr, types}
}

// --

type DynamicRef struct {
	Ref    *Schema
	Anchor string // "" if not specified
}

func newSchema(up urlPtr) *Schema {
	return &Schema{up: up, Location: up.String()}
}
