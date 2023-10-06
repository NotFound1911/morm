package morm

import (
	errs "github.com/NotFound1911/morm/internal/pkg/errors"
	"reflect"
	"sync"
)

// 元数据注册中心
type registry struct {
	models sync.Map // sync.Map 解决并发问题
}

func (r *registry) get(val any) (*model, error) {
	return r.parseModel(val)
}
func (r *registry) parseModel(val any) (*model, error) {
	typ := reflect.TypeOf(val)
	m, ok := r.models.Load(typ)
	if !ok {
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
		m = &model{
			tableName: underscoreName(typ.Name()),
			fieldMap:  fds,
		}
	}
	r.models.Store(typ, m)
	return m.(*model), nil
}
