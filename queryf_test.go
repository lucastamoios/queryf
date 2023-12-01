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
	suite.Equal(Format(`SELECT $1, $2, $1`, 4, 5), `SELECT 4, 5, 4`)
	var arg *int64
	suite.Equal(Format(`SELECT $1`, arg), `SELECT NULL`)
	num := 5
	suite.Equal(Format(`SELECT $1`, &num), `SELECT 5`)
	t, err := time.Parse(time.RFC3339, "2022-02-10T00:00:00Z")
	suite.Nil(err)
	suite.Equal(Format(`SELECT $1`, t), `SELECT '2022-02-10T00:00:00Z'`)
	a := []int64{1, 2, 3}
	suite.Equal(Format(`SELECT $1`, a), `SELECT '{1,2,3}'`)

	// List with nil values
	v := int64(5)
	b := []*int64{&v, nil}
	suite.Equal(Format(`SELECT $1`, b), `SELECT '{5,NULL}'`)
}

func TestQueryfTestSuite(t *testing.T) {
	suite.Run(t, new(QueryfTestSuite))
}
