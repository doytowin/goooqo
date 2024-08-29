package rdb

import (
	. "github.com/doytowin/goooqo/core"
	"reflect"
	"testing"
)

func TestOr(t *testing.T) {

	t.Run("Or Clause", func(t *testing.T) {
		actual, _ := fpForOr.Process(reflect.ValueOf(&TestCond{Username: PStr("f0rb"), Email: PStr("f0rb")}))
		expect := "(username = ? OR email = ?)"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})

	t.Run("Or Interface", func(t *testing.T) {
		query := TestQuery{TestOr: &TestCond{Username: PStr("f0rb"), Email: PStr("f0rb")}, Deleted: PBool(true)}
		actual, args := BuildWhereClause(query)
		expect := " WHERE (username = ? OR email = ?) AND deleted = ?"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !(len(args) == 3 && args[0] == "f0rb" && args[1] == "f0rb" && args[2] == true) {
			t.Errorf("Unexpected args: %v", args)
		}
	})

	t.Run("Build Or Clause with And", func(t *testing.T) {
		accountAnd := TestCond{Email: PStr("f0rb@qq.com"), Mobile: PStr("01008888")}
		query := TestQuery{TestOr: &TestCond{Username: PStr("f0rb"), TestAnd: &accountAnd}, Deleted: PBool(true)}
		actual, args := BuildWhereClause(query)
		expect := " WHERE (username = ? OR email = ? AND mobile = ?) AND deleted = ?"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !(len(args) == 4 && args[0] == "f0rb" &&
			args[1] == "f0rb@qq.com" && args[2] == "01008888" && args[3] == true) {
			t.Errorf("Unexpected args: %v", args)
		}
	})

	t.Run("Build Or Clause for array", func(t *testing.T) {
		query := TestQuery{EmailEndOr: &[]string{"icloud.com", "gmail.com"}, Deleted: PBool(true)}
		actual, args := BuildWhereClause(query)
		expect := " WHERE (email LIKE ? OR email LIKE ?) AND deleted = ?"
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
		if !(len(args) == 3 && args[0] == "%icloud.com" &&
			args[1] == "%gmail.com" && args[2] == true) {
			t.Errorf("Unexpected args: %v", args)
		}
	})

}
