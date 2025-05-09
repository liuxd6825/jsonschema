package jsonschema

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

type TagExtension interface {
	TagName() string
}

type WithJsonOptions func(compiler *Compiler)

func NewSchemaWithJson(fileName string, jsonText string, opts ...WithJsonOptions) (*Schema, error) {
	return NewSchemaWithBytes(fileName, []byte(jsonText), opts...)
}

func NewSchemaWithBytes(fileName string, fileBytes []byte, opts ...WithJsonOptions) (*Schema, error) {
	compiler := NewCompiler()
	compiler.AssertFormat()
	compiler.AssertContent()

	for _, o := range opts {
		if o != nil {
			o(compiler)
		}
	}

	data, err := UnmarshalJSON(bytes.NewReader(fileBytes))
	if err != nil {
		return nil, err
	}
	if err = compiler.AddResource(fileName, data); err != nil {
		return nil, err
	}
	sch, err := compiler.Compile(fileName)
	if err != nil {
		return nil, err
	}
	return sch, err
}

func NewSchemaWithStruct(fileName string, obj any, opts ...WithJsonOptions) (*Schema, error) {
	jsonBytes, err := ReflectJson(obj)
	if err != nil {
		return nil, err
	}
	return NewSchemaWithJson(fileName, string(jsonBytes), opts...)
}

// GetAllProperties
//
//	@Description: 取得所有属性
//	@receiver sch
//	@return map[string]*Schema
func (sch *Schema) GetAllProperties() map[string]*Schema {
	props := make(map[string]*Schema)
	sch.addMap(sch, props)
	// 添加allOf属性
	for _, item := range sch.AllOf {
		sch.addMap(item, props)
	}
	return props
}

// GetSortProperties
//
//	@Description: 取得所有排序后的属性
//	@receiver sch
//	@return []*Schema
func (sch *Schema) GetSortProperties() []*Schema {
	props := sch.GetAllProperties()
	return sch.SortSchemas(props)
}

func (sch *Schema) GetExtensions(tagName string) any {
	if len(sch.Extensions) == 0 {
		return nil
	}
	for _, v := range sch.Extensions {
		if ext, ok := v.(TagExtension); ok {
			if ext.TagName() == tagName {
				return v
			}
		}
	}
	return nil
}

func (sch *Schema) GetView() *SchemaView {
	return NewSchemaView(sch)
}

// SortSchemas
//
//	@Description: 对schema进行排序
//	@receiver sch
//	@param fieldsMap
//	@return []*Schema
func (sch *Schema) SortSchemas(schMaps map[string]*Schema) []*Schema {
	// Convert map to slice
	fields := make([]*Schema, 0, len(schMaps))
	for _, field := range schMaps {
		//field.Name = name
		fields = append(fields, field)
	}

	// Sort slice by "Order"
	sort.Slice(fields, func(i, j int) bool {
		if fields[i].Order != nil && fields[j].Order != nil {
			iOrder := fields[i].Order
			jOrder := fields[j].Order
			return *iOrder < *jOrder
		}
		return false
	})

	return fields
}

func (sch *Schema) addMap(source *Schema, target map[string]*Schema) map[string]*Schema {
	for key, srcProp := range source.Properties {
		if _, ok := target[key]; ok {
			continue
		} else {
			target[key] = srcProp
		}

	}
	// 添加ref属性
	if source.Ref != nil {
		if source.Ref.Properties != nil {
			sch.addMap(source.Ref, target)
		}
		for _, item := range source.Ref.AllOf {
			sch.addMap(item, target)
		}
		for _, item := range source.Ref.AnyOf {
			sch.addMap(item, target)
		}
	}
	return target
}

func (sch *Schema) GetType() string {
	list := sch.Types.ToStrings()
	for _, item := range list {
		if item != "null" {
			return item
		}
	}
	return "null"
}

// IsRequired
//
//	@Description:  Helper function to check if a field is in the "required" list
//	@receiver sch
//	@param name
//	@return bool
func (sch *Schema) IsRequired(name string) bool {
	for _, r := range sch.Required {
		if r == name {
			return true
		}
	}
	return false
}

// GetFields
//
//	@Description:  深度获取schema的属性
//	@receiver sch
//	@param ...JsonType  想要获取的类型
//	@return map[string]any
func (sch *Schema) GetFields(types ...JsonType) map[string]any {
	fields := map[string]any{}
	sch.getFieldsType(sch, fields, "", types...)
	return fields
}

func (sch *Schema) GetTimeFields() map[string]any {
	if sch.timeFields == nil {
		sch.timeFields = sch.GetFields(JsonType_DateTimeType, JsonType_DateType)
	}
	return sch.timeFields
}

