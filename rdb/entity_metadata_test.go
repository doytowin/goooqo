package rdb

import (
	. "github.com/doytowin/goooqo/core"
	. "github.com/doytowin/goooqo/test"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestBuildStmt(t *testing.T) {
	log.SetLevel(log.DebugLevel)

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
		query := UserQuery{IdGt: PInt(5), MemoNull: PBool(true)}
		actual, args := BuildWhereClause(query)
		expect := " WHERE id > ? AND memo IS NULL"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !(len(args) == 1 && args[0] == 5) {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Build Where Clause", func(t *testing.T) {
		query := UserQuery{IdGt: PInt(5), MemoNull: PBool(false)}
		actual, args := BuildWhereClause(query)
		expect := " WHERE id > ? AND memo IS NOT NULL"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !(len(args) == 1 && args[0] == 5) {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Build Select Statement", func(t *testing.T) {
		em := buildEntityMetadata[UserEntity]()
		query := UserQuery{IdGt: PInt(5), ScoreLt: PInt(60)}
		actual, args := em.buildSelect(&query)
		expect := "SELECT id, score, memo FROM User WHERE id > ? AND score < ?"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !(len(args) == 2 && args[0] == 5) || args[1] != 60 {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Build Select Without Where", func(t *testing.T) {
		em := buildEntityMetadata[UserEntity]()
		query := UserQuery{}
		actual, args := em.buildSelect(&query)
		expect := "SELECT id, score, memo FROM User"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if len(args) != 0 {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Build Select with Page Clause", func(t *testing.T) {
		em := buildEntityMetadata[UserEntity]()
		query := UserQuery{PageQuery: PageQuery{PageNumber: PInt(1), PageSize: PInt(10)}}
		actual, args := em.buildSelect(&query)
		expect := "SELECT id, score, memo FROM User LIMIT 10 OFFSET 0"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if len(args) != 0 {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Build Select with Sort Clause", func(t *testing.T) {
		em := buildEntityMetadata[UserEntity]()
		query := UserQuery{PageQuery: PageQuery{PageSize: PInt(5), Sort: PStr("id")}}
		actual, args := em.buildSelect(&query)
		expect := "SELECT id, score, memo FROM User ORDER BY id LIMIT 5 OFFSET 0"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
			return
		}
		if len(args) != 0 {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Build Count", func(t *testing.T) {
		em := buildEntityMetadata[UserEntity]()
		query := UserQuery{ScoreLt: PInt(60)}
		actual, args := em.buildCount(&query)
		expect := "SELECT count(0) FROM User WHERE score < ?"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !(len(args) == 1 && args[0] == 60) {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Build Create Stmt", func(t *testing.T) {
		em := buildEntityMetadata[UserEntity]()
		entity := UserEntity{Score: PInt(90), Memo: PStr("Great")}
		actual, args := em.buildCreate(entity)
		expect := "INSERT INTO User (score, memo) VALUES (?, ?)"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !(len(args) == 2 && args[0] == 90 && args[1] == "Great") {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Build Update Stmt", func(t *testing.T) {
		em := buildEntityMetadata[UserEntity]()
		entity := UserEntity{Int64Id: NewIntId(2), Score: PInt(90), Memo: PStr("Great")}
		actual, args := em.buildUpdate(entity)
		expect := "UPDATE User SET score = ?, memo = ? WHERE id = ?"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !(len(args) == 3 && args[0] == 90 && args[1] == "Great" && args[2] == int64(2)) {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Build Patch Stmt", func(t *testing.T) {
		em := buildEntityMetadata[UserEntity]()
		entity := UserEntity{Int64Id: NewIntId(2), Memo: PStr("Great")}
		actual, args := em.buildPatchById(entity)
		expect := "UPDATE User SET memo = ? WHERE id = ?"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !(len(args) == 2 && args[0] == "Great" && args[1] == int64(2)) {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Support tag subquery", func(t *testing.T) {
		em := buildEntityMetadata[UserEntity]()
		query := UserQuery{ScoreLt1: &UserQuery{MemoLike: PStr("Well")}}
		actual, args := em.buildSelect(&query)
		expect := "SELECT id, score, memo FROM User WHERE score < (SELECT avg(score) FROM User WHERE memo LIKE ?)"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
			return
		}
		if !(len(args) == 1 && args[0] == "Well") {
			t.Errorf("Args are not expected: %s", args)
		}
	})
}
