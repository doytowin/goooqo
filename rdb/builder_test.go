/*
 * The Clear BSD License
 *
 * Copyright (c) 2024, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package rdb

import (
	. "github.com/doytowin/goooqo/core"
	. "github.com/doytowin/goooqo/test"
	"reflect"
	"testing"
)

func TestBuildWhereClause(t *testing.T) {

	tests := []struct {
		name       string
		query      any
		expect     string
		expectArgs []any
	}{
		{
			name:       "Support custom condition",
			query:      TestQuery{Account: P("f0rb"), Deleted: P(true)},
			expect:     " WHERE (username = ? OR email = ?) AND deleted = ?",
			expectArgs: []any{"f0rb", "f0rb", true},
		},
		{
			name:       "Given field with type *bool and suffix Null, when assigned true, then map to IS NULL",
			query:      TestQuery{EmailNull: P(true)},
			expect:     " WHERE email IS NULL",
			expectArgs: []any{},
		},
		{
			name:       "Given field with type *bool and suffix Null, when assigned false, then map to IS NOT NULL",
			query:      TestQuery{EmailNull: P(false)},
			expect:     " WHERE email IS NOT NULL",
			expectArgs: []any{},
		},
		{
			name:       "Given field with type *bool and suffix Null, when not assigned, then map nothing",
			query:      TestQuery{},
			expect:     "",
			expectArgs: []any{},
		},
		{
			name:       "Given field with type *string mapped to LIKE, when assigned blank string, then map nothing",
			query:      TestQuery{EmailStart: P(" ")},
			expect:     "",
			expectArgs: []any{},
		},
		{
			"Query User by Role ID",
			UserQuery{Role: &RoleQuery{Id: P(1)}},
			" WHERE id IN (SELECT user_id FROM a_user_and_role WHERE role_id IN (SELECT id FROM t_role WHERE id = ?))",
			[]any{1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, args := BuildWhereClause(tt.query)
			if actual != tt.expect {
				t.Errorf("\nExpected: %s\nBut got : %s", tt.expect, actual)
			}
			if !reflect.DeepEqual(args, tt.expectArgs) {
				t.Errorf("BuildWhereClause() args = %v, expect %v", args, tt.expectArgs)
			}
		})
	}

}
