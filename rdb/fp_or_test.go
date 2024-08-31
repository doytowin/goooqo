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
	"reflect"
	"testing"
)

func TestOr(t *testing.T) {

	t.Run("Build Or Condition", func(t *testing.T) {
		actual, _ := fpForOr.Process(reflect.ValueOf(&TestCond{Username: PStr("f0rb"), Email: PStr("f0rb")}))
		expect := "(username = ? OR email = ?)"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})

	t.Run("Build OR Clause for struct", func(t *testing.T) {
		query := TestQuery{TestOr: &TestCond{Username: PStr("f0rb"), Email: PStr("f0rb")}, Deleted: PBool(true)}
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
		accountAnd := TestCond{Email: PStr("f0rb@qq.com"), Mobile: PStr("01008888")}
		query := TestQuery{TestOr: &TestCond{Username: PStr("f0rb"), TestAnd: &accountAnd}, Deleted: PBool(true)}
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
		query := TestQuery{EmailEndOr: &[]string{"icloud.com", "gmail.com"}, Deleted: PBool(true)}
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
		condArr := []TestCond{{Username: PStr("f0rb")}, {Username: PStr("test2"), Email: PStr("test2@qq.com")}}
		query := TestQuery{TestsOr: &condArr, Deleted: PBool(true)}
		actual, args := BuildWhereClause(query)
		expect := " WHERE (username = ? OR username = ? AND email = ?) AND deleted = ?"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !reflect.DeepEqual(args, []any{"f0rb", "test2", "test2@qq.com", true}) {
			t.Errorf("Unexpected args: %v", args)
		}
	})

}
