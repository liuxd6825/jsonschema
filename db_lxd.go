package jsonschema

import "math/big"
import "encoding/json"
import "fmt"

type DBField struct {
	//Name   string // 属性名
	DBName string // 数据库字段名
	//DataType   DataType // 数据类型
	PrimaryKey bool   // 是主键
	NotNull    bool   // 不能为空
	Unique     bool   // 唯一键值
	Size       *int   // 字段大小
	Creatable  bool   // 是否可创建
	Updatable  bool   // 是否可更新
	Readable   bool   // 可读取的字段
	Order      string // 是索引类型 asc desc none
	obj        map[string]any
}

type DBOrder = string

const (
	DBOrder_Asc  DBOrder = "asc"
	DBOrder_Desc DBOrder = "desc"
)

func NewDBField(name string, m map[string]any) *DBField {
	field := &DBField{
		obj:        m,
		DBName:     name,
		PrimaryKey: false,
		NotNull:    false,
		Unique:     false,
		Size:       nil,
		Creatable:  true,
		Updatable:  true,
		Readable:   true,
	}
	if m != nil {
		if val := field.string("dbName"); val != "" {
			field.DBName = val
		}
		field.PrimaryKey = field.boolean("primaryKey")
		field.NotNull = field.boolean("notNull")
		field.Unique = field.boolean("unique")
		field.Updatable = field.boolean("updatable")
		field.Readable = field.boolean("readable")
		field.Size = field.intVal("size")
		field.Order = field.string("order")
	}

	return field
}

// value helpers --

func (c *DBField) boolVal(pname string) *bool {
	v, ok := c.obj[pname]
	if !ok {
		return nil
	}
	b, ok := v.(bool)
	if !ok {
		return nil
	}
	return &b
}

func (c *DBField) boolean(pname string) bool {
	b := c.boolVal(pname)
	return b != nil && *b
}

func (c *DBField) strVal(pname string) *string {
	v, ok := c.obj[pname]
	if !ok {
		return nil
	}
	s, ok := v.(string)
	if !ok {
		return nil
	}
	return &s
}

func (c *DBField) string(pname string) string {
	if s := c.strVal(pname); s != nil {
		return *s
	}
	return ""
}

func (c *DBField) numVal(pname string) *big.Rat {
	v, ok := c.obj[pname]
	if !ok {
		return nil
	}
	switch v.(type) {
	case json.Number, float32, float64, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		if n, ok := new(big.Rat).SetString(fmt.Sprint(v)); ok {
			return n
		}
	}
	return nil
}

func (c *DBField) intVal(pname string) *int {
	if n := c.numVal(pname); n != nil && n.IsInt() {
		n := int(n.Num().Int64())
		return &n
	}
	return nil
}

func (c *DBField) objVal(pname string) map[string]any {
	v, ok := c.obj[pname]
	if !ok {
		return nil
	}
	obj, ok := v.(map[string]any)
	if !ok {
		return nil
	}
	return obj
}

func (c *DBField) arrVal(pname string) []any {
	v, ok := c.obj[pname]
	if !ok {
		return nil
	}
	arr, ok := v.([]any)
	if !ok {
		return nil
	}
	return arr
}
