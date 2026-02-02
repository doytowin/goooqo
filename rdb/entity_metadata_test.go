/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2026, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package rdb

import (
	"reflect"
	"testing"

	. "github.com/doytowin/goooqo/core"
	. "github.com/doytowin/goooqo/test"
	log "github.com/sirupsen/logrus"
)

func TestBuildStmt(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	em := buildEntityMetadata[UserEntity]()

	t.Run("Build with Custom Table Name", func(t *testing.T) {
		em := buildEntityMetadata[TestEntity]()
		actual := em.TableName
		expect := "t_user"
		if actual != expect {
			t.Errorf("\nExpected: %s\n     Got: %s", expect, actual)
		}
	})

	t.Run("Support snake_case_column", func(t *testing.T) {
		em := buildEntityMetadata[TestEntity]()
		actual := em.ColStr
		expect := "id, username, email, mobile, create_time"
		if actual != expect {
			t.Errorf("\nExpected: %s\n     Got: %s", expect, actual)
		}
	})

	t.Run("Build Where Clause", func(t *testing.T) {
		query := UserQuery{IdGt: P(5), MemoNull: P(true)}
		actual, args := BuildWhereClause(query)
		expect := " WHERE id > ? AND memo IS NULL"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !reflect.DeepEqual(args, []any{5}) {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Build Where Clause", func(t *testing.T) {
		query := UserQuery{IdGt: P(5), MemoNull: P(false)}
		actual, args := BuildWhereClause(query)
		expect := " WHERE id > ? AND memo IS NOT NULL"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !reflect.DeepEqual(args, []any{5}) {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Build Select Statement", func(t *testing.T) {
		query := UserQuery{IdGt: P(5), ScoreLt: P(60)}
		actual, args := em.buildSelect(&query)
		expect := "SELECT id, score, memo FROM t_user WHERE id > ? AND score < ?"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !reflect.DeepEqual(args, []any{5, 60}) {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Build Select Without Where", func(t *testing.T) {
		query := UserQuery{}
		actual, args := em.buildSelect(&query)
		expect := "SELECT id, score, memo FROM t_user"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if len(args) != 0 {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Build Select with Page Clause", func(t *testing.T) {
		query := UserQuery{PageQuery: PageQuery{Page: 1, Size: 10}}
		actual, args := em.buildSelect(&query)
		expect := "SELECT id, score, memo FROM t_user LIMIT 10 OFFSET 0"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if len(args) != 0 {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Build Select with Sort Clause", func(t *testing.T) {
		query := UserQuery{PageQuery: PageQuery{Size: 5, Sort: "id"}}
		actual, args := em.buildSelect(&query)
		expect := "SELECT id, score, memo FROM t_user ORDER BY id LIMIT 5 OFFSET 0"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
			return
		}
		if len(args) != 0 {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Build Count", func(t *testing.T) {
		query := UserQuery{ScoreLt: P(60)}
		actual, args := em.buildCount(&query)
		expect := "SELECT count(0) FROM t_user WHERE score < ?"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !reflect.DeepEqual(args, []any{60}) {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Build Create Stmt", func(t *testing.T) {
		entity := UserEntity{Score: P(90), Memo: P("Great")}
		actual, args := em.buildCreate(entity)
		expect := "INSERT INTO t_user (score, memo) VALUES (?, ?)"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !reflect.DeepEqual(args, []any{90, "Great"}) {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Build Update Stmt", func(t *testing.T) {
		entity := UserEntity{Int64Id: NewInt64Id(2), Score: P(90), Memo: P("Great")}
		actual, args := em.buildUpdate(entity)
		expect := "UPDATE t_user SET score = ?, memo = ? WHERE id = ?"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !reflect.DeepEqual(args, []any{90, "Great", int64(2)}) {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Build Patch Stmt", func(t *testing.T) {
		entity := UserEntity{Int64Id: NewInt64Id(2), Memo: P("Great")}
		actual, args := em.buildPatchById(entity)
		expect := "UPDATE t_user SET memo = ? WHERE id = ?"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !reflect.DeepEqual(args, []any{"Great", int64(2)}) {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Support tag subquery", func(t *testing.T) {
		query := UserQuery{ScoreLtAvg: &UserQuery{MemoLike: P("Well")}}
		actual, args := em.buildSelect(&query)
		expect := "SELECT id, score, memo FROM t_user WHERE score < (SELECT avg(score) FROM t_user WHERE memo LIKE ?)"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
			return
		}
		if !reflect.DeepEqual(args, []any{"Well"}) {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Support tag subquery with Any", func(t *testing.T) {
		query := UserQuery{ScoreLtAny: &UserQuery{MemoLike: P("Well")}}
		actual, args := em.buildSelect(&query)
		expect := "SELECT id, score, memo FROM t_user WHERE score < ANY(SELECT score FROM t_user WHERE memo LIKE ?)"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
			return
		}
		if !reflect.DeepEqual(args, []any{"Well"}) {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Support tag subquery with All", func(t *testing.T) {
		query := UserQuery{ScoreLtAll: &UserQuery{MemoLike: P("Well")}}
		actual, args := em.buildSelect(&query)
		expect := "SELECT id, score, memo FROM t_user WHERE score < ALL(SELECT score FROM t_user WHERE memo LIKE ?)"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
			return
		}
		if !reflect.DeepEqual(args, []any{"Well"}) {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Support tag select and from", func(t *testing.T) {
		query := UserQuery{ScoreGtAvg: &UserQuery{MemoLike: P("Well")}}
		actual, args := em.buildSelect(&query)
		expect := "SELECT id, score, memo FROM t_user WHERE score > (SELECT avg(score) FROM t_user WHERE memo LIKE ?)"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
			return
		}
		if !reflect.DeepEqual(args, []any{"Well"}) {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Support subquery by fieldname: ScoreInScoreOfUser", func(t *testing.T) {
		query := UserQuery{ScoreInScoreOfUser: &UserQuery{Deleted: P(true)}}
		actual, args := em.buildSelect(&query)
		expect := "SELECT id, score, memo FROM t_user WHERE score IN (SELECT score FROM t_user WHERE deleted = ?)"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
			return
		}
		if !reflect.DeepEqual(args, []any{true}) {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Support subquery by fieldname: ScoreGtAvgScoreOfUser", func(t *testing.T) {
		RegisterEntity("t_user", "t_user")

		query := UserQuery{ScoreGtAvgScoreOfUser: &UserQuery{Deleted: P(true)}}
		actual, args := em.buildSelect(&query)
		expect := "SELECT id, score, memo FROM t_user WHERE score > (SELECT AVG(score) FROM t_user WHERE deleted = ?)"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
			return
		}
		if !reflect.DeepEqual(args, []any{true}) {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Error: Build UPDATE without SET columns", func(t *testing.T) {
		entity := UserEntity{Score: nil}
		_, _, err := em.buildPatchByQuery(entity, UserQuery{})
		expect := "at least one field should be updated"
		if err != nil && err.Error() != expect {
			t.Errorf("\nExpected: %s\nBut got <nil>", expect)
		}
	})

	t.Run("Error: Build UPDATE without WHERE clause", func(t *testing.T) {
		entity := UserEntity{Score: P(90)}
		_, _, err := em.buildPatchByQuery(entity, UserQuery{})
		expect := "deletion of all records is restricted"
		if err != nil && err.Error() != expect {
			t.Errorf("\nExpected: %s\nBut got <nil>", expect)
		}
	})
}
