package morm

import (
	errs "github.com/NotFound1911/morm/internal/pkg/errors"
	"unicode"
)

type ModelOpt func(model *Model) error

type Model struct {
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

// 支持的tag 标签
const (
	tagKeyColumn = "column"
)

// TableName 用户实现这个接口来返回自定义的表名
type TableName interface {
	TableName() string
}

func ModelWitTableName(name string) ModelOpt {
	return func(model *Model) error {
		model.tableName = name
		return nil
	}
}

func ModelWithColumnName(field string, columnName string) ModelOpt {
	return func(model *Model) error {
		fd, ok := model.fieldMap[field]
		if !ok {
			return errs.NewErrUnknownField(field)
		}
		fd.colName = columnName
		return nil
	}
}
