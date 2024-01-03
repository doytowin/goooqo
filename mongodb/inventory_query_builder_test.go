package mongodb

import "go.mongodb.org/mongo-driver/bson"

func (q InventoryQuery) BuildFilter() []bson.D {
	d := make([]bson.D, 0, 10)
	if q.QtyGt != nil {
		d = append(d, bson.D{{"qty", bson.D{{"$gt", q.QtyGt}}}})
	}
	return d
}
