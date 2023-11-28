package goquery

import (
	"github.com/doytowin/doyto-query-go-sql/util"
	"testing"
)

func TestBuildPageClause(t *testing.T) {
	t.Run("Build Page Clause", func(t *testing.T) {
		pageQuery := PageQuery{util.PInt(3), util.PInt(10)}
		actual := pageQuery.buildPageClause()
		expect := " LIMIT 10 OFFSET 20"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})
}
