package mongodb

import (
	"github.com/doytowin/goooqo"
)

type InventoryQuery struct {
	goooqo.PageQuery
	QtyGt *int
}

type InventoryEntity struct {
	MongoId `bson:",inline"`
	Item    string
	Size    struct {
		H   float64
		W   float64
		Uom string
	}
	Qty    int
	Status string
}

func (r InventoryEntity) Database() string {
	return "doytowin"
}

func (r InventoryEntity) Collection() string {
	return "inventory"
}
