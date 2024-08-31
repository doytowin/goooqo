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
	"regexp"
	"strings"
)

func BuildPageClause(sql *string, offset int, size int) string {
	return fmt.Sprintf("%s LIMIT %d OFFSET %d", *sql, size, offset)
}

var sortRgx = regexp.MustCompile("(?i)(\\w+)(,(asC|dEsc))?;?")

func BuildSortClause(sort *string) string {
	if sort == nil {
		return ""
	}
	submatch := sortRgx.FindAllStringSubmatch(*sort, -1)
	var orderBy = make([]string, len(submatch))
	for i, group := range submatch {
		orderBy[i] = group[1]
		if group[3] != "" {
			orderBy[i] += " " + strings.ToUpper(group[3])
		}
	}
	return " ORDER BY " + strings.Join(orderBy, ", ")
}
