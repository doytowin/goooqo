/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2026, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package core

import (
	"testing"
)

func TestBuildPageClause(t *testing.T) {
	t.Run("When PageQuery is empty then should not do paging", func(t *testing.T) {
		pageQuery := PageQuery{}
		actual := pageQuery.NeedPaging()
		if actual != false {
			t.Errorf("\nExpected: %t\nBut got : %t", false, actual)
		}
	})
	t.Run("When Size is set then should do paging", func(t *testing.T) {
		pageQuery := PageQuery{Size: 20}
		actual := pageQuery.NeedPaging()
		if actual != true {
			t.Errorf("\nExpected: %t\nBut got : %t", true, actual)
		}
	})
	t.Run("When Page is set then should do paging", func(t *testing.T) {
		pageQuery := PageQuery{Page: 1}
		actual := pageQuery.NeedPaging()
		if actual != true {
			t.Errorf("\nExpected: %t\nBut got : %t", true, actual)
		}
	})
	t.Run("When Page and Size are set then should do paging", func(t *testing.T) {
		pageQuery := PageQuery{Page: 1, Size: 10}
		actual := pageQuery.NeedPaging()
		if actual != true {
			t.Errorf("\nExpected: %t\nBut got : %t", true, actual)
		}
	})
	t.Run("CalcOffset", func(t *testing.T) {
		pageQuery := PageQuery{Page: 1, Size: 10}
		actual := pageQuery.CalcOffset()
		expect := 0
		if actual != expect {
			t.Errorf("\nExpected: %d\nBut got : %d", expect, actual)
		}
	})
	t.Run("CalcOffset", func(t *testing.T) {
		pageQuery := PageQuery{Page: 1, Size: 10}
		if pageQuery.GetSort() != "" {
			t.Errorf("\nExpected: empty string\nBut got : %s", pageQuery.GetSort())
		}
	})
}
