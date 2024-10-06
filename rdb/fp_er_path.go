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
	"github.com/doytowin/goooqo/core"
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
		joinIds[i] = FormatJoinId(domain)
	}
	joinTables := make([]string, l-1)
	for i := 0; i < l-1; i++ {
		joinTables[i] = FormatJoinTable(path[i], path[i+1])
	}
	targetTable := FormatTable(path[l-1])
	localFieldColumn := core.ConvertToColumnCase(field.Tag.Get("localField"))
	if localFieldColumn == "" {
		localFieldColumn = "id"
	}
	foreignFieldColumn := core.ConvertToColumnCase(field.Tag.Get("foreignField"))
	if foreignFieldColumn == "" {
		foreignFieldColumn = "id"
	}
	return &fpERPath{ERPath{path, joinTables, joinIds, targetTable, localFieldColumn, foreignFieldColumn}}
}

func (fp *fpERPath) Process(value reflect.Value) (string, []any) {
	args := make([]any, 0)

	l := len(fp.path)
	sql := fp.localField + " IN ("
	closeParesis := strings.Repeat(")", l)
	for i := 0; i < l-1; i++ {
		queryValue := value.FieldByName(core.Capitalize(fp.path[i]) + "Query")
		if queryValue.IsValid() && !queryValue.IsNil() {
			where0, args0 := BuildWhereClause(queryValue.Interface())
			sql += "SELECT id FROM " + FormatTable(fp.path[i]) + where0 + "\nINTERSECT "
			args = append(args, args0...)
		}
		sql += "SELECT " + fp.joinIds[i] + " FROM " + fp.joinTables[i] + " WHERE " + fp.joinIds[i+1] + " IN ("
	}
	where, args0 := BuildWhereClause(value.Interface())
	args = append(args, args0...)
	return sql + "SELECT " + fp.foreignField + " FROM " + fp.targetTable + where + closeParesis, args
}
