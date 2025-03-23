package jsonschema

import "math/big"
import "encoding/json"
import "fmt"
import "strings"

type DBProperty struct {
	//Name   string // 属性名
	Name string // 数据库字段名
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

func NewDBProperty(name string, m map[string]any) *DBProperty {

	field := &DBProperty{
		obj:        m,
		Name:       SnakeString(name),
		PrimaryKey: false,
		NotNull:    false,
		Unique:     false,
		Size:       nil,
		Creatable:  true,
		Updatable:  true,
		Readable:   true,
	}
	if m != nil {
		if val := field.string("name"); val != "" {
			field.Name = val
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

func (c *DBProperty) boolVal(pname string) *bool {
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

func (c *DBProperty) boolean(pname string) bool {
	b := c.boolVal(pname)
	return b != nil && *b
}

func (c *DBProperty) strVal(pname string) *string {
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

func (c *DBProperty) string(pname string) string {
	if s := c.strVal(pname); s != nil {
		return *s
	}
	return ""
}

func (c *DBProperty) numVal(pname string) *big.Rat {
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

func (c *DBProperty) intVal(pname string) *int {
	if n := c.numVal(pname); n != nil && n.IsInt() {
		n := int(n.Num().Int64())
		return &n
	}
	return nil
}

func (c *DBProperty) objVal(pname string) map[string]any {
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

func (c *DBProperty) arrVal(pname string) []any {
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

// SnakeString
// @Description: 驼峰转蛇形
// @param s 要转换的字符串
// @return string
func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		// or通过ASCII码进行大小写的转化
		// 65-90（A-Z），97-122（a-z）
		//判断如果字母为大写的A-Z就在前面拼接一个_
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	//ToLower把大写字母统一转小写
	res := strings.ToLower(string(data[:]))
	if strings.HasPrefix(res, "_") {
		return res[1:]
	}
	return res
}
