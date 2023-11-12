//go:build e2e

package integration

import (
	"github.com/NotFound1911/morm"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite

	driver string
	dsn    string

	db *morm.DB
}

func (s *Suite) SetupSuite() {
	db, err := morm.Open(s.driver, s.dsn)
	require.NoError(s.T(), err)
	err = db.Wait()
	require.NoError(s.T(), err)
	s.db = db
}
