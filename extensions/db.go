package extensions

import "encoding/json"

type DBOrder string

const (
	DBOrder_None DBOrder = ""
	DBOrder_Asc  DBOrder = "asc"
	DBOrder_Desc DBOrder = "desc"
)

type DBField struct {
	Name       string  `json:"name"`
	Size       int64   `json:"size"`
	PrimaryKey bool    `json:"primaryKey"`
	NotNull    bool    `json:"notNull"`
	Unique     bool    `json:"unique"`
	Creatable  bool    `json:"creatable"`
	Updatable  bool    `json:"updatable"`
	Readable   bool    `json:"readable"`
	Order      DBOrder `json:"order"`
}

type DBTable struct {
	Name string `json:"name"`
}

func (db *DBTable) init(values map[string]any) error {
	var err error
	for k, v := range values {
		switch k {
		case "name":
			db.Name = v.(string)
		}
	}
	return err
}

func (db *DBField) init(values map[string]any) error {
	var err error
	for key, value := range values {
		switch key {
		case "primaryKey":
			if val, ok := value.(bool); ok {
				db.PrimaryKey = val
			}
		case "creatable":
			if val, ok := value.(bool); ok {
				db.Creatable = val
			}
		case "size":
			if val, ok := value.(json.Number); ok {
				db.Size, err = val.Int64()
			}
		case "name":
			if val, ok := value.(string); ok {
				db.Name = val
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
		if err != nil {
			return err
		}
	}
	return err
}
