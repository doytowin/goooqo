/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2026, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package mongodb

import (
	"reflect"
	"testing"

	. "go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_buildSort(t *testing.T) {
	tests := []struct {
		input  string
		expect D
	}{
		{"item,desc", D{{"item", -1}}},
		{"item,asc", D{{"item", 1}}},
		{"item", D{{"item", 1}}},
		{"item,desc;qty,asc", D{{"item", -1}, {"qty", 1}}},
		{"item;qty,asc", D{{"item", 1}, {"qty", 1}}},
	}
	for _, tt := range tests {
		t.Run("Sort:"+tt.input, func(t *testing.T) {
			if got := buildSort(tt.input); !reflect.DeepEqual(got, tt.expect) {
				t.Errorf("buildSort() = %v, want %v", got, tt.expect)
			}
		})
	}
}
