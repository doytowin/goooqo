/*
 * The Clear BSD License
 *
 * Copyright (c) 2024, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package mongodb

import . "go.mongodb.org/mongo-driver/bson/primitive"

func (q InventoryQuery) BuildFilter() A {
	d := make(A, 0, 4)
	if q.QtyGt != nil {
		d = append(d, D{{"qty", D{{"$gt", q.QtyGt}}}})
	}
	if q.QtyOr != nil {
		or := make(A, 0, 4)
		if q.QtyOr.QtyLt != nil {
			or = append(or, D{{"qty", D{{"$lt", q.QtyOr.QtyLt}}}})
		}
		if q.QtyOr.QtyGe != nil {
			or = append(or, D{{"qty", D{{"$gte", q.QtyOr.QtyGe}}}})
		}
		if len(or) > 1 {
			d = append(d, D{{"$or", or}})
		} else if len(or) == 1 {
			d = append(d, or[0])
		}
	}
	return d
}
