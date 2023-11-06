package querylog

import (
	"context"
	"database/sql"
	"github.com/NotFound1911/morm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// 测试打印效果
func TestMiddlewareBuilder(t *testing.T) {
	var query string // sql
	var args []any   // 参数
	m := NewMiddlewareBuilder().LogFunc(func(q string, as []any) {
		query = q
		args = as
	})
	db, err := morm.Open("sqlite3",
		"file:test.db?cache=share&mode=memory",
		morm.DBWithMiddleware(m.Build()))
	require.NoError(t, err)
	_, _ = morm.NewSelector[TestModel](db).Where(morm.C("Id").EQ(10)).Get(context.Background())
	assert.Equal(t, "SELECT * FROM `test_model` WHERE `id` = ?;", query)
	assert.Equal(t, []any{10}, args)

	morm.NewInserter[TestModel](db).Values(&TestModel{Id: 10}).Exec(context.Background())
	assert.Equal(t, "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?);", query)
	assert.Equal(t, []any{int64(10), "", int8(0), (*sql.NullString)(nil)}, args)

}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}
