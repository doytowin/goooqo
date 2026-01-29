/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2026, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package rdb

import (
	"strings"

	"github.com/doytowin/goooqo/core"
)

func BuildSortClause(sort string) string {
	if strings.TrimSpace(sort) == "" {
		return ""
	}
	groups := core.SortRgx.FindAllStringSubmatch(sort, -1)
	var orderBy = make([]string, len(groups))
	for i, group := range groups {
		orderBy[i] = group[1]
		if group[3] != "" {
			orderBy[i] += " " + strings.ToUpper(group[3])
		}
	}
	return " ORDER BY " + strings.Join(orderBy, ", ")
}
