package model

import (
	"database/sql"
	errs "github.com/NotFound1911/morm/internal/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}

func (t TestModel) TableName() string {
	return "TestModel"
}
func Test_underscoreName(t *testing.T) {
	testCases := []struct {
		name    string
		srcStr  string
		wantStr string
	}{
		// 确定 ID 只能转化为 i_d
		{
			name:    "upper cases",
			srcStr:  "ID",
			wantStr: "i_d",
		},
		{
			name:    "use number",
			srcStr:  "Table1Name",
			wantStr: "table1_name",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := underscoreName(tc.srcStr)
			assert.Equal(t, tc.wantStr, res)
		})
	}
}

func TestModelWithTableName(t *testing.T) {
	testCases := []struct {
		name          string
		val           any
		opt           ModelOpt
		wantTableName string
		wantErr       error
	}{
		{
			name:          "empty string",
			val:           &TestModel{},
			opt:           ModelWitTableName(""),
			wantTableName: "",
		},
		{
			name:          "table name",
			val:           &TestModel{},
			opt:           ModelWitTableName("test_model_table_name"),
			wantTableName: "test_model_table_name",
		},
	}
	r := NewRegistry()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.Register(tc.val, tc.opt)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				assert.Equal(t, tc.wantTableName, m.TableName)
			}
		})
	}
}

func TestModelWithColumnName(t *testing.T) {
	testCases := []struct {
		name        string
		val         any
		opt         ModelOpt
		field       string
		wantColName string
		wantErr     error
	}{
		{
			name:        "new name",
			val:         &TestModel{},
			opt:         ModelWithColumnName("FirstName", "test_first_name"),
			field:       "FirstName",
			wantColName: "test_first_name",
		},
		{
			name:        "empty new name",
			val:         &TestModel{},
			opt:         ModelWithColumnName("FirstName", ""),
			field:       "FirstName",
			wantColName: "",
		},
		{
			// 不存在的字段
			name:    "invaild field name",
			val:     &TestModel{},
			opt:     ModelWithColumnName("FirstNameTest", ""),
			field:   "FirstNameTest",
			wantErr: errs.NewErrUnknownField("FirstNameTest"),
		},
	}
	r := NewRegistry()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.Register(tc.val, tc.opt)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			fd := m.FieldMap[tc.field]
			assert.Equal(t, tc.wantColName, fd.ColName)
		})
	}
}
