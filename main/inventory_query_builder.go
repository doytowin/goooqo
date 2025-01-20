package main

import . "go.mongodb.org/mongo-driver/bson/primitive"
import . "github.com/doytowin/goooqo/mongodb"

func (q InventoryQuery) BuildFilter(connector string) D {
	d := make(A, 0, 4)
	if q.Id != nil {
		d = append(d, D{{"_id", D{{"$eq", q.Id}}}})
	}
	if q.IdNe != nil {
		d = append(d, D{{"_id", D{{"$ne", q.IdNe}}}})
	}
	if q.IdIn != nil {
		d = append(d, D{{"_id", D{{"$in", q.IdIn}}}})
	}
	if q.IdNotIn != nil {
		d = append(d, D{{"_id", D{{"$nin", q.IdNotIn}}}})
	}
	if q.Qty != nil {
		d = append(d, D{{"qty", D{{"$eq", q.Qty}}}})
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
	if q.Size != nil {
		d = append(d, q.Size.BuildFilter("$and"))
	}
	if q.StatusNull != nil {
		if *q.StatusNull {
			d = append(d, D{{"status", D{{"$type", 10}}}})
		} else {
			d = append(d, D{{"status", D{{"$not", D{{"$type", 10}}}}}})
		}
	}
	if q.ItemContain != nil && *q.ItemContain != "" {
		d = append(d, D{{"item", D{{"$regex", q.ItemContain}}}})
	}
	if q.ItemNotContain != nil && *q.ItemNotContain != "" {
		d = append(d, D{{"item", D{{"$not", D{{"$regex", q.ItemNotContain}}}}}})
	}
	if q.ItemStart != nil && *q.ItemStart != "" {
		d = append(d, D{{"item", D{{"$regex", "^" + *q.ItemStart}}}})
	}
	if q.ItemNotStart != nil && *q.ItemNotStart != "" {
		d = append(d, D{{"item", D{{"$not", D{{"$regex", "^" + *q.ItemNotStart}}}}}})
	}
	if q.ItemEnd != nil && *q.ItemEnd != "" {
		d = append(d, D{{"item", D{{"$regex", *q.ItemEnd + "$"}}}})
	}
	if q.ItemNotEnd != nil && *q.ItemNotEnd != "" {
		d = append(d, D{{"item", D{{"$not", D{{"$regex", *q.ItemNotEnd + "$"}}}}}})
	}
	if q.CustomFilter != nil {
		d = append(d, *q.CustomFilter)
	}
	if q.QtyOr != nil {
		d = append(d, q.QtyOr.BuildFilter("$or"))
	}
	if q.Search != nil {
		d = append(d, D{{"$text", D{{"$search", *q.Search}}}})
	}
	return CombineConditions(connector, d)
}

func (q SizeQuery) BuildFilter(connector string) D {
	d := make(A, 0, 4)
	if q.HLt != nil {
		d = append(d, D{{"size.h", D{{"$lt", q.HLt}}}})
	}
	if q.HGe != nil {
		d = append(d, D{{"size.h", D{{"$gte", q.HGe}}}})
	}
	if q.Unit != nil {
		d = append(d, q.Unit.BuildFilter("$and"))
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
	if q.Size != nil {
		d = append(d, q.Size.BuildFilter("$and"))
	}
	if q.SizeOr != nil {
		d = append(d, q.SizeOr.BuildFilter("$or"))
	}
	return CombineConditions(connector, d)
}

func (q Unit) BuildFilter(connector string) D {
	d := make(A, 0, 4)
	if q.Name != nil {
		d = append(d, D{{"size.unit.name", D{{"$eq", q.Name}}}})
	}
	if q.NameNull != nil {
		if *q.NameNull {
			d = append(d, D{{"size.unit.name", D{{"$type", 10}}}})
		} else {
			d = append(d, D{{"size.unit.name", D{{"$not", D{{"$type", 10}}}}}})
		}
	}
	return CombineConditions(connector, d)
}
