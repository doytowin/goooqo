package goquery

import (
	. "github.com/doytowin/goquery/util"
	"testing"
)

func TestBuildStmt(t *testing.T) {
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
		actual, args := em.buildPatch(entity)
		expect := "UPDATE User SET memo = ? WHERE id = ?"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !(len(args) == 2 && args[0] == "Great" && args[1] == int64(2)) {
			t.Errorf("Args are not expected: %s", args)
		}
	})

}
