package main

import . "go.mongodb.org/mongo-driver/bson/primitive"

func (q InventoryQuery) BuildFilter() []D {
	d := make([]D, 0, 4)
	if q.Id != nil {
		d = append(d, D{{"_id", D{{"$eq", q.Id}}}})
	}
	if q.IdNot != nil {
		d = append(d, D{{"_id", D{{"$ne", q.IdNot}}}})
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
		if q.Size.HLt != nil {
			d = append(d, D{{"size.h", D{{"$lt", q.Size.HLt}}}})
		}
		if q.Size.HGe != nil {
			d = append(d, D{{"size.h", D{{"$gte", q.Size.HGe}}}})
		}
		if q.Size.Unit != nil {
			if q.Size.Unit.Name != nil {
				d = append(d, D{{"size.unit.name", D{{"$eq", q.Size.Unit.Name}}}})
			}
			if q.Size.Unit.NameNull {
				d = append(d, D{{"size.unit.name", D{{"$type", 10}}}})
			}
		}
	}
	if q.StatusNull {
		d = append(d, D{{"status", D{{"$type", 10}}}})
	}
	if q.StatusNotNull {
		d = append(d, D{{"status", D{{"$not", D{{"$type", 10}}}}}})
	}
	return d
}
