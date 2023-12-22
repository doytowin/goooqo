package rdb

import (
	. "github.com/doytowin/go-query/core"
	"testing"
)

func TestBuildWhereClause(t *testing.T) {

	t.Run("support custom condition", func(t *testing.T) {
		query := TestQuery{Account: PStr("f0rb"), Deleted: PBool(true)}
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
