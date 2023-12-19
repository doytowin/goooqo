package mongodb

import (
	"context"
	"log"
	"testing"
)

func TestMongoDataAccess(t *testing.T) {

	ctx := context.Background()
	var client = Connect(ctx, "app.properties")
	defer Disconnect(client, ctx)

	tm := NewMongoTransactionManager(client)
	mongoDataAccess := NewMongoDataAccess[InventoryEntity](tm, func() InventoryEntity { return InventoryEntity{} })

	t.Run("Support Basic Query", func(t *testing.T) {
		tc := mongoDataAccess.StartTransaction(ctx)
		defer tc.Rollback()
		actual, err := mongoDataAccess.Query(tc, InventoryQuery{})
		expect := "journal"
		if !(err == nil && len(actual) == 5) {
			t.Errorf("%s\nExpected: %d\n     Got: %d", err, 5, len(actual))
		} else if !(actual[0].Item == expect) {
			t.Errorf("%s\nExpected: %s\n     Got: %s", err, expect, actual[0].Item)
		} else if !(actual[0].Size.H == 14.) {
			t.Errorf("%s\nExpected: %f\n     Got: %f", err, 14., actual[0].Size.H)
		}
		log.Println(actual)
	})
}
