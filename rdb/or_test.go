package rdb

import (
	. "github.com/doytowin/goooqo/core"
	"testing"
)

func TestOr(t *testing.T) {

	t.Run("Or Clause", func(t *testing.T) {
		actual, _ := ProcessOr(AccountOr{Username: PStr("f0rb"), Email: PStr("f0rb")})
		expect := "(username = ? OR email = ?)"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})

	t.Run("Or Interface", func(t *testing.T) {
		query := TestQuery{AccountOr: &AccountOr{Username: PStr("f0rb"), Email: PStr("f0rb")}, Deleted: PBool(true)}
		actual, args := BuildWhereClause(query)
		expect := " WHERE (username = ? OR email = ?) AND deleted = ?"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !(len(args) == 3 && args[0] == "f0rb" && args[1] == "f0rb" && args[2] == true) {
			t.Errorf("Unexpected args: %v", args)
		}
	})

}
