/*
 * The Clear BSD License
 *
 * Copyright (c) 2024, DoytoWin, Inc.
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
}
