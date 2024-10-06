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
		name   string
		query  any
		expect string
		args   []any
	}{
		{
			"Support custom condition",
			TestQuery{Account: P("f0rb"), Deleted: P(true)},
			" WHERE (username = ? OR email = ?) AND deleted = ?",
			[]any{"f0rb", "f0rb", true},
		},
		{
			"Given field with type *bool and suffix Null, when assigned true, then map to IS NULL",
			TestQuery{EmailNull: P(true)},
			" WHERE email IS NULL",
			[]any{},
		},
		{
			"Given field with type *bool and suffix Null, when assigned false, then map to IS NOT NULL",
			TestQuery{EmailNull: P(false)},
			" WHERE email IS NOT NULL",
			[]any{},
		},
		{
			"Given field with type *bool and suffix Null, when not assigned, then map nothing",
			TestQuery{},
			"",
			[]any{},
		},
		{
			"Given field with type *string mapped to LIKE, when assigned blank string, then map nothing",
			TestQuery{EmailStart: P(" ")},
			"",
			[]any{},
		},
		{
			"Query User by Role ID",
			UserQuery{Role: &RoleQuery{Id: P(1)}},
			" WHERE id IN (SELECT user_id FROM a_user_and_role WHERE role_id IN (SELECT id FROM t_role WHERE id = ?))",
			[]any{1},
		},
		{
			"Query User by Permission id",
			UserQuery{Perm: &PermQuery{Code: P("user:list")}},
			" WHERE id IN (SELECT user_id FROM a_user_and_role WHERE role_id IN " +
				"(SELECT role_id FROM a_role_and_perm WHERE perm_id IN (SELECT id FROM t_perm WHERE code = ?)))",
			[]any{"user:list"},
		},
		{
			"Query User by valid Role and Permission id",
			UserQuery{Perm: &PermQuery{Code: P("user:list"), RoleQuery: &RoleQuery{Valid: P(true)}}},
			` WHERE id IN (SELECT user_id FROM a_user_and_role WHERE role_id IN (SELECT id FROM t_role WHERE valid = ?
INTERSECT SELECT role_id FROM a_role_and_perm WHERE perm_id IN (SELECT id FROM t_perm WHERE code = ?)))`,
			[]any{true, "user:list"},
		},
		{
			"Query Role by User id",
			RoleQuery{User: &UserQuery{IdIn: &[]int{1, 3, 4}}},
			" WHERE id IN (SELECT role_id FROM a_user_and_role WHERE user_id IN (SELECT id FROM t_user WHERE id IN (?, ?, ?)))",
			[]any{1, 3, 4},
		},
		{
			"Query children menu by parent id | one-to-many",
			MenuQuery{Parent: &MenuQuery{Id: P(1)}},
			" WHERE parent_id IN (SELECT id FROM t_menu WHERE id = ?)",
			[]any{1},
		},
		{
			"Query parent menu by child id | many-to-one",
			MenuQuery{Children: &MenuQuery{Id: P(1)}},
			" WHERE id IN (SELECT parent_id FROM t_menu WHERE id = ?)",
			[]any{1},
		},
	}
	RegisterJoinTable("role", "user", "a_user_and_role")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, args := BuildWhereClause(tt.query)
			if actual != tt.expect {
				t.Errorf("\nExpected: %s\nBut got : %s", tt.expect, actual)
			}
			if !reflect.DeepEqual(args, tt.args) {
				t.Errorf("BuildWhereClause() args = %v, expect %v", args, tt.args)
			}
		})
	}

}
