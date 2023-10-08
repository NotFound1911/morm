package morm

import (
	errs "github.com/NotFound1911/morm/internal/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegistry_get(t *testing.T) {
	mp := func() any {
		val := &TestModel{}
		return &val
	}
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
				tableName: "TestModel",
				fieldMap: map[string]*field{
					"Id": {
						colName: "id",
					},
					"FirstName": {
						colName: "first_name",
					},
					"Age": {
						colName: "age",
					},
					"LastName": {
						colName: "last_name",
					},
				},
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
				tableName: "column_tag",
				fieldMap: map[string]*field{
					"ID": {
						colName: "id",
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
				tableName: "empty_column",
				fieldMap: map[string]*field{
					"FirstName": {
						colName: "first_name",
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
			// 如果用户设置了一些奇奇怪怪的内容，这部分内容我们会忽略掉
			name: "ignore tag",
			val: func() any {
				// 我们把测试结构体定义在方法内部，防止被其它用例访问
				type IgnoreTag struct {
					FirstName uint64 `orm:"abc=abc"`
				}
				return &IgnoreTag{}
			}(),
			wantModel: &Model{
				tableName: "ignore_tag",
				fieldMap: map[string]*field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
		},
		// interface test
		{
			name: "custom table name",
			val:  &CustomTableName{},
			wantModel: &Model{
				tableName: "test_custom_table_name",
				fieldMap: map[string]*field{
					"Name": {
						colName: "name",
					},
				},
			},
		},
		{
			name: "custom table name ptr",
			val:  &CustomTableNamePtr{},
			wantModel: &Model{
				tableName: "test_custom_table_name_ptr",
				fieldMap: map[string]*field{
					"Name": {
						colName: "name",
					},
				},
			},
		},
		{
			name: "empty table name",
			val:  &EmptyTableName{},
			wantModel: &Model{
				tableName: "empty_table_name",
				fieldMap: map[string]*field{
					"Name": {
						colName: "name",
					},
				},
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
