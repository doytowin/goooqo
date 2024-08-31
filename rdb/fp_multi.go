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
	"reflect"
	"strings"
)

// build multiple conditions mapped by a struct and connect them
type fpMultiConditions struct {
	connect func([]string) string
}

var fpForOr = &fpMultiConditions{connect: func(conditions []string) string {
	return "(" + strings.Join(conditions, " OR ") + ")"
}}

var fpForAnd = &fpMultiConditions{connect: func(conditions []string) string {
	return strings.Join(conditions, " AND ")
}}

func (fp *fpMultiConditions) Process(value reflect.Value) (string, []any) {
	conditions, args := buildConditions(value.Interface())
	return fp.connect(conditions), args
}