// getFieldsType
//
//	@Description: 深度获取schema的属性
//	@receiver sch
//	@param s
//	@param fields
//	@param propName
func (sch *Schema) getFieldsType(s *Schema, fields map[string]any, parentKey string, types ...JsonType) {
	if s == nil {
		return
	}
	if s.Types == nil {
		panic(fmt.Sprintf("jsonschema types is null, location:%s", s.Location))
	}
	if s.Types != nil {
		if s.Types.Contains(JsonType_ArrayType) {
			items := s.Items2020
			if items != nil {
				propsTimeFields := map[string]any{}
				props := items.GetAllProperties()
				sch.getFieldsProps(props, propsTimeFields, parentKey, types...)
				for k, v := range propsTimeFields {
					fields[k] = v
				}
			}
		} else if s.Types.Contains(JsonType_ObjectType) {
			props := sch.GetAllProperties()
			sch.getFieldsProps(props, fields, parentKey, types...)
		}
	}
	/*
		if s.Ref != nil {
			sch.getFieldsType(s.Ref, fields, propName, types...)
		}
		if s.AllOf != nil {
			for _, item := range s.AllOf {
				if item.Ref != nil {
					sch.getFieldsType(item.Ref, fields, propName, types...)
				}
			}
		}
	*/
}

func (sch *Schema) getName() string {
	// 查找 "properties/" 的最后一个位置
	index := strings.LastIndex(sch.Location, "properties/")
	if index == -1 {
		// 如果没有找到 "properties/"，返回空字符串
		return ""
	}

	// 提取 "properties/" 之后的内容
	return sch.Location[index+len("properties/"):]
}

func (sch *Schema) Name() string {
	return sch.name
}

// GetTitleByLocation 查找属性标题，支持多级查找
func (sch *Schema) GetTitleByLocation(location []string) string {
	findSch := sch
	findOk := true
	titles := make([]string, 0)
	for _, name := range location {
		if isInteger(name) {
			titles = append(titles, "第"+name+"项")
			continue
		}
		p, ok := findSch.Properties[name]
		if !ok && findSch.Items2020 != nil {
			if findSch.Items2020.Ref != nil {
				p, ok = findSch.Items2020.Ref.Properties[name]
			} else {
				p, ok = findSch.Items2020.Properties[name]
			}
		}
		if ok {
			title := p.Title
			if len(title) == 0 {
				title = p.Name()
			}
			titles = append(titles, title)
			findSch = p
		} else {
			findOk = false
			break
		}
	}
	if !findOk {
		return ""
	}
	return strings.Join(titles, ".")
}

func (sch *Schema) GetPropertyByLocation(location []string) (list []*Schema, findOk bool) {
	list, findOk = GetPropertyByLocation(sch, location)
	return list, findOk
}

// GetPropertyByLocation 查找属性，支持多级查找
func GetPropertyByLocation(rootSch *Schema, location []string) (list []*Schema, findOk bool) {
	list = make([]*Schema, 0)
	sch := rootSch
	findOk = true
	for _, name := range location {
		if isInteger(name) {
			continue
		}
		p, ok := sch.Properties[name]
		if ok {
			list = append(list, p)
			sch = p
		} else {
			findOk = false
		}
	}
	return list, findOk
}

// getFieldsProps
//
//	@Description: 返回属性中的指定类型，支持深度搜索
//	@receiver sch
//	@param props 要搜索的属性集
//	@param fields 返回的字段集
//	@param types 要搜索的类型
func (sch *Schema) getFieldsProps(props map[string]*Schema, fields map[string]any, parentKey string, types ...JsonType) {
	if props == nil {
		return
	}

	for k, p := range props {
		if p == nil {
			continue
		}
		if p.Types != nil {
			if len(types) == 0 {
				fields[k] = p
			} else {
				for _, t := range types {
					if p.Types.Contains(t) {
						fields[k] = p
						break
					}
				}
			}
			if p.Types.Contains(JsonType_ObjectType) {
				subFields := map[string]any{}
				p.getFieldsType(p, subFields, k, types...)
				if len(subFields) > 0 {
					fields[k] = subFields
				}
			} else if p.Types.Contains(JsonType_ArrayType) {
				subFields := map[string]any{}
				p.getFieldsType(p, subFields, k, types...)
				if len(subFields) > 0 {
					fields[k] = subFields
				}
			}
		} else if p.Ref != nil {
			subFields := map[string]any{}
			props := p.Ref.GetAllProperties()
			p.Ref.getFieldsProps(props, subFields, k, types...)
			if len(subFields) > 0 {
				fields[k] = subFields
			}
		} else if p.Items2020 != nil {
			subFields := map[string]any{}
			props := p.Items2020.GetAllProperties()
			p.Ref.getFieldsProps(props, subFields, k, types...)
			if len(subFields) > 0 {
				fields[k] = subFields
			}
		}
	}
}
