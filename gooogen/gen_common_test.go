/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2025, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

import (
	"testing"
)

func TestExampleCommentMap(t *testing.T) {
	tests := []struct {
		input, output, expect string
		generator             Generator
	}{
		{input: "../main/inventory.go", output: "../main/inventory_query_builder.go", expect: `package main

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
`, generator: NewMongoGenerator()},
		{input: "../main/user.go", output: "../main/user_query_builder.go", expect: `package main

import . "github.com/doytowin/goooqo/rdb"
import "strings"

func (q UserQuery) BuildConditions() ([]string, []any) {
	conditions := make([]string, 0, 4)
	args := make([]any, 0, 4)
	if q.IdGt != nil {
		conditions = append(conditions, "id > ?")
		args = append(args, *q.IdGt)
	}
	if q.IdIn != nil {
		phs := make([]string, 0, len(*q.IdIn))
		for _, arg := range *q.IdIn {
			args = append(args, arg)
			phs = append(phs, "?")
		}
		conditions = append(conditions, "id IN ("+strings.Join(phs, ", ")+")")
	}
	if q.IdNotIn != nil {
		phs := make([]string, 0, len(*q.IdNotIn))
		for _, arg := range *q.IdNotIn {
			args = append(args, arg)
			phs = append(phs, "?")
		}
		conditions = append(conditions, "id NOT IN ("+strings.Join(phs, ", ")+")")
	}
	if q.Cond != nil {
		conditions = append(conditions, "(score = ? OR memo = ?)")
		args = append(args, *q.Cond)
		args = append(args, *q.Cond)
	}
	if q.ScoreLt != nil {
		conditions = append(conditions, "score < ?")
		args = append(args, *q.ScoreLt)
	}
	if q.MemoNull != nil {
		if *q.MemoNull {
			conditions = append(conditions, "memo IS NULL")
		} else {
			conditions = append(conditions, "memo IS NOT NULL")
		}
	}
	if q.MemoLike != nil {
		conditions = append(conditions, "memo LIKE ?")
		args = append(args, *q.MemoLike)
	}
	if q.Deleted != nil {
		conditions = append(conditions, "deleted = ?")
		args = append(args, *q.Deleted)
	}
	if q.MemoContain != nil {
		conditions = append(conditions, "memo LIKE ?")
		args = append(args, "%" + *q.MemoContain + "%")
	}
	if q.MemoNotContain != nil {
		conditions = append(conditions, "memo NOT LIKE ?")
		args = append(args, "%" + *q.MemoNotContain + "%")
	}
	if q.MemoStart != nil {
		conditions = append(conditions, "memo LIKE ?")
		args = append(args, *q.MemoStart + "%")
	}
	if q.MemoNotStart != nil {
		conditions = append(conditions, "memo NOT LIKE ?")
		args = append(args, *q.MemoNotStart + "%")
	}
	if q.ScoreLtAvg != nil {
		where, args1 := BuildWhereClause(q.ScoreLtAvg)
		conditions = append(conditions, "score < (SELECT avg(score) FROM t_user"+where+")")
		args = append(args, args1...)
	}
	if q.ScoreLtAny != nil {
		where, args1 := BuildWhereClause(q.ScoreLtAny)
		conditions = append(conditions, "score < ANY(SELECT score FROM t_user"+where+")")
		args = append(args, args1...)
	}
	if q.ScoreLtAll != nil {
		where, args1 := BuildWhereClause(q.ScoreLtAll)
		conditions = append(conditions, "score < ALL(SELECT score FROM t_user"+where+")")
		args = append(args, args1...)
	}
	if q.ScoreGtAvg != nil {
		where, args1 := BuildWhereClause(q.ScoreGtAvg)
		conditions = append(conditions, "score > (SELECT avg(score) FROM t_user"+where+")")
		args = append(args, args1...)
	}
	return conditions, args
}
`, generator: NewSqlGenerator()},
	}
	for _, tt := range tests {
		t.Run("Generate for "+tt.input, func(t *testing.T) {
			code := GenerateCode(tt.input, tt.generator)
			if code != tt.expect {
				t.Fatalf("Got \n%s", code)
			}
			_ = WriteFile(tt.output, code)
		})
	}
}
