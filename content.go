package jsonschema

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
)

// Decoder specifies how to decode specific contentEncoding.
type Decoder struct {
	// Name of contentEncoding.
	Name string
	// Decode given string to byte array.
	Decode func(string) ([]byte, error) `json:"-"`
}

var decoders = map[string]*Decoder{
	"base64": {
		Name: "base64",
		Decode: func(s string) ([]byte, error) {
			return base64.StdEncoding.DecodeString(s)
		},
	},
}

// MediaType specified how to validate bytes against specific contentMediaType.
type MediaType struct {
	// Name of contentMediaType.
	Name string `json:"name"`

	// Validate checks whether bytes conform to this mediatype.
	Validate func([]byte) error `json:"-"`

	// UnmarshalJSON unmarshals bytes into json value.
	// This must be nil if this mediatype is not compatible
	// with json.
	UnmarshalJSON func([]byte) (any, error) `json:"-"`
}

var mediaTypes = map[string]*MediaType{
	"application/json": {
		Name: "application/json",
		Validate: func(b []byte) error {
			var v any
			return json.Unmarshal(b, &v)
		},
		UnmarshalJSON: func(b []byte) (any, error) {
			return UnmarshalJSON(bytes.NewReader(b))
		},
	},
}
