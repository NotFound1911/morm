//go:build e2e

package integration

import (
	"context"
	"database/sql"
	"github.com/NotFound1911/morm"
	"github.com/NotFound1911/morm/internal/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SelectTestSuite struct {
	Suite
}

func (s *SelectTestSuite) SetupSuite() {
	s.Suite.SetupSuite()
	res := morm.NewInserter[test.SimpleStruct](s.db).
		Values(test.NewSimpleStruct(1), test.NewSimpleStruct(2)).Exec(context.Background())
	require.NoError(s.T(), res.Err())
}

func (s *SelectTestSuite) TearDownSuite() {
	res := morm.RawQuery[any](s.db, "TRUNCATE TABLE `simple_struct`").Exec(context.Background())
	require.NoError(s.T(), res.Err())
}

func (s *SelectTestSuite) TestSet() {
	testCase := []struct {
		name    string
		s       *morm.Selector[test.SimpleStruct]
		wantErr error
		wantRes *test.SimpleStruct
	}{
		{
			name:    "not found",
			s:       morm.NewSelector[test.SimpleStruct](s.db).Where(morm.C("Id").EQ(18)),
			wantErr: sql.ErrNoRows,
		},
		{
			name:    "find",
			s:       morm.NewSelector[test.SimpleStruct](s.db).Where(morm.C("Id").EQ(1)),
			wantRes: test.NewSimpleStruct(1),
		},
	}

	for _, tc := range testCase {
		s.T().Run(tc.name, func(t *testing.T) {
			res, err := tc.s.Get(context.Background())
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}

}
func TestSelectMySQL8(t *testing.T) {
	suite.Run(t, &SelectTestSuite{
		Suite: Suite{
			driver: "mysql",
			dsn:    "root:root@tcp(localhost:13306)/integration_test",
		},
	})
}
