/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2026, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package rdb

import (
	"testing"

	. "github.com/doytowin/goooqo/core"
)

func TestBuildPageClause(t *testing.T) {
	t.Run("Build Page Clause", func(t *testing.T) {
		pageQuery := PageQuery{PageNumber: P(3), PageSize: P(10)}
		actual := Dialect.BuildPageClause("", pageQuery.CalcOffset(), pageQuery.GetPageSize())
		expect := " LIMIT 10 OFFSET 20"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})
	t.Run("Build Page Clause with PageNumber Only", func(t *testing.T) {
		pageQuery := PageQuery{PageNumber: P(3)}
		actual := Dialect.BuildPageClause("", pageQuery.CalcOffset(), pageQuery.GetPageSize())
		expect := " LIMIT 10 OFFSET 20"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})
	t.Run("Build Page Clause with PageSize Only", func(t *testing.T) {
		pageQuery := PageQuery{PageSize: P(20)}
		actual := Dialect.BuildPageClause("", pageQuery.CalcOffset(), pageQuery.GetPageSize())
		expect := " LIMIT 20 OFFSET 0"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})
	t.Run("Build Page Clause with PageNumber less than 1", func(t *testing.T) {
		pageQuery := PageQuery{PageNumber: P(0)}
		actual := Dialect.BuildPageClause("", pageQuery.CalcOffset(), pageQuery.GetPageSize())
		expect := " LIMIT 10 OFFSET 0"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})
	t.Run("Build Page Clause with PageSize less than 0", func(t *testing.T) {
		pageQuery := PageQuery{PageSize: P(-1)}
		actual := Dialect.BuildPageClause("", pageQuery.CalcOffset(), pageQuery.GetPageSize())
		expect := " LIMIT 10 OFFSET 0"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})
	t.Run("Build Sort Clause", func(t *testing.T) {
		pageQuery := PageQuery{Sort: P("id,desc;score,asc;memo")}
		actual := BuildSortClause(pageQuery.GetSort())
		expect := " ORDER BY id DESC, score ASC, memo"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})
}
