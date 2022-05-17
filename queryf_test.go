package queryf

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type QueryfTestSuite struct {
	suite.Suite
}

func (suite *QueryfTestSuite) TestResult() {
	suite.Equal(Print(`SELECT $1, $2, $1`, 4, 5), `SELECT 4, 5, 4`)
	var arg *int64
	suite.Equal(Print(`SELECT $1`, arg), `SELECT NULL`)
	num := 5
	suite.Equal(Print(`SELECT $1`, &num), `SELECT 5`)
}

func TestQueryfTestSuite(t *testing.T) {
	suite.Run(t, new(QueryfTestSuite))
}
