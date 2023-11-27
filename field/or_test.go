package field

import (
	"testing"
)

type AccountOr struct {
	Username string
	Email    string
	Mobile   string
}

func TestOr(t *testing.T) {

	t.Run("Or Clause", func(t *testing.T) {
		actual, _ := ProcessOr(AccountOr{Username: "f0rb", Email: "f0rb"})
		expect := "username = ? OR email = ?"
		if actual != expect {
			t.Errorf("Expected: %s, but got %s", expect, actual)
		}
	})

}
