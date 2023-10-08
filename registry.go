package morm

import (
	errs "github.com/NotFound1911/morm/internal/pkg/errors"
	"reflect"
	"strings"
	"sync"
)

// Registry 元数据注册中心的抽象
type Registry interface {
	// Get 查询元数据
	Get(val any) (*Model, error)
	// Register 注册一个模型
	Register(val any, opts ...ModelOpt) (*Model, error)
}

// 元数据注册中心
type registry struct {
	models sync.Map // sync.Map 解决并发问题
}

func NewRegistry() Registry {
	return &registry{}
}
func (r *registry) Get(val any) (*Model, error) {
	typ := reflect.TypeOf(val)
	m, ok := r.models.Load(typ)
	if ok {
		return m.(*Model), nil
	}
	return r.Register(val)
}

func (r *registry) Register(val any, opts ...ModelOpt) (*Model, error) {
	m, err := r.parseModel(val)
	if err != nil {
		return nil, err
	}
	for _, opt := range opts {
		if err := opt(m); err != nil {
			return nil, err
		}
	}
	typ := reflect.TypeOf(val)
	r.models.Store(typ, m)
	return m, nil
}

// parseModel 支持从标签中提取自定义设置
// 标签形式 morm:"key1=value1,key2=value2"
func (r *registry) parseModel(val any) (*Model, error) {
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
		// 解析tag
		tags, err := r.parseTag(fdType.Tag)
		if err != nil {
			return nil, err
		}
		colName := tags[tagKeyColumn]
		if colName == "" {
			colName = underscoreName(fdType.Name)
		}
		fds[fdType.Name] = &field{
			colName: colName,
		}
	}
	var tableName string
	if tn, ok := val.(TableName); ok {
		tableName = tn.TableName()
	}
	if tableName == "" {
		tableName = underscoreName(typ.Name())
	}
	return &Model{
		tableName: tableName,
		fieldMap:  fds,
	}, nil
}
func (r *registry) parseTag(tag reflect.StructTag) (map[string]string, error) {
	ormTag := tag.Get("morm")
	if ormTag == "" {
		return map[string]string{}, nil
	}
	res := make(map[string]string, 1)

	pairs := strings.Split(ormTag, ",")
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			return nil, errs.NewErrInvalidTagContent(pair)
		}
		res[kv[0]] = kv[1]
	}
	return res, nil
}
