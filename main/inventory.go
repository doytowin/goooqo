package main

import (
	"github.com/doytowin/goooqo"
	"github.com/doytowin/goooqo/mongodb"
)

type InventoryQuery struct {
	goooqo.PageQuery
	Qty   *int
	QtyGt *int
	QtyLt *int
	QtyGe *int
}

type SizeDoc struct {
	H   float64 `json:"h,omitempty" bson:"h"`
	W   float64 `json:"w,omitempty" bson:"w"`
	Uom string  `json:"uom,omitempty" bson:"uom"`
}

type InventoryEntity struct {
	mongodb.MongoId `bson:",inline"`
	Item            string  `json:"item,omitempty" bson:"item"`
	Size            SizeDoc `json:"size" bson:"size"`
	Qty             int     `json:"qty,omitempty" bson:"qty"`
	Status          string  `json:"status,omitempty" bson:"status"`
}

func (r InventoryEntity) Database() string {
	return "doytowin"
}

func (r InventoryEntity) Collection() string {
	return "inventory"
}
