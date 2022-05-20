package queryf

import (
	"testing"
	"time"

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
	t, err := time.Parse(time.RFC3339, "2022-02-10T00:00:00Z")
	suite.Nil(err)
	suite.Equal(Print(`SELECT $1`, t), `SELECT '2022-02-10T00:00:00Z'`)
	type Int64Slice []int64
	a := Int64Slice{1, 2, 3}
	suite.Equal(Print(`SELECT $1`, a), `SELECT '{1,2,3}'`)

}

func TestQueryfTestSuite(t *testing.T) {
	suite.Run(t, new(QueryfTestSuite))
}
