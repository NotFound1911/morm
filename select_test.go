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

func TestSelector_Build(t *testing.T) {
	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name: "no from",
			q:    NewSelector[TestModel](),
			wantQuery: &Query{
				SQL: "SELECT * FROM `TestModel`;",
			},
		},
		{
			name: "with from",
			q:    NewSelector[TestModel]().From("`test_model_t`"),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model_t`;",
			},
		},
		{
			name: "empty from",
			q:    NewSelector[TestModel]().From(""),
			wantQuery: &Query{
				SQL: "SELECT * FROM `TestModel`;",
			},
		},
		{
			name: "with db",
			q:    NewSelector[TestModel]().From("`test_db`.`test_model`"),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_db`.`test_model`;",
			},
		},
		{
			name: "single and simple predicate",
			q:    NewSelector[TestModel]().From("`test_model_t`").Where(C("Id").EQ(1)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model_t` WHERE `Id` = ?;",
				Args: []any{1},
			},
		},
		{
			name: "multi predicates",
			q:    NewSelector[TestModel]().Where(C("Age").GT(18), C("Age").LT(35)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel` WHERE (`Age` > ?) AND (`Age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			name: "use and",
			q:    NewSelector[TestModel]().Where(C("Age").GT(18).And(C("Age").LT(35))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel` WHERE (`Age` > ?) AND (`Age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			name: "use or",
			q:    NewSelector[TestModel]().Where(C("Age").GT(18).Or(C("Age").LT(35))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel` WHERE (`Age` > ?) OR (`Age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			name: "use not",
			q:    NewSelector[TestModel]().Where(Not(C("Age").GT(18))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `TestModel` WHERE  NOT (`Age` > ?);",
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
