package model

import (
	errs "github.com/NotFound1911/morm/internal/pkg/errors"
	"reflect"
	"unicode"
)

type Opt func(model *Model) error

type Model struct {
	FieldMap  map[string]*Field // 字段（go）
	TableName string            // 表名
	ColumnMap map[string]*Field // 列名（sql）
	Fields    []*Field
}

// Field 字段
type Field struct {
	ColName string
	// Offset 相对于对象起始地址的字段偏移量
	Offset uintptr
	// Type 类型
	Type reflect.Type
	// Go字段名
	GoName string
	Index  int
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

func WitTableName(name string) Opt {
	return func(model *Model) error {
		model.TableName = name
		return nil
	}
}

func WithColumnName(field string, columnName string) Opt {
	return func(model *Model) error {
		fd, ok := model.FieldMap[field]
		if !ok {
			return errs.NewErrUnknownField(field)
		}
		fd.ColName = columnName
		return nil
	}
}
