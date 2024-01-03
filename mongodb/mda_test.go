package mongodb

import (
	"context"
	. "github.com/doytowin/goooqo/core"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestMongoDataAccess(t *testing.T) {

	ctx := context.Background()
	var client = Connect(ctx, "app.properties")
	defer Disconnect(client, ctx)

	tm := NewMongoTransactionManager(client)
	mongoDataAccess := NewMongoDataAccess[InventoryEntity](tm, func() InventoryEntity { return InventoryEntity{} })

	t.Run("Support Basic Query", func(t *testing.T) {
		tc, _ := mongoDataAccess.StartTransaction(ctx)
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
		log.Debugln(actual)
	})

	t.Run("Support Custom Query Builder", func(t *testing.T) {
		tc, _ := mongoDataAccess.StartTransaction(ctx)
		defer tc.Rollback()
		actual, err := mongoDataAccess.Query(tc, InventoryQuery{QtyGt: PInt(70)})
		expect := 2
		if !(err == nil && len(actual) == expect) {
			t.Errorf("%s\nExpected: %d\n     Got: %d", err, expect, len(actual))
		} else if !(actual[0].Qty == 100) {
			t.Errorf("\nExpected: %f\n     Got: %d", 100., actual[0].Qty)
		}
		log.Debugln(actual)
	})

	t.Run("Support Page Query", func(t *testing.T) {
		tc, _ := mongoDataAccess.StartTransaction(ctx)
		defer tc.Rollback()
		actual, err := mongoDataAccess.Page(tc, InventoryQuery{QtyGt: PInt(70)})
		expect := 2
		if !(err == nil && len(actual.List) == expect) {
			t.Errorf("%s\nExpected: %d\n     Got: %d", err, expect, len(actual.List))
		} else if !(actual.Total == int64(expect)) {
			t.Errorf("\nExpected: %d\n     Got: %d", expect, actual.Total)
		} else if !(actual.List[0].Qty == 100) {
			t.Errorf("\nExpected: %f\n     Got: %d", 100., actual.List[0].Qty)
		}
		log.Debugln(actual)
	})
}
