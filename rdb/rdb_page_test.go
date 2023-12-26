package rdb

import (
	. "github.com/doytowin/goooqo/core"
	"testing"
)

func TestBuildPageClause(t *testing.T) {
	t.Run("Build Page Clause", func(t *testing.T) {
		pageQuery := PageQuery{PageNumber: PInt(3), PageSize: PInt(10)}
		actual := BuildPageClause(PStr(""), pageQuery.CalcOffset(), pageQuery.GetPageSize())
		expect := " LIMIT 10 OFFSET 20"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})
	t.Run("Build Page Clause with PageNumber Only", func(t *testing.T) {
		pageQuery := PageQuery{PageNumber: PInt(3)}
		actual := BuildPageClause(PStr(""), pageQuery.CalcOffset(), pageQuery.GetPageSize())
		expect := " LIMIT 10 OFFSET 20"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})
	t.Run("Build Page Clause with PageSize Only", func(t *testing.T) {
		pageQuery := PageQuery{PageSize: PInt(20)}
		actual := BuildPageClause(PStr(""), pageQuery.CalcOffset(), pageQuery.GetPageSize())
		expect := " LIMIT 20 OFFSET 0"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})
	t.Run("Build Page Clause with PageNumber less than 1", func(t *testing.T) {
		pageQuery := PageQuery{PageNumber: PInt(0)}
		actual := BuildPageClause(PStr(""), pageQuery.CalcOffset(), pageQuery.GetPageSize())
		expect := " LIMIT 10 OFFSET 0"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})
	t.Run("Build Page Clause with PageSize less than 0", func(t *testing.T) {
		pageQuery := PageQuery{PageSize: PInt(-1)}
		actual := BuildPageClause(PStr(""), pageQuery.CalcOffset(), pageQuery.GetPageSize())
		expect := " LIMIT 10 OFFSET 0"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})
	t.Run("Build Sort Clause", func(t *testing.T) {
		pageQuery := PageQuery{Sort: PStr("id,desc;score,asc;memo")}
		actual := BuildSortClause(pageQuery.GetSort())
		expect := " ORDER BY id DESC, score ASC, memo"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})
}
