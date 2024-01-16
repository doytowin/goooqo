package mongodb

import (
	"context"
	"encoding/json"
	. "github.com/doytowin/goooqo/core"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func PF64(f float64) *float64 {
	return &f
}

func TestMongoDataAccess(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	ctx := context.Background()
	var client = Connect(ctx, "app.properties")
	defer Disconnect(client, ctx)

	tm := NewMongoTransactionManager(client)
	inventoryDataAccess := NewMongoDataAccess[InventoryEntity](tm)

	t.Run("Support Basic Query", func(t *testing.T) {
		tc, _ := inventoryDataAccess.StartTransaction(ctx)
		defer tc.Rollback()
		actual, err := inventoryDataAccess.Query(tc, InventoryQuery{})
		expect := "journal"
		if !(err == nil && len(actual) == 5) {
			t.Errorf("%s\nExpected: %d\n     Got: %d", err, 5, len(actual))
		} else if !(*actual[0].Item == expect) {
			t.Errorf("%s\nExpected: %s\n     Got: %s", err, expect, *actual[0].Item)
		} else if !(*actual[0].Size.H == 14.) {
			t.Errorf("%s\nExpected: %f\n     Got: %f", err, 14., *actual[0].Size.H)
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
		} else if !(*actual[0].Qty == 100) {
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
		} else if !(*actual.List[0].Qty == 100) {
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
		} else if !(*actual.List[0].Qty == 75) {
			t.Errorf("\nExpected: %d\n     Got: %d", 75, actual.List[0].Qty)
		}
		log.Debugln(actual)
	})

	t.Run("Support Page Query with Sort", func(t *testing.T) {
		tc, _ := inventoryDataAccess.StartTransaction(ctx)
		defer tc.Rollback()

		inventoryQuery := InventoryQuery{QtyGt: PInt(70)}
		inventoryQuery.Sort = PStr("qty,desc")

		actual, err := inventoryDataAccess.Page(tc, inventoryQuery)
		if !(err == nil) {
			t.Error(err)
		} else if !(actual.Total == int64(2)) {
			t.Errorf("\nExpected: %d\n     Got: %d", 2, actual.Total)
		} else if !(*actual.List[1].Qty == 75) {
			t.Errorf("\nExpected: %d\n     Got: %d", 75, actual.List[1].Qty)
		}
	})

	t.Run("Support OR Query", func(t *testing.T) {
		tc, _ := inventoryDataAccess.StartTransaction(ctx)
		defer tc.Rollback()

		inventoryQuery := InventoryQuery{QtyOr: &QtyOr{QtyLt: PInt(30), QtyGe: PInt(80)}}

		page, err := inventoryDataAccess.Page(tc, inventoryQuery)
		if !(err == nil) {
			t.Error(err)
		} else if !(page.Total == int64(2)) {
			t.Errorf("\nExpected: %d\n     Got: %d", 2, page.Total)
		} else if !(*page.List[0].Qty == 25) {
			t.Errorf("\nExpected: %d\n     Got: %d", 25, page.List[0].Qty)
		} else if !(*page.List[1].Qty == 100) {
			t.Errorf("\nExpected: %d\n     Got: %d", 100, page.List[1].Qty)
		}
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

	t.Run("Support Delete by Page", func(t *testing.T) {
		tc, _ := inventoryDataAccess.StartTransaction(ctx)
		defer tc.Rollback()
		inventoryQuery := InventoryQuery{QtyGt: PInt(30)}
		inventoryQuery.PageNumber = PInt(2)
		inventoryQuery.PageSize = PInt(2)
		actual, err := inventoryDataAccess.DeleteByQuery(tc, inventoryQuery)
		expect := int64(2)
		if !(err == nil && actual == expect) {
			t.Fatalf("%s\nExpected: %d\n     Got: %d", err, expect, actual)
		}

		entities, err := inventoryDataAccess.Query(tc, InventoryQuery{})
		if !(err == nil) {
			t.Error(err)
		}
		assertEquals(t, 3, len(entities))
		assertEquals(t, "657bbb49675e5c32a2b8af74", entities[2].Id.Hex())
	})

	t.Run("Support Create", func(t *testing.T) {
		tc, _ := inventoryDataAccess.StartTransaction(ctx)
		defer tc.Rollback()
		entity := InventoryEntity{
			Item:   PStr("eraser"),
			Size:   &SizeDoc{PF64(3.5), PF64(2), PStr("cm")},
			Qty:    PInt(20),
			Status: PStr("A"),
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

		inventory.Qty = &newQty
		actual, err := inventoryDataAccess.Update(tc, *inventory)

		if !(err == nil && actual == 1) {
			t.Errorf("%s\nExpected: %d\n     Got: %d", err, 1, actual)
		}

		newE, err := inventoryDataAccess.Get(tc, Id)
		if !(err == nil) {
			t.Error(err)
		} else if !(*newE.Qty == newQty) {
			t.Errorf("\nExpected: %d\n     Got: %d", newQty, newE.Qty)
		}
		log.Println(actual)
	})

	t.Run("Support Create Multiple Entities", func(t *testing.T) {
		tc, _ := inventoryDataAccess.StartTransaction(ctx)
		defer tc.Rollback()
		entities := []InventoryEntity{
			{
				Item:   PStr("eraser"),
				Size:   &SizeDoc{PF64(3.5), PF64(2), PStr("cm")},
				Qty:    PInt(20),
				Status: PStr("A"),
			},
			{
				Item:   PStr("keyboard"),
				Size:   &SizeDoc{PF64(40), PF64(15.5), PStr("cm")},
				Qty:    PInt(10),
				Status: PStr("D"),
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

	t.Run("Support Patch", func(t *testing.T) {
		tc, _ := inventoryDataAccess.StartTransaction(ctx)
		defer tc.Rollback()

		Id, _ := primitive.ObjectIDFromHex("657bbb49675e5c32a2b8af72")
		inventory := InventoryEntity{MongoId: NewMongoId(&Id)}
		inventory.Qty = PInt(123)
		inventory.Size = &SizeDoc{H: PF64(20.5)}

		cnt, err := inventoryDataAccess.Patch(tc, inventory)

		if !(err == nil && cnt == 1) {
			t.Errorf("%s\nExpected: %d\n     Got: %d", err, 1, cnt)
		}

		newE, err := inventoryDataAccess.Get(tc, Id)
		data, _ := json.Marshal(newE)
		actual := string(data)
		expect := `{"id":"657bbb49675e5c32a2b8af72","item":"journal","size":{"h":20.5,"w":21,"uom":"cm"},"qty":123,"status":"A"}`

		if !(err == nil) {
			t.Error(err)
		} else if !(actual == expect) {
			t.Errorf("\nExpected: %s\n     Got: %s", expect, actual)
		}
		log.Println(actual)
	})

	t.Run("Support Patch by Query", func(t *testing.T) {
		tc, _ := inventoryDataAccess.StartTransaction(ctx)
		defer tc.Rollback()

		newQty := 70
		cnt, err := inventoryDataAccess.PatchByQuery(tc, InventoryEntity{Qty: &newQty}, InventoryQuery{QtyGt: &newQty})

		if !(err == nil && cnt == 2) {
			t.Errorf("%s\nExpected: %d\n     Got: %d", err, 2, cnt)
		}

		cnt, err = inventoryDataAccess.Count(tc, InventoryQuery{QtyGt: &newQty})
		if !(err == nil && cnt == 0) {
			t.Errorf("%s\nExpected: %d\n     Got: %d", err, 0, cnt)
		}
	})

	t.Run("Support Patch by Page", func(t *testing.T) {
		tc, _ := inventoryDataAccess.StartTransaction(ctx)
		defer tc.Rollback()

		newQty := 30
		inventoryQuery := InventoryQuery{QtyGt: &newQty}
		inventoryQuery.PageNumber = PInt(2)
		inventoryQuery.PageSize = PInt(2)
		cnt, err := inventoryDataAccess.PatchByQuery(tc, InventoryEntity{Qty: &newQty}, inventoryQuery)

		assertNoError(t, err)
		assertEquals(t, int64(2), cnt)

		entities, _ := inventoryDataAccess.Query(tc, InventoryQuery{})
		assertEquals(t, 100, *entities[2].Qty)
		assertEquals(t, 30, *entities[4].Qty)
	})

}

func assertNoError(t *testing.T, err error) {
	if !(err == nil) {
		t.Error(err)
	}
}

func assertEquals(t *testing.T, expect any, actual any) {
	if !(actual == expect) {
		t.Error("\nExpected: ", expect, "\n\t Got: ", actual)
	}
}
