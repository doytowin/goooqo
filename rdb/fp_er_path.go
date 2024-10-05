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
	"reflect"
	"strings"
)

type ERPath struct {
	path         []string
	joinTables   []string
	joinIds      []string
	targetTable  string
	localField   string
	foreignField string
}

type fpERPath struct {
	ERPath
}

func buildFpErPath(field reflect.StructField) FieldProcessor {
	path := strings.Split(field.Tag.Get("erpath"), ",")
	l := len(path)
	joinIds := make([]string, l)
	for i, domain := range path {
		joinIds[i] = fmt.Sprintf("%s_id", domain)
	}
	joinTables := make([]string, l-1)
	for i := 0; i < l-1; i++ {
		joinTables[i] = fmt.Sprintf("a_%s_and_%s", path[i], path[i+1])
	}
	targetTable := fmt.Sprintf("t_%s", path[l-1])
	return &fpERPath{ERPath{path, joinTables, joinIds, targetTable, "id", "id"}}
}

func (fp *fpERPath) Process(value reflect.Value) (condition string, args []any) {
	where, args := BuildWhereClause(value.Interface())
	return fp.localField + " IN (SELECT " + fp.joinIds[0] + " FROM " + fp.joinTables[0] + " WHERE " +
		fp.joinIds[1] + " IN (SELECT " + fp.foreignField + " FROM " + fp.targetTable + where + "))", args
}
