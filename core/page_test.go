package core

import (
	"testing"
)

func TestBuildPageClause(t *testing.T) {
	t.Run("Build Page Clause", func(t *testing.T) {
		pageQuery := PageQuery{PInt(3), PInt(10)}
		actual := pageQuery.BuildPageClause()
		expect := " LIMIT 10 OFFSET 20"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})
	t.Run("When PageQuery is empty then should not do paging", func(t *testing.T) {
		pageQuery := PageQuery{}
		actual := pageQuery.NeedPaging()
		if actual != false {
			t.Errorf("\nExpected: %t\nBut got : %t", false, actual)
		}
	})
	t.Run("When PageSize is set then should do paging", func(t *testing.T) {
		pageQuery := PageQuery{PageSize: PInt(20)}
		actual := pageQuery.NeedPaging()
		if actual != true {
			t.Errorf("\nExpected: %t\nBut got : %t", true, actual)
		}
	})
	t.Run("When PageNumber is set then should do paging", func(t *testing.T) {
		pageQuery := PageQuery{PageNumber: PInt(1)}
		actual := pageQuery.NeedPaging()
		if actual != true {
			t.Errorf("\nExpected: %t\nBut got : %t", true, actual)
		}
	})
	t.Run("When PageNumber and PageSize are set then should do paging", func(t *testing.T) {
		pageQuery := PageQuery{PageNumber: PInt(1), PageSize: PInt(10)}
		actual := pageQuery.NeedPaging()
		if actual != true {
			t.Errorf("\nExpected: %t\nBut got : %t", true, actual)
		}
	})
	t.Run("Build Page Clause with PageNumber Only", func(t *testing.T) {
		pageQuery := PageQuery{PageNumber: PInt(3)}
		actual := pageQuery.BuildPageClause()
		expect := " LIMIT 10 OFFSET 20"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})
	t.Run("Build Page Clause with PageSize Only", func(t *testing.T) {
		pageQuery := PageQuery{PageSize: PInt(20)}
		actual := pageQuery.BuildPageClause()
		expect := " LIMIT 20 OFFSET 0"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})
	t.Run("Build Page Clause with PageNumber less than 1", func(t *testing.T) {
		pageQuery := PageQuery{PageNumber: PInt(0)}
		actual := pageQuery.BuildPageClause()
		expect := " LIMIT 10 OFFSET 0"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})
	t.Run("Build Page Clause with PageSize less than 0", func(t *testing.T) {
		pageQuery := PageQuery{PageSize: PInt(-1)}
		actual := pageQuery.BuildPageClause()
		expect := " LIMIT 10 OFFSET 0"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})
}
