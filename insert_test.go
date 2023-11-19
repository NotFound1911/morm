package morm

import (
	"database/sql"
	errs "github.com/NotFound1911/morm/internal/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInserter_Build(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name       string
		q          QueryBuilder
		wantQuerry *Query
		wantErr    error
	}{
		{
			name:    "no value",
			q:       NewInserter[TestModel](db).Values(),
			wantErr: errs.NewErrInsertZeroRow(),
		},
		{
			name: "single values",
			q: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        1,
					FirstName: "test",
					Age:       19,
					LastName:  &sql.NullString{String: "do", Valid: true},
				},
			),
			wantQuerry: &Query{
				SQL:  "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?);",
				Args: []any{int64(1), "test", int8(19), &sql.NullString{String: "do", Valid: true}},
			},
		},
		{
			name: "multiple values",
			q: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        1,
					FirstName: "test",
					Age:       19,
					LastName:  &sql.NullString{String: "do", Valid: true},
				},
				&TestModel{
					Id:        2,
					FirstName: "practice",
					Age:       20,
					LastName:  &sql.NullString{String: "do", Valid: true},
				},
			),
			wantQuerry: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?),(?,?,?,?);",
				Args: []any{int64(1), "test", int8(19), &sql.NullString{String: "do", Valid: true},
					int64(2), "practice", int8(20), &sql.NullString{String: "do", Valid: true}},
			},
		},
		// 指定列
		{
			name: "specify  columns",
			q: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        1,
					FirstName: "test",
					Age:       19,
					LastName:  &sql.NullString{String: "do", Valid: true},
				}).Cloumns("FirstName", "LastName"),
			wantQuerry: &Query{
				SQL:  "INSERT INTO `test_model`(`first_name`,`last_name`) VALUES(?,?);",
				Args: []any{"test", &sql.NullString{String: "do", Valid: true}},
			},
		},
		{
			name: "invalid columns",
			q: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        1,
					FirstName: "test",
					Age:       19,
					LastName:  &sql.NullString{String: "do", Valid: true},
				}).Cloumns("FirstName", "Invalid"),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
		// upsert
		{
			name: "upsert",
			q: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        1,
					FirstName: "test",
					Age:       19,
					LastName:  &sql.NullString{String: "do", Valid: true},
				}).OnDuplicateKey().Update(Assign("FirstName", "practice")),
			wantQuerry: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?) " +
					"ON DUPLICATE KEY UPDATE `first_name`=?;",
				Args: []any{int64(1), "test", int8(19), &sql.NullString{String: "do", Valid: true}, value{val: "practice"}},
			},
		},
		{
			name: "upsert invalid column",
			q: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        1,
					FirstName: "test",
					Age:       19,
					LastName:  &sql.NullString{String: "do", Valid: true},
				}).OnDuplicateKey().Update(Assign("Invalid", "invalid")),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
		{
			name: "upsert use insert value",
			q: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        1,
					FirstName: "test",
					Age:       19,
					LastName:  &sql.NullString{String: "do", Valid: true}},
				&TestModel{
					Id:        2,
					FirstName: "practice",
					Age:       19,
					LastName:  &sql.NullString{String: "te", Valid: true}},
			).OnDuplicateKey().Update(C("FirstName"), C("LastName")),
			wantQuerry: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?),(?,?,?,?) " +
					"ON DUPLICATE KEY UPDATE `first_name`=VALUES(`first_name`),`last_name`=VALUES(`last_name`);",
				Args: []any{int64(1), "test", int8(19), &sql.NullString{String: "do", Valid: true},
					int64(2), "practice", int8(19), &sql.NullString{String: "te", Valid: true}},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.q.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuerry, query)
		})
	}
}

func TestUpsert_SQLite3_Build(t *testing.T) {
	db := memoryDB(t, DBWithDialect(SQLite3))
	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			// upsert
			name: "upsert",
			q: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        1,
					FirstName: "test",
					Age:       19,
					LastName:  &sql.NullString{String: "do", Valid: true},
				}).OnDuplicateKey().ConflictColumns("Id").
				Update(Assign("FirstName", "practice")),
			wantQuery: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?) " +
					"ON CONFLICT(`id`) DO UPDATE SET `first_name`=?;",
				Args: []any{int64(1), "test", int8(19), &sql.NullString{String: "do", Valid: true}, value{val: "practice"}},
			},
		},
		{
			// upsert invalid column
			name: "upsert invalid column",
			q: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        1,
					FirstName: "test",
					Age:       20,
					LastName:  &sql.NullString{String: "do", Valid: true},
				}).OnDuplicateKey().ConflictColumns("Id").
				Update(Assign("Invalid", "")),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
		{
			// conflict invalid column
			name: "conflict invalid column",
			q: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        1,
					FirstName: "test",
					Age:       18,
					LastName:  &sql.NullString{String: "do", Valid: true},
				}).OnDuplicateKey().ConflictColumns("Invalid").
				Update(Assign("FirstName", "practice")),
			wantErr: errs.NewErrUnknownField("Invalid"),
		},
		{
			// 使用原本插入的值
			name: "upsert use insert value",
			q: NewInserter[TestModel](db).Values(
				&TestModel{
					Id:        1,
					FirstName: "test",
					Age:       18,
					LastName:  &sql.NullString{String: "do", Valid: true},
				},
				&TestModel{
					Id:        2,
					FirstName: "practice",
					Age:       19,
					LastName:  &sql.NullString{String: "did", Valid: true},
				}).OnDuplicateKey().ConflictColumns("Id").
				Update(C("FirstName"), C("LastName")),
			wantQuery: &Query{
				SQL: "INSERT INTO `test_model`(`id`,`first_name`,`age`,`last_name`) VALUES(?,?,?,?),(?,?,?,?) " +
					"ON CONFLICT(`id`) DO UPDATE SET `first_name`=excluded.`first_name`,`last_name`=excluded.`last_name`;",
				Args: []any{int64(1), "test", int8(18), &sql.NullString{String: "do", Valid: true},
					int64(2), "practice", int8(19), &sql.NullString{String: "did", Valid: true}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, err := tc.q.Build()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}
