package model

import (
	errs "github.com/NotFound1911/morm/internal/pkg/errors"
	"unicode"
)

type ModelOpt func(model *Model) error

type Model struct {
	FieldMap  map[string]*field // 字段
	TableName string            // 表名
}

// field 字段
type field struct {
	ColName string
}

// underscoreName 驼峰转字符串命名
func underscoreName(name string) string {
	var buf []byte
	for i, v := range name {
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
		model.TableName = name
		return nil
	}
}

func ModelWithColumnName(field string, columnName string) ModelOpt {
	return func(model *Model) error {
		fd, ok := model.FieldMap[field]
		if !ok {
			return errs.NewErrUnknownField(field)
		}
		fd.ColName = columnName
		return nil
	}
}
