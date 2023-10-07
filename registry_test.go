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
		wantModel *model
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
			wantModel: &model{
				tableName: "test_model",
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
			val:  &ColumnTag{},
			wantModel: &model{
				tableName: "column_tag",
				fieldMap: map[string]*field{
					"Id": {
						colName: "id",
					},
				},
			},
		},
		{
			name: "empty column",
			val:  &EmptyColumn{},
			wantModel: &model{
				tableName: "empty_column",
				fieldMap: map[string]*field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
		},
		{
			// 设置了column 但没有赋值
			name:    "invalid tag",
			val:     &InvalidTag{},
			wantErr: errs.NewErrInvalidTagContent("column"),
		},
		{
			// 设置了非定义tag 进行忽略
			name: "ignore tag",
			val:  &IgnoreTag{},
			wantModel: &model{
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
			wantModel: &model{
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
			wantModel: &model{
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
			wantModel: &model{
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
			m, err := r.get(tc.val)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantModel, m)
		})
	}
}

type InvalidTag struct {
	FirstName uint64 `morm:"column"`
}
type ColumnTag struct {
	Id uint64 `morm:"column=id"`
}

type EmptyColumn struct {
	FirstName uint64 `morm:"column="`
}

type IgnoreTag struct {
	FirstName uint64 `morm:"aaa=bbb"`
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
