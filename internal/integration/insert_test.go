//go:build e2e

package integration

import (
	"context"
	"github.com/NotFound1911/morm"
	"github.com/NotFound1911/morm/internal/test"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type InsertTestSuite struct {
	Suite
}

func (i *InsertTestSuite) TearDownTest() {
	res := morm.RawQuery[any](i.db, "TRUNCATE TABLE `simple_struct`").Exec(context.Background())
	require.NoError(i.T(), res.Err())
}
func TestInsertMySQL8(t *testing.T) {
	suite.Run(t, &InsertTestSuite{
		Suite: Suite{
			driver: "mysql",
			dsn:    "root:root@tcp(localhost:13306)/integration_test",
		},
	})
}

func (i *InsertTestSuite) TestInsert() {
	testCases := []struct {
		name         string
		i            *morm.Inserter[test.SimpleStruct]
		wantData     *test.SimpleStruct
		rowsAffected int64
		wantErr      error
	}{
		{
			name:         "id only",
			i:            morm.NewInserter[test.SimpleStruct](i.db).Values(&test.SimpleStruct{Id: 1}),
			rowsAffected: 1,
			wantData:     &test.SimpleStruct{Id: 1},
		},
		{
			name:         "all field",
			i:            morm.NewInserter[test.SimpleStruct](i.db).Values(test.NewSimpleStruct(2)),
			rowsAffected: 1,
			wantData:     test.NewSimpleStruct(2),
		},
	}
	for _, tc := range testCases {
		i.T().Run(tc.name, func(t *testing.T) {
			res := tc.i.Exec(context.Background())
			affected, err := res.RowsAffected()
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.rowsAffected, affected)
			data, err := morm.NewSelector[test.SimpleStruct](i.db).
				Where(morm.C("Id").EQ(tc.wantData.Id)).Get(context.Background())
			require.NoError(t, err)
			assert.Equal(t, tc.wantData, data)
		})
	}
}
