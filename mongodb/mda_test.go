package mongodb

import (
	"context"
	. "github.com/doytowin/goooqo/core"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestMongoDataAccess(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	ctx := context.Background()
	var client = Connect(ctx, "app.properties")
	defer Disconnect(client, ctx)

	tm := NewMongoTransactionManager(client)
	inventoryDataAccess := NewMongoDataAccess[InventoryEntity](tm, func() InventoryEntity { return InventoryEntity{} })

	t.Run("Support Basic Query", func(t *testing.T) {
		tc, _ := inventoryDataAccess.StartTransaction(ctx)
		defer tc.Rollback()
		actual, err := inventoryDataAccess.Query(tc, InventoryQuery{})
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
		tc, _ := inventoryDataAccess.StartTransaction(ctx)
		defer tc.Rollback()
		actual, err := inventoryDataAccess.Query(tc, InventoryQuery{QtyGt: PInt(70)})
		expect := 2
		if !(err == nil && len(actual) == expect) {
			t.Errorf("%s\nExpected: %d\n     Got: %d", err, expect, len(actual))
		} else if !(actual[0].Qty == 100) {
			t.Errorf("\nExpected: %f\n     Got: %d", 100., actual[0].Qty)
		}
		log.Debugln(actual)
	})

	t.Run("Support Page Query", func(t *testing.T) {
		tc, _ := inventoryDataAccess.StartTransaction(ctx)
		defer tc.Rollback()
		actual, err := inventoryDataAccess.Page(tc, InventoryQuery{QtyGt: PInt(70)})
		expect := 2
		if !(err == nil && len(actual.List) == expect) {
			t.Errorf("%s\nExpected: %d\n     Got: %d", err, expect, len(actual.List))
		} else if !(actual.Total == int64(expect)) {
			t.Errorf("\nExpected: %d\n     Got: %d", expect, actual.Total)
		} else if !(actual.List[0].Qty == 100) {
			t.Errorf("\nExpected: %d\n     Got: %d", 100, actual.List[0].Qty)
		}
		log.Debugln(actual)
	})

	t.Run("Support Page Query with Pagination", func(t *testing.T) {
		tc, _ := inventoryDataAccess.StartTransaction(ctx)
		defer tc.Rollback()
		inventoryQuery := InventoryQuery{QtyGt: PInt(70)}
		inventoryQuery.PageSize = PInt(1)
		inventoryQuery.PageNumber = PInt(2)
		actual, err := inventoryDataAccess.Page(tc, inventoryQuery)
		if !(err == nil && len(actual.List) == 1) {
			t.Errorf("%s\nExpected: %d\n     Got: %d", err, 1, len(actual.List))
		} else if !(actual.Total == int64(2)) {
			t.Errorf("\nExpected: %d\n     Got: %d", 2, actual.Total)
		} else if !(actual.List[0].Qty == 75) {
			t.Errorf("\nExpected: %d\n     Got: %d", 75, actual.List[0].Qty)
		}
		log.Debugln(actual)
	})

	t.Run("Support Delete by id", func(t *testing.T) {
		tc, _ := inventoryDataAccess.StartTransaction(ctx)
		defer tc.Rollback()
		id, _ := primitive.ObjectIDFromHex("657bbb49675e5c32a2b8af73")
		actual, err := inventoryDataAccess.Delete(tc, id)
		expect := int64(1)
		if !(err == nil && actual == expect) {
			t.Errorf("%s\nExpected: %d\n     Got: %d", err, expect, actual)
		}
		inventory, err := inventoryDataAccess.Get(tc, "657bbb49675e5c32a2b8af73")
		if !(err != nil && inventory == nil) {
			t.Errorf("%s\nExpected: %v\n     Got: %v", err, nil, inventory)
		}
		log.Debugln(actual)
	})

	t.Run("Support Delete by Query", func(t *testing.T) {
		tc, _ := inventoryDataAccess.StartTransaction(ctx)
		defer tc.Rollback()
		actual, err := inventoryDataAccess.DeleteByQuery(tc, InventoryQuery{QtyGt: PInt(70)})
		expect := int64(2)
		if !(err == nil && actual == expect) {
			t.Errorf("%s\nExpected: %d\n     Got: %d", err, expect, actual)
		}
		cnt, err := inventoryDataAccess.Count(tc, InventoryQuery{})
		if !(err == nil && cnt == int64(3)) {
			t.Errorf("%s\nExpected: %d\n     Got: %d", err, 3, cnt)
		}
		log.Debugln(actual)
	})

	t.Run("Support Create", func(t *testing.T) {
		tc, _ := inventoryDataAccess.StartTransaction(ctx)
		defer tc.Rollback()
		entity := InventoryEntity{
			Item:   "eraser",
			Size:   SizeDoc{3.5, 2, "cm"},
			Qty:    20,
			Status: "A",
		}
		actual, err := inventoryDataAccess.Create(tc, &entity)
		if !(err == nil && !entity.Id.IsZero()) {
			t.Errorf("%s\n     Got: %s", err, entity.Id.Hex())
		}
		cnt, err := inventoryDataAccess.Count(tc, InventoryQuery{})
		if !(err == nil && cnt == int64(6)) {
			t.Errorf("%s\nExpected: %d\n     Got: %d", err, 6, cnt)
		}
		log.Debugln(actual)
	})

	t.Run("Support Update", func(t *testing.T) {
		tc, _ := inventoryDataAccess.StartTransaction(ctx)
		defer tc.Rollback()

		newQty := 123
		Id, _ := primitive.ObjectIDFromHex("657bbb49675e5c32a2b8af72")
		inventory, _ := inventoryDataAccess.Get(tc, Id)

		inventory.Qty = newQty
		actual, err := inventoryDataAccess.Update(tc, *inventory)

		if !(err == nil && actual == 1) {
			t.Errorf("%s\nExpected: %d\n     Got: %d", err, 1, actual)
		}

		newE, err := inventoryDataAccess.Get(tc, Id)
		if !(err == nil) {
			t.Error(err)
		} else if !(newE.Qty == newQty) {
			t.Errorf("\nExpected: %d\n     Got: %d", newQty, newE.Qty)
		}
		log.Println(actual)
	})

	t.Run("Support Create Multiple Entities", func(t *testing.T) {
		tc, _ := inventoryDataAccess.StartTransaction(ctx)
		defer tc.Rollback()
		entities := []InventoryEntity{
			{
				Item:   "eraser",
				Size:   SizeDoc{3.5, 2, "cm"},
				Qty:    20,
				Status: "A",
			},
			{
				Item:   "keyboard",
				Size:   SizeDoc{40, 15.5, "cm"},
				Qty:    10,
				Status: "D",
			},
		}
		actual, err := inventoryDataAccess.CreateMulti(tc, entities)
		if !(err == nil) {
			t.Fatal(err)
		} else if entities[0].Id == nil {
			t.Fatal("id should not be nil")
		}

		cnt, err := inventoryDataAccess.Count(tc, InventoryQuery{})
		if !(err == nil && cnt == int64(7)) {
			t.Errorf("%s\nExpected: %d\n     Got: %d", err, 7, cnt)
		}
		log.Debugln(actual)
	})
}
