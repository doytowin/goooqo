package mongodb

import (
	. "go.mongodb.org/mongo-driver/bson/primitive"
)

func (q InventoryQuery) BuildFilter() A {
	d := make(A, 0, 4)
	if q.QtyGt != nil {
		d = append(d, D{{"qty", D{{"$gt", q.QtyGt}}}})
	}
	return d
}
