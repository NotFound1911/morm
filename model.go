package morm

import (
	errs "github.com/NotFound1911/morm/internal/pkg/errors"
	"reflect"
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

func parseModel(val any) (*model, error) {
	typ := reflect.TypeOf(val)
	// 只支持一级指针
	if typ.Kind() != reflect.Ptr || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.NewErrPointerOnly(val)
	}
	typ = typ.Elem()
	numField := typ.NumField()
	fds := make(map[string]*field, numField)
	for i := 0; i < numField; i++ {
		fdType := typ.Field(i)
		fds[fdType.Name] = &field{
			colName: underscoreName(fdType.Name),
		}
	}
	return &model{
		tableName: underscoreName(typ.Name()),
		fieldMap:  fds,
	}, nil
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
