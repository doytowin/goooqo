/*
 * The Clear BSD License
 *
 * Copyright (c) 2024, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package rdb

import (
	"fmt"
	"github.com/doytowin/goooqo/core"
	"strings"
)

func BuildPageClause(sql *string, offset int, size int) string {
	return fmt.Sprintf("%s LIMIT %d OFFSET %d", *sql, size, offset)
}

func BuildSortClause(sort *string) string {
	if sort == nil {
		return ""
	}
	groups := core.SortRgx.FindAllStringSubmatch(*sort, -1)
	var orderBy = make([]string, len(groups))
	for i, group := range groups {
		orderBy[i] = group[1]
		if group[3] != "" {
			orderBy[i] += " " + strings.ToUpper(group[3])
		}
	}
	return " ORDER BY " + strings.Join(orderBy, ", ")
}
