package query

import (
	"testing"
)

type UserQuery struct {
	idGt     int
	memoNull bool
}

func TestBuild(t *testing.T) {

	t.Run("Build Where Clause", func(t *testing.T) {
		actual := BuildConditions(UserQuery{idGt: 5, memoNull: true})
		expect := "id > ? AND memo IS NULL"
		if actual != expect {
			t.Errorf("Expected: %s, but got %s", expect, actual)
		}
	})

}
