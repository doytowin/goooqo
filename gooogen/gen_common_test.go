/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2026, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

import (
	"os"
	"testing"
)

func TestExampleCommentMap(t *testing.T) {
	tests := []struct {
		input, expect string
		generator     Generator
	}{
		{input: "../main/inventory.go", expect: `inventory_query_builder.tpl`, generator: NewMongoGenerator()},
		{input: "../main/user.go", expect: `user_query_builder.tpl`, generator: NewSqlGenerator()},
	}
	for _, tt := range tests {
		t.Run("Generate for "+tt.input, func(t *testing.T) {
			code := GenerateCode(tt.input, tt.generator)
			expect, _ := os.ReadFile(tt.expect)
			if code != string(expect) {
				t.Fatalf("Got \n%s", code)
			}
			_ = WriteFile(tt.expect, code)
		})
	}
}
