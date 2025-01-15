/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2025, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package core

import (
	"reflect"
	"testing"
)

func TestIntId_GetId(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  int
	}{
		{"Support int", 5, 5},
		{"Support int64", int64(6), 6},
		{"Support string", "7", 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := NewIntId(0)
			id.SetId(&id, tt.input)
			if got := id.GetId(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetId() = %v, want %v", got, tt.want)
			}
		})
	}
}
