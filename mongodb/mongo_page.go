/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2025, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package mongodb

import (
	"github.com/doytowin/goooqo/core"
	. "go.mongodb.org/mongo-driver/bson/primitive"
)

func buildSort(sort string) D {
	submatch := core.SortRgx.FindAllStringSubmatch(sort, -1)
	result := make(D, len(submatch))
	for i, group := range submatch {
		if group[3] != "" {
			result[i] = E{group[1], 7 - len(group[3])*2}
		} else {
			result[i] = E{group[1], 1}
		}
	}
	return result
}
