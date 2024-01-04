package main

import (
	"github.com/doytowin/goooqo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InventoryEntity struct {
	Id   primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Item string             `json:"item,omitempty"`
	Size struct {
		H   float64 `json:"h,omitempty"`
		W   float64 `json:"w,omitempty"`
		Uom string  `json:"uom,omitempty"`
	} `json:"size"`
	Qty    int    `json:"qty,omitempty"`
	Status string `json:"status,omitempty"`
}

func (r InventoryEntity) GetId() any {
	return r.Id
}

func (r InventoryEntity) SetId(self any, id any) error {
	panic("implement me")
}

func (r InventoryEntity) Database() string {
	return "doytowin"
}

func (r InventoryEntity) Collection() string {
	return "inventory"
}

type InventoryQuery struct {
	goooqo.PageQuery
}
