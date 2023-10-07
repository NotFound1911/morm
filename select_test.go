package morm

import (
	"database/sql"
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
func TestSelector_Build(t *testing.T) {
	db, err := NewDB()
	if err != nil {
		t.Fatal(err)
	}
	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name: "no from",
			q:    NewSelector[TestModel](db),
			wantQuery: &Query{
				SQL: "SELECT * FROM `TestModel`;",
			},
		},
		{
			name: "with from",
			q:    NewSelector[TestModel](db).From("`test_model_t`"),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model_t`;",
			},
		},
		{
			name: "empty from",
			q:    NewSelector[TestModel](db).From(""),
			wantQuery: &Query{
				SQL: "SELECT * FROM `TestModel`;",
			},
		},
		{
			name: "with db",
			q:    NewSelector[TestModel](db).From("`test_db`.`test_model`"),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_db`.`test_model`;",
			},
		},
		{
			name: "single and simple predicate",
			q:    NewSelector[TestModel](db).From("`test_model_t`").Where(C("Id").EQ(1)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model_t` WHERE `id` = ?;",
				Args: []any{1},
			},
		},
		{
			name: "multi predicates",
			q:    NewSelector[TestModel](db).Where(C("Age").GT(18), C("Age").LT(35)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel` WHERE (`age` > ?) AND (`age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			name: "use and",
			q:    NewSelector[TestModel](db).Where(C("Age").GT(18).And(C("Age").LT(35))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel` WHERE (`age` > ?) AND (`age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			name: "use or",
			q:    NewSelector[TestModel](db).Where(C("Age").GT(18).Or(C("Age").LT(35))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel` WHERE (`age` > ?) OR (`age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			name: "use not",
			q:    NewSelector[TestModel](db).Where(Not(C("Age").GT(18))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel` WHERE  NOT (`age` > ?);",
				Args: []any{18},
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
