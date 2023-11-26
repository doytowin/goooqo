package query

import (
	"testing"
)

type UserQuery struct {
	idGt     *int
	scoreLt  *int
	memoNull *bool
}

func intPtr(o int) *int {
	return &o
}

func boolPtr(o bool) *bool {
	return &o
}

func TestBuild(t *testing.T) {
	t.Run("Build Where Clause", func(t *testing.T) {
		query := UserQuery{idGt: intPtr(5), memoNull: boolPtr(true)}
		actual := BuildConditions(query)
		expect := "id > ? AND memo IS NULL"
		if actual != expect {
			t.Errorf("Expected: %s, but got %s", expect, actual)
		}
	})

}
