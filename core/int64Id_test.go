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
	"reflect"
	"testing"
)

func TestInt64Id_GetId(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  int64
	}{
		{"Support int", 5, int64(5)},
		{"Support int64", int64(6), int64(6)},
		{"Support string", "7", int64(7)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := NewInt64Id(0)
			id.SetId(&id, tt.input)
			if got := id.GetId(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetId() = %v, want %v", got, tt.want)
			}
		})
	}
}
