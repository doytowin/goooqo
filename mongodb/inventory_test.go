package mongodb

import (
	"github.com/doytowin/goooqo"
	"github.com/doytowin/goooqo/core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InventoryQuery struct {
	goooqo.PageQuery
	QtyGt *int
}

type InventoryEntity struct {
	Id   primitive.ObjectID `bson:"_id,omitempty"`
	Item string
	Size struct {
		H   float64
		W   float64
		Uom string
	}
	Qty    int
	Status string
}

func (r InventoryEntity) GetId() any {
	return r.Id
}

func (r InventoryEntity) SetId(self any, id any) {
	objectID, err := resolveId(id)
	if core.NoError(err) {
		self.(*InventoryEntity).Id = objectID
	}
}

func (r InventoryEntity) Database() string {
	return "doytowin"
}

func (r InventoryEntity) Collection() string {
	return "inventory"
}
