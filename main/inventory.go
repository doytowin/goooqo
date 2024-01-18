package main

import (
	"github.com/doytowin/goooqo"
	"github.com/doytowin/goooqo/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InventoryQuery struct {
	goooqo.PageQuery
	Id             *primitive.ObjectID
	IdNot          *primitive.ObjectID
	IdNe           *primitive.ObjectID
	IdIn           *[]primitive.ObjectID
	IdNotIn        *[]primitive.ObjectID
	Qty            *int
	QtyGt          *int
	QtyLt          *int
	QtyGe          *int
	QtyLe          *int
	Size           *SizeQuery
	StatusNull     *bool
	ItemContain    *string
	ItemNotContain *string
	ItemStart      *string
	ItemNotStart   *string
	ItemEnd        *string
	ItemNotEnd     *string
	CustomFilter   *primitive.M
	*QtyOr
	Search *string
}

type QtyOr struct {
	QtyLt  *int
	QtyGe  *int
	Size   *SizeQuery
	SizeOr *SizeQuery
}

type Unit struct {
	Name     *string
	NameNull *bool `column:"size.unit.name"`
}

type SizeQuery struct {
	HLt  *float64
	HGe  *float64
	Unit *Unit
}

type SizeDoc struct {
	H   float64 `json:"h,omitempty" bson:"h"`
	W   float64 `json:"w,omitempty" bson:"w"`
	Uom string  `json:"uom,omitempty" bson:"uom"`
}

type InventoryEntity struct {
	mongodb.MongoId `bson:",inline"`
	Item            string  `json:"item,omitempty" bson:"item" column:"item,index"`
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
