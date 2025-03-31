package jsonschema

type SchemaView struct {
	Name        string                 `json:"name,omitempty"`
	Title       string                 `json:"title,omitempty"`
	Type        []string               `json:"type,omitempty"`
	Meta        any                    `json:"meta,omitempty"`
	Properties  map[string]*SchemaView `json:"properties,omitempty"`
	Order       *int                   `json:"order,omitempty"`
	Items       map[string]*SchemaView `json:"items,omitempty"`
	MinLength   *int                   `json:"minLength,omitempty"`
	MaxLength   *int                   `json:"maxLength,omitempty"`
	Description string                 `json:"description,omitempty"`
	Default     any                    `json:"default,omitempty"`
	Comment     string                 `json:"comment,omitempty"`
	ReadOnly    bool                   `json:"readOnly,omitempty"`
	WriteOnly   bool                   `json:"writeOnly,omitempty"`
	Examples    []any                  `json:"examples,omitempty"`
	Deprecated  bool                   `json:"deprecated,omitempty"`
	Required    []string               `json:"required,omitempty"`
}

func NewSchemaView(sch *Schema) *SchemaView {
	schemaView := &SchemaView{
		Name:        sch.Name(),
		Title:       sch.Title,
		Type:        sch.Types.ToStrings(),
		Order:       sch.Order,
		Meta:        sch.GetExtensions("meta"),
		MinLength:   sch.MinLength,
		MaxLength:   sch.MaxLength,
		Description: sch.Description,
		Default:     sch.Default,
		ReadOnly:    sch.ReadOnly,
		WriteOnly:   sch.WriteOnly,
		Examples:    sch.Examples,
		Deprecated:  sch.Deprecated,
		Required:    sch.Required,
	}

	props := sch.GetAllProperties()
	if len(props) > 0 {
		schemaView.Properties = map[string]*SchemaView{}
		for k, o := range props {
			schemaView.Properties[k] = NewSchemaView(o)
		}
	}

	return schemaView
}
