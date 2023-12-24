package rdb

import (
	"context"
	"testing"
)

func TestWeb(t *testing.T) {

	db := Connect("app.properties")
	defer Disconnect(db)

	ctx := context.Background()
	tm := NewTransactionManager(db)

	t.Run("Should not start tx repeated", func(t *testing.T) {
		tc, _ := tm.StartTransaction(ctx)
		tc2, _ := tm.StartTransaction(tc)
		defer tc2.Rollback()

		if tc2 != tc {
			t.Error("Should not start tx repeated")
		}
	})

	t.Run("Support get parent context from TransactionContext", func(t *testing.T) {
		tc, _ := tm.StartTransaction(ctx)
		defer tc.Rollback()

		if tc.Parent() != ctx {
			t.Error("Should return parent context")
		}
	})
}
