package morm

import (
	"unicode"
)

type model struct {
	fieldMap  map[string]*field // 字段
	tableName string            // 表名
}

// field 字段
type field struct {
	colName string
}

// underscoreName 驼峰转字符串命名
func underscoreName(tableName string) string {
	var buf []byte
	for i, v := range tableName {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}

	}
	return string(buf)
}
