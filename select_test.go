package morm

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
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
	return "test_model"
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
				SQL: "SELECT * FROM `test_model`;",
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
				SQL: "SELECT * FROM `test_model`;",
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
				SQL:  "SELECT * FROM `test_model` WHERE (`age` > ?) AND (`age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			name: "use and",
			q:    NewSelector[TestModel](db).Where(C("Age").GT(18).And(C("Age").LT(35))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` > ?) AND (`age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			name: "use or",
			q:    NewSelector[TestModel](db).Where(C("Age").GT(18).Or(C("Age").LT(35))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE (`age` > ?) OR (`age` < ?);",
				Args: []any{18, 35},
			},
		},
		{
			name: "use not",
			q:    NewSelector[TestModel](db).Where(Not(C("Age").GT(18))),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE  NOT (`age` > ?);",
				Args: []any{18},
			},
		},
		{
			name:    "invalid column",
			q:       NewSelector[TestModel](db).Where(Not(C("Invalid").GT(18))),
			wantErr: errs.NewErrUnknownField("Invalid"),
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

func TestSelector_Get(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = mockDB.Close() }()
	db, err := OpenDB(mockDB)
	if err != nil {
		t.Fatal(err)
	}
	testCases := []struct {
		name     string
		query    string
		mockErr  error
		mockRows *sqlmock.Rows
		wantErr  error
		wantVal  *TestModel
	}{
		{
			name:    "query error",
			mockErr: errors.New("invalid error"),
			wantErr: errors.New("invalid error"),
			query:   "SELECT .*",
		},
		{
			name:     "no row",
			wantErr:  sql.ErrNoRows,
			query:    "SELECT .*",
			mockRows: sqlmock.NewRows([]string{"id"}),
		},
		{
			name:    "too many column",
			wantErr: errs.NewErrTooManyReturnedColumns([]string{"id", "first_name", "age", "last_name", "extra_column"}),
			mockRows: func() *sqlmock.Rows {
				res := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name", "extra_column"})
				res.AddRow([]byte("1"), []byte("Da"), []byte("18"), []byte("Ming"), []byte("nothing"))
				return res
			}(),
		},
		{
			name:  "get data",
			query: "SELECT .*",
			mockRows: func() *sqlmock.Rows {
				res := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
				res.AddRow([]byte("1"), []byte("Da"), []byte("18"), []byte("Ming"))
				return res
			}(),
			wantVal: &TestModel{
				Id:        1,
				FirstName: "Da",
				Age:       18,
				LastName:  &sql.NullString{String: "Ming", Valid: true},
			},
		},
	}
	for _, tc := range testCases {
		exp := mock.ExpectQuery(tc.query)
		if tc.mockErr != nil {
			exp.WillReturnError(tc.mockErr)
		} else {
			exp.WillReturnRows(tc.mockRows)
		}
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := NewSelector[TestModel](db).Get(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, res)
		})
	}
}

func TestSelector_GetMulti(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = mockDB.Close() }()
	db, err := OpenDB(mockDB)
	if err != nil {
		t.Fatal(err)
	}
	testCases := []struct {
		name     string
		query    string
		mockErr  error
		mockRows *sqlmock.Rows
		wantErr  error
		wantVal  []*TestModel
	}{
		{
			name:    "query error",
			mockErr: errors.New("invalid error"),
			wantErr: errors.New("invalid error"),
			query:   "SELECT .*",
		},
		{
			name:     "no row",
			wantErr:  sql.ErrNoRows,
			query:    "SELECT .*",
			mockRows: sqlmock.NewRows([]string{"id"}),
		},
		{
			name:    "too many column",
			wantErr: errs.NewErrTooManyReturnedColumns([]string{"id", "first_name", "age", "last_name", "extra_column"}),
			mockRows: func() *sqlmock.Rows {
				res := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name", "extra_column"})
				res.AddRow([]byte("1"), []byte("Da"), []byte("18"), []byte("Ming"), []byte("nothing"))
				return res
			}(),
		},
		{
			name:  "get data",
			query: "SELECT .*",
			mockRows: func() *sqlmock.Rows {
				res := sqlmock.NewRows([]string{"id", "first_name", "age", "last_name"})
				res.AddRow([]byte("1"), []byte("Da"), []byte("18"), []byte("Ming"))
				res.AddRow([]byte("2"), []byte("Hao"), []byte("19"), []byte("Yang"))
				res.AddRow([]byte("3"), []byte("Niu"), []byte("18"), []byte("Zero"))
				return res
			}(),

			wantVal: []*TestModel{
				{
					Id:        1,
					FirstName: "Da",
					Age:       18,
					LastName:  &sql.NullString{String: "Ming", Valid: true},
				},
				{
					Id:        2,
					FirstName: "Hao",
					Age:       19,
					LastName:  &sql.NullString{String: "Yang", Valid: true},
				},
				{
					Id:        3,
					FirstName: "Niu",
					Age:       18,
					LastName:  &sql.NullString{String: "Zero", Valid: true},
				},
			},
		},
	}
	for _, tc := range testCases {
		exp := mock.ExpectQuery(tc.query)
		if tc.mockErr != nil {
			exp.WillReturnError(tc.mockErr)
		} else {
			exp.WillReturnRows(tc.mockRows)
		}
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := NewSelector[TestModel](db).GetMulti(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantVal, res)
		})
	}
}

func TestSelector_Select(t *testing.T) {
	db := memoryDB(t)
	testCases := []struct {
		name      string
		q         QueryBuilder
		wantQuery *Query
		wantErr   error
	}{
		{
			name: "all",
			q:    NewSelector[TestModel](db),
			wantQuery: &Query{
				SQL: "SELECT * FROM `test_model`;",
			},
		},
		{
			name:    "invalid column",
			q:       NewSelector[TestModel](db).Select(Avg("invalid")),
			wantErr: errs.NewErrUnknownField("invalid"),
		},
		{
			name: "partial columns",
			q:    NewSelector[TestModel](db).Select(C("Id"), C("FirstName")),
			wantQuery: &Query{
				SQL: "SELECT `id`,`first_name` FROM `test_model`;",
			},
		},
		{
			name: "avg fn",
			q:    NewSelector[TestModel](db).Select(Avg("Age")),
			wantQuery: &Query{
				SQL: "SELECT AVG(`age`) FROM `test_model`;",
			},
		},
		{
			name: "raw expression",
			q:    NewSelector[TestModel](db).Select(Raw("COUNT(DISTINCT `first_name`)")),
			wantQuery: &Query{
				SQL: "SELECT COUNT(DISTINCT `first_name`) FROM `test_model`;",
			},
		},
		// 别名
		{
			name: "alias",
			q:    NewSelector[TestModel](db).Select(C("Id").As("my_id"), Avg("Age").As("avg_age")),
			wantQuery: &Query{
				SQL: "SELECT `id` AS `my_id`,AVG(`age`) AS `avg_age` FROM `test_model`;",
			},
		},
		// where 忽略别名
		{
			name: "where ignore alias",
			q:    NewSelector[TestModel](db).Where(C("Id").As("my_id").LT(1001)),
			wantQuery: &Query{
				SQL:  "SELECT * FROM `test_model` WHERE `id` < ?;",
				Args: []any{1001},
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
