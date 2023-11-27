package query

import (
	"testing"
)

type UserQuery struct {
	idGt     *int
	scoreLt  *int
	memoNull bool
}

func intPtr(o int) *int {
	return &o
}

func TestBuild(t *testing.T) {
	t.Run("Build Where Clause", func(t *testing.T) {
		query := UserQuery{idGt: intPtr(5), memoNull: true}
		actual, args := BuildConditions(query)
		expect := "id > ? AND memo IS NULL"
		if actual != expect {
			t.Errorf("Expected: %s, but got %s", expect, actual)
		}
		if len(args) != 1 || args[0] != int64(5) {
			t.Errorf("Args are not expected: %s", args)
		}
	})

}
