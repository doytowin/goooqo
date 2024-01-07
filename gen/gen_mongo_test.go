package gen

import (
	"testing"
)

func TestExampleCommentMap(t *testing.T) {
	tests := []struct{ input, output, expect string }{
		{input: "../main/inventory.go", output: "../main/inventory_query_builder.go", expect: `package main

import . "go.mongodb.org/mongo-driver/bson/primitive"

func (q InventoryQuery) BuildFilter() []D {
	d := make([]D, 0, 10)
	if q.Id != nil {
		d = append(d, D{{"_id", D{{"$eq", q.Id}}}})
	}
	if q.IdIn != nil {
		d = append(d, D{{"_id", D{{"$in", q.IdIn}}}})
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
	return d
}
`},
	}
	for _, tt := range tests {
		t.Run("Generate for "+tt.input, func(t *testing.T) {
			code := GenerateCode(tt.input)
			if code != tt.expect {
				t.Fatalf("Got \n%s", code)
			}
			_ = Generate(tt.input, tt.output)
		})
	}
}
