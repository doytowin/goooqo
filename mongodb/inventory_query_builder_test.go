package mongodb

import . "go.mongodb.org/mongo-driver/bson/primitive"

func (q InventoryQuery) BuildFilter(connector string) D {
	d := make(A, 0, 4)
	if q.QtyGt != nil {
		d = append(d, D{{"qty", D{{"$gt", q.QtyGt}}}})
	}
	if q.QtyOr != nil {
		d = append(d, q.QtyOr.BuildFilter("$or"))
	}
	return CombineConditions(connector, d)
}

func (q QtyOr) BuildFilter(connector string) D {
	d := make(A, 0, 4)
	if q.QtyLt != nil {
		d = append(d, D{{"qty", D{{"$lt", q.QtyLt}}}})
	}
	if q.QtyGe != nil {
		d = append(d, D{{"qty", D{{"$gte", q.QtyGe}}}})
	}
	return CombineConditions(connector, d)
}
