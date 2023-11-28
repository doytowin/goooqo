package goquery

import (
	. "github.com/doytowin/goquery/util"
	"testing"
)

func TestBuildPageClause(t *testing.T) {
	t.Run("Build Page Clause", func(t *testing.T) {
		pageQuery := PageQuery{PInt(3), PInt(10)}
		actual := pageQuery.buildPageClause()
		expect := " LIMIT 10 OFFSET 20"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})
	t.Run("When PageQuery is empty then should not do paging", func(t *testing.T) {
		pageQuery := PageQuery{}
		actual := pageQuery.needPaging()
		if actual != false {
			t.Errorf("\nExpected: %t\nBut got : %t", false, actual)
		}
	})
	t.Run("When PageSize is set then should do paging", func(t *testing.T) {
		pageQuery := PageQuery{PageSize: PInt(20)}
		actual := pageQuery.needPaging()
		if actual != true {
			t.Errorf("\nExpected: %t\nBut got : %t", true, actual)
		}
	})
	t.Run("When PageNumber is set then should do paging", func(t *testing.T) {
		pageQuery := PageQuery{PageNumber: PInt(1)}
		actual := pageQuery.needPaging()
		if actual != true {
			t.Errorf("\nExpected: %t\nBut got : %t", true, actual)
		}
	})
	t.Run("When PageNumber and PageSize are set then should do paging", func(t *testing.T) {
		pageQuery := PageQuery{PageNumber: PInt(1), PageSize: PInt(10)}
		actual := pageQuery.needPaging()
		if actual != true {
			t.Errorf("\nExpected: %t\nBut got : %t", true, actual)
		}
	})
}
