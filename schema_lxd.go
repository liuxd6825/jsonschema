package jsonschema

func (sch *Schema) GetAllProperties() map[string]*Schema {
	props := make(map[string]*Schema)
	sch.addMap(sch, props)
	// 添加allOf属性
	for _, item := range sch.AllOf {
		sch.addMap(item, props)
	}

	return props
}

func (sch *Schema) addMap(source *Schema, target map[string]*Schema) map[string]*Schema {
	for key, value := range source.Properties {
		target[key] = value
	}
	// 添加ref属性
	if source.Ref != nil && source.Ref.Properties != nil {
		sch.addMap(source.Ref, target)
	}
	return target
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

	if s.Types != nil {
		if s.Types.Contains(JsonType_ArrayType) {
			items := s.Items2020
			if items != nil && items.Types.Contains(JsonType_ObjectType) {
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

// getFieldsProps
//
//	@Description:
//	@receiver sch
//	@param props
//	@param fields
//	@param types
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
			}
		}
	}
}
