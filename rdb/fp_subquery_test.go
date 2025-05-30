/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2025, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package rdb

import (
	"testing"
)

func TestTrimFieldName(t *testing.T) {
	if got := trimFieldName("ScoreLtAvg"); got != "ScoreLt" {
		t.Errorf("trimFieldName() = %v, want %v", got, "ScoreLt")
	}
}
