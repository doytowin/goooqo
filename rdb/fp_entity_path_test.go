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

func Test_fpEntityPath_buildQuery(t *testing.T) {
	epField := reflect.TypeOf(test.UserEntity{}).Field(3)
	tests := []struct {
		name  string
		field reflect.StructField
		query Query
		sql   string
		args  []any
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
			sql, args := fp.buildQuery(tt.query)
			if sql != tt.sql {
				t.Errorf("buildSql()\n got : %v\n want: %v", sql, tt.sql)
			}
			if !reflect.DeepEqual(args, tt.args) {
				t.Errorf("buildSql()\n got : %v\n want: %v", args, tt.args)
			}
		})
	}
}

func Test_fpEntityPath_buildSql(t *testing.T) {
	RegisterJoinTable("product", "order", "a_order_and_product")
	RegisterVirtualEntity("friend", "user")
	tests := []struct {
		name string
		aep  string
		sql  string
	}{
		{
			"Build SELECT for user's products",
			"user,user_id<-order,product",
			"SELECT * FROM t_product WHERE id IN (SELECT product_id FROM a_order_and_product WHERE order_id IN (SELECT id FROM t_order WHERE user_id = ?))",
		},
		{
			"Build SELECT for product's buyer",
			"product,order->user_id,user",
			"SELECT * FROM t_user WHERE id IN (SELECT user_id FROM t_order WHERE id IN (SELECT order_id FROM a_order_and_product WHERE product_id = ?))",
		},
		{
			"Support user's friend",
			"user,friend,friend,friend",
			"SELECT * FROM t_user WHERE id IN (SELECT friend_id FROM a_user_and_friend WHERE user_id IN (SELECT friend_id FROM a_user_and_friend WHERE user_id IN (SELECT friend_id FROM a_user_and_friend WHERE user_id = ?)))",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := &fpEntityPath{*BuildEntityPathStr(tt.aep)}
			sql := fp.buildSql("*")
			if sql != tt.sql {
				t.Errorf("buildSql()\n got : %v\n want: %v", sql, tt.sql)
			}
		})
	}
}
