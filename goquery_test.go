package goquery

import (
	. "github.com/doytowin/goquery/field"
	. "github.com/doytowin/goquery/util"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestBuild(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	t.Run("Build Where Clause", func(t *testing.T) {
		query := UserQuery{IdGt: PInt(5), MemoNull: true}
		actual, args := BuildWhereClause(query)
		expect := " WHERE id > ? AND memo IS NULL"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !(len(args) == 2 && args[0] == int64(5) && args[1] == true) {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Build Select Statement", func(t *testing.T) {
		em := buildEntityMetadata[UserEntity](UserEntity{})
		query := UserQuery{IdGt: PInt(5), ScoreLt: PInt(60)}
		actual, args := em.buildSelect(query)
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
		actual, args := em.buildSelect(query)
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
		actual, args := em.buildSelect(query)
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
		actual, args := em.buildCount(query)
		expect := "SELECT count(0) FROM User WHERE score < ?"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !(len(args) == 1 && args[0] != 60) {
			t.Errorf("Args are not expected: %s", args)
		}
	})

}
