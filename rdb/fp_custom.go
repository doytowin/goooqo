/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2025, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package rdb

import (
	. "github.com/doytowin/goooqo/core"
	"reflect"
	"strings"
)

type fpCustom struct {
	field     *reflect.StructField
	condition string
	phCnt     int
}

func buildFpCustom(field reflect.StructField) FieldProcessor {
	condition := field.Tag.Get("condition")
	phCnt := strings.Count(condition, "?")
	return &fpCustom{&field, condition, phCnt}
}

func (fp *fpCustom) Process(value reflect.Value) (string, []any) {
	arr := make([]any, 0, fp.phCnt)
	arg := ReadValue(value)
	for j := 0; j < fp.phCnt; j++ {
		arr = append(arr, arg)
	}
	return fp.condition, arr
}
