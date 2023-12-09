package goquery

import (
	"github.com/doytowin/goquery/field"
	. "github.com/doytowin/goquery/util"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestBuildStmt(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	t.Run("Build Where Clause", func(t *testing.T) {
		query := UserQuery{IdGt: PInt(5), MemoNull: true}
		actual, args := field.BuildWhereClause(query)
		expect := " WHERE id > ? AND memo IS NULL"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !(len(args) == 1 && args[0] == int64(5)) {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Build Select Statement", func(t *testing.T) {
		em := buildEntityMetadata[UserEntity](UserEntity{})
		query := UserQuery{IdGt: PInt(5), ScoreLt: PInt(60)}
		actual, args := em.buildSelect(&query)
		expect := "SELECT id, score, memo FROM User WHERE id > ? AND score < ?"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !(len(args) == 2 && args[0] == int64(5)) || args[1] != int64(60) {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Build Select Without Where", func(t *testing.T) {
		em := buildEntityMetadata[UserEntity](UserEntity{})
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
		em := buildEntityMetadata[UserEntity](UserEntity{})
		query := UserQuery{PageQuery: PageQuery{PInt(1), PInt(10)}}
		actual, args := em.buildSelect(&query)
		expect := "SELECT id, score, memo FROM User LIMIT 10 OFFSET 0"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if len(args) != 0 {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Build Count", func(t *testing.T) {
		em := buildEntityMetadata[UserEntity](UserEntity{})
		query := UserQuery{ScoreLt: PInt(60)}
		actual, args := em.buildCount(&query)
		expect := "SELECT count(0) FROM User WHERE score < ?"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !(len(args) == 1 && args[0] != 60) {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Build Create Stmt", func(t *testing.T) {
		em := buildEntityMetadata[UserEntity](UserEntity{})
		entity := UserEntity{Score: PInt(90), Memo: PStr("Great")}
		actual, args := em.buildCreate(entity)
		expect := "INSERT INTO User (score, memo) VALUES (?, ?)"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !(len(args) == 2 && args[0] == int64(90) && args[1] == "Great") {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Build Update Stmt", func(t *testing.T) {
		em := buildEntityMetadata[UserEntity](UserEntity{})
		entity := UserEntity{2, PInt(90), PStr("Great")}
		actual, args := em.buildUpdate(entity)
		expect := "UPDATE User SET score = ?, memo = ? WHERE id = ?"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !(len(args) == 3 && args[0] == int64(90) && args[1] == "Great" && args[2] == int64(2)) {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Build Patch Stmt", func(t *testing.T) {
		em := buildEntityMetadata[UserEntity](UserEntity{})
		entity := UserEntity{Id: 2, Memo: PStr("Great")}
		actual, args := em.buildPatchById(entity)
		expect := "UPDATE User SET memo = ? WHERE id = ?"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !(len(args) == 2 && args[0] == "Great" && args[1] == int64(2)) {
			t.Errorf("Args are not expected: %s", args)
		}
	})

}
