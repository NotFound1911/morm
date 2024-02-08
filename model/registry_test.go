package model

import (
	"github.com/NotFound1911/morm/errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestRegistry_get(t *testing.T) {
	mp := func() any {
		val := &TestModel{}
		return &val
	}
	var tm TestModel
	testCases := []struct {
		name      string
		val       any
		wantModel *Model
		wantErr   error
	}{
		{
			name:    "test model",
			val:     TestModel{},
			wantErr: errs.NewErrPointerOnly(TestModel{}),
		},
		{
			name: "pointer",
			val:  &TestModel{},
			wantModel: &Model{
				TableName: "TestModel",
				FieldMap: map[string]*Field{
					"Id":        tm.IdField(),
					"FirstName": tm.FirstNameField(),
					"Age":       tm.AgeField(),
					"LastName":  tm.LastNameField(),
				},
				ColumnMap: map[string]*Field{
					"id":         tm.IdField(),
					"first_name": tm.FirstNameField(),
					"age":        tm.AgeField(),
					"last_name":  tm.LastNameField(),
				},
				Fields: []*Field{tm.IdField(), tm.FirstNameField(), tm.AgeField(), tm.LastNameField()},
			},
		},
		{
			name:    "multiple pointer",
			val:     mp,
			wantErr: errs.NewErrPointerOnly(mp),
		},
		{
			name:    "map",
			val:     map[string]any{},
			wantErr: errs.NewErrPointerOnly(map[string]any{}),
		},
		{
			name:    "slice",
			val:     []int{},
			wantErr: errs.NewErrPointerOnly([]int{}),
		},
		{
			name:    "basic type",
			val:     0,
			wantErr: errs.NewErrPointerOnly(0),
		},
		// tag test
		{
			name: "column tag",
			val: func() any {
				// 我们把测试结构体定义在方法内部，防止被其它用例访问
				type ColumnTag struct {
					ID uint64 `morm:"column=id"`
				}
				return &ColumnTag{}
			}(),
			wantModel: &Model{
				TableName: "column_tag",
				FieldMap: map[string]*Field{
					"ID": {
						ColName: "id",
						Type:    reflect.TypeOf(uint64(0)),
						GoName:  "ID",
					},
				},
				ColumnMap: map[string]*Field{
					"id": {
						ColName: "id",
						Type:    reflect.TypeOf(uint64(0)),
						GoName:  "ID",
					},
				},
				Fields: []*Field{
					{
						ColName: "id",
						Type:    reflect.TypeOf(uint64(0)),
						GoName:  "ID",
					},
				},
			},
		},
		{
			// 如果用户设置了 column，但是传入一个空字符串，那么会用默认的名字
			name: "empty column",
			val: func() any {
				// 我们把测试结构体定义在方法内部，防止被其它用例访问
				type EmptyColumn struct {
					FirstName uint64 `morm:"column="`
				}
				return &EmptyColumn{}
			}(),
			wantModel: &Model{
				TableName: "empty_column",
				FieldMap: map[string]*Field{
					"FirstName": {
						ColName: "first_name",
						Type:    reflect.TypeOf(uint64(0)),
						GoName:  "FirstName",
					},
				},
				ColumnMap: map[string]*Field{
					"first_name": {
						ColName: "first_name",
						Type:    reflect.TypeOf(uint64(0)),
						GoName:  "FirstName",
					},
				},
				Fields: []*Field{
					{
						ColName: "first_name",
						Type:    reflect.TypeOf(uint64(0)),
						GoName:  "FirstName",
					},
				},
			},
		},
		{
			// 如果用户设置了 column，但是没有赋值
			name: "invalid tag",
			val: func() any {
				// 我们把测试结构体定义在方法内部，防止被其它用例访问
				type InvalidTag struct {
					FirstName uint64 `morm:"column"`
				}
				return &InvalidTag{}
			}(),
			wantErr: errs.NewErrInvalidTagContent("column"),
		},
		{
			// 如果用户设置了一些不规范的内容，这部分内容会忽略掉
			name: "ignore tag",
			val: func() any {
				type IgnoreTag struct {
					FirstName uint64 `orm:"abc=abc"`
				}
				return &IgnoreTag{}
			}(),
			wantModel: &Model{
				TableName: "ignore_tag",
				FieldMap: map[string]*Field{
					"FirstName": {
						ColName: "first_name",
						Type:    reflect.TypeOf(uint64(0)),
						GoName:  "FirstName",
					},
				},
				ColumnMap: map[string]*Field{
					"first_name": {
						ColName: "first_name",
						Type:    reflect.TypeOf(uint64(0)),
						GoName:  "FirstName",
					},
				},
				Fields: []*Field{{
					ColName: "first_name",
					Type:    reflect.TypeOf(uint64(0)),
					GoName:  "FirstName",
				}},
			},
		},
		// interface test
		{
			name: "custom table name",
			val:  &CustomTableName{},
			wantModel: &Model{
				TableName: "test_custom_table_name",
				FieldMap: map[string]*Field{
					"Name": {
						ColName: "name",
						Type:    reflect.TypeOf(""),
						GoName:  "Name",
					},
				},
				ColumnMap: map[string]*Field{
					"name": {
						ColName: "name",
						Type:    reflect.TypeOf(""),
						GoName:  "Name",
					},
				},
				Fields: []*Field{{
					ColName: "name",
					GoName:  "Name",
					Type:    reflect.TypeOf(""),
				}},
			},
		},
		{
			name: "custom table name ptr",
			val:  &CustomTableNamePtr{},
			wantModel: &Model{
				TableName: "test_custom_table_name_ptr",
				FieldMap: map[string]*Field{
					"Name": {
						ColName: "name",
						Type:    reflect.TypeOf(""),
						GoName:  "Name",
					},
				},
				ColumnMap: map[string]*Field{
					"name": {
						ColName: "name",
						Type:    reflect.TypeOf(""),
						GoName:  "Name",
					},
				},
				Fields: []*Field{{
					ColName: "name",
					GoName:  "Name",
					Type:    reflect.TypeOf(""),
				}},
			},
		},
		{
			name: "empty table name",
			val:  &EmptyTableName{},
			wantModel: &Model{
				TableName: "empty_table_name",
				FieldMap: map[string]*Field{
					"Name": {
						ColName: "name",
						Type:    reflect.TypeOf(""),
						GoName:  "Name",
					},
				},
				ColumnMap: map[string]*Field{
					"name": {
						ColName: "name",
						Type:    reflect.TypeOf(""),
						GoName:  "Name",
					},
				},
				Fields: []*Field{{
					ColName: "name",
					GoName:  "Name",
					Type:    reflect.TypeOf(""),
				}},
			},
		},
	}
	r := &registry{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.Get(tc.val)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantModel, m)
		})
	}
}

type CustomTableName struct {
	Name string
}

func (c CustomTableName) TableName() string {
	return "test_custom_table_name"
}

type CustomTableNamePtr struct {
	Name string
}

func (c *CustomTableNamePtr) TableName() string {
	return "test_custom_table_name_ptr"
}

type EmptyTableName struct {
	Name string
}

func (c *EmptyTableName) TableName() string {
	return ""
}
