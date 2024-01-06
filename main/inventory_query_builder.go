package main

import . "go.mongodb.org/mongo-driver/bson/primitive"

func (q InventoryQuery) BuildFilter() []D {
	d := make([]D, 0, 10)
	if q.Qty != nil {
		d = append(d, D{{"qty", q.Qty}})
	}
	if q.QtyGt != nil {
		d = append(d, D{{"qty", D{{"$gt", q.QtyGt}}}})
	}
	if q.QtyLt != nil {
		d = append(d, D{{"qty", D{{"$lt", q.QtyLt}}}})
	}
	if q.QtyGe != nil {
		d = append(d, D{{"qty", D{{"$gte", q.QtyGe}}}})
	}
	if q.QtyLe != nil {
		d = append(d, D{{"qty", D{{"$lte", q.QtyLe}}}})
	}
	return d
}
