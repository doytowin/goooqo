package goquery

import (
	"testing"
)

func TestBuildStmt(t *testing.T) {
	t.Run("Build Create Stmt", func(t *testing.T) {
		em := buildEntityMetadata[UserEntity](UserEntity{})
		entity := UserEntity{Score: 90, Memo: "Great"}
		actual, args := em.buildCreate(entity)
		expect := "INSERT INTO User (score, memo) VALUES (?, ?)"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !(len(args) == 2 && args[0] == int64(90) && args[1] == "Great") {
			t.Errorf("Args are not expected: %s", args)
		}
	})
}
