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

	t.Run("Support Basic Query", func(t *testing.T) {
		mongoDataAccess := BuildMongoDataAccess[context.Context, InventoryEntity](client, func() InventoryEntity { return InventoryEntity{} })
		actual, err := mongoDataAccess.Query(ctx, InventoryQuery{})
		expect := "journal"
		if !(err == nil && len(actual) == 5) {
			t.Errorf("%s\nExpected: %d\n     Got: %d", err, 4, len(actual))
		} else if !(actual[0].Item == expect) {
			t.Errorf("%s\nExpected: %s\n     Got: %s", err, expect, actual[0].Item)
		} else if !(actual[0].Size.H == 14.) {
			t.Errorf("%s\nExpected: %f\n     Got: %f", err, 14., actual[0].Size.H)
		}
		log.Println(actual)
	})
}
