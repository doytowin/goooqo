/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2026, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

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
