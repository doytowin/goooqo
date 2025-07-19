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
	"reflect"
	"testing"
)

func TestOr(t *testing.T) {

	t.Run("Build Or Condition", func(t *testing.T) {
		actual, _ := fpForOr.Process(reflect.ValueOf(&TestQuery{Username: P("f0rb"), Email: P("f0rb")}))
		expect := "(username = ? OR email = ?)"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})

	t.Run("Build OR Clause for struct", func(t *testing.T) {
		query := TestQuery{Or: &TestQuery{Username: P("f0rb"), Email: P("f0rb")}, Deleted: P(true)}
		actual, args := BuildWhereClause(query)
		expect := " WHERE (username = ? OR email = ?) AND deleted = ?"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !reflect.DeepEqual(args, []any{"f0rb", "f0rb", true}) {
			t.Errorf("Unexpected args: %v", args)
		}
	})

	t.Run("Build OR Clause with And", func(t *testing.T) {
		accountAnd := TestQuery{Email: P("f0rb@qq.com"), Mobile: P("01008888")}
		query := TestQuery{Or: &TestQuery{Username: P("f0rb"), And: &accountAnd}, Deleted: P(true)}
		actual, args := BuildWhereClause(query)
		expect := " WHERE (username = ? OR email = ? AND mobile = ?) AND deleted = ?"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !reflect.DeepEqual(args, []any{"f0rb", "f0rb@qq.com", "01008888", true}) {
			t.Errorf("Unexpected args: %v", args)
		}
	})

	t.Run("Build OR Clause for basic array", func(t *testing.T) {
		query := TestQuery{EmailEndOr: &[]string{"icloud.com", "gmail.com"}, Deleted: P(true)}
		actual, args := BuildWhereClause(query)
		expect := " WHERE (email LIKE ? OR email LIKE ?) AND deleted = ?"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !reflect.DeepEqual(args, []any{"%icloud.com", "%gmail.com", true}) {
			t.Errorf("Unexpected args: %v", args)
		}
	})

	t.Run("Build OR Clause for struct array", func(t *testing.T) {
		condArr := []TestQuery{{Username: P("f0rb")}, {Username: P("test2"), Email: P("test2@qq.com")}}
		query := TestQuery{TestsOr: &condArr, Deleted: P(true)}
		actual, args := BuildWhereClause(query)
		expect := " WHERE (username = ? OR username = ? AND email = ?) AND deleted = ?"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !reflect.DeepEqual(args, []any{"f0rb", "test2", "test2@qq.com", true}) {
			t.Errorf("Unexpected args: %v", args)
		}
	})

	t.Run("Build for named Or", func(t *testing.T) {
		condArr := TestQuery{EmailStart: P("test"), EmailNull: P(true)}
		query := TestQuery{Or: &condArr, Deleted: P(true)}
		actual, args := BuildWhereClause(query)
		expect := " WHERE (email LIKE ? OR email IS NULL) AND deleted = ?"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !reflect.DeepEqual(args, []any{"test%", true}) {
			t.Errorf("Unexpected args: %v", args)
		}
	})

}
