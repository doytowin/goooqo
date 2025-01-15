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
	. "github.com/doytowin/goooqo/core"
	"github.com/doytowin/goooqo/test"
	"reflect"
	"testing"
)

func Test_fpEntityPath_buildSql(t *testing.T) {
	epField := reflect.TypeOf(test.UserEntity{}).Field(3)
	tests := []struct {
		name  string
		field reflect.StructField
		query Query
		want  string
		want1 []any
	}{
		{
			"Build SELECT FROM t_role with conditions",
			epField,
			test.RoleQuery{Valid: P(true)},
			"SELECT id, role_name, role_code, create_user_id FROM t_role WHERE id IN (SELECT role_id FROM a_user_and_role WHERE user_id = ?) AND valid = ?",
			[]any{true},
		},
		{
			"Build SELECT FROM t_role with paging and sorting",
			epField,
			test.RoleQuery{PageQuery: PageQuery{P(10), P(5), P("role_name,desc")}, Valid: P(true)},
			"SELECT id, role_name, role_code, create_user_id FROM t_role WHERE id IN (SELECT role_id FROM a_user_and_role WHERE user_id = ?) AND valid = ? ORDER BY role_name DESC LIMIT 5 OFFSET 45",
			[]any{true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := BuildRelationEntityPath(tt.field)
			got, got1 := fp.buildSql(tt.query)
			if got != tt.want {
				t.Errorf("buildSql()\n got : %v,\n want: %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("buildSql() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
