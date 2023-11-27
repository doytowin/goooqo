package goquery

import (
	fp "github.com/doytowin/doyto-query-go-sql/field"
	"testing"
)

type UserEntity struct {
	Id    int
	Score int
	Memo  string
}

type UserQuery struct {
	IdGt     *int
	ScoreLt  *int
	MemoNull bool
}

func intPtr(o int) *int {
	return &o
}

func TestBuild(t *testing.T) {
	t.Run("Build Where Clause", func(t *testing.T) {
		query := UserQuery{IdGt: intPtr(5), MemoNull: true}
		actual, args := fp.BuildWhereClause(query)
		expect := " WHERE id > ? AND memo IS NULL"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if len(args) != 1 || args[0] != int64(5) {
			t.Errorf("Args are not expected: %s", args)
		}
	})

	t.Run("Build Select Statement", func(t *testing.T) {
		em := BuildEntityMetadata[UserEntity](UserEntity{})
		query := UserQuery{IdGt: intPtr(5), ScoreLt: intPtr(60)}
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
		em := BuildEntityMetadata[UserEntity](UserEntity{})
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

}
