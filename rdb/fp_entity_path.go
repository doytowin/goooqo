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
	. "github.com/doytowin/goooqo/core"
	"reflect"
	"strings"
)

type fpEntityPath struct {
	EntityPath
}

func buildFpEntityPath(field reflect.StructField) FieldProcessor {
	return &fpEntityPath{*BuildEntityPath(field)}
}

func (fp *fpEntityPath) Process(value reflect.Value) (string, []any) {
	args := make([]any, 0)

	l := len(fp.Path)
	sql := fp.LocalField + " IN ("
	closeParesis := strings.Repeat(")", l)
	for i := 0; i < l-1; i++ {
		queryValue := value.FieldByName(Capitalize(fp.Path[i]) + "Query")
		if queryValue.IsValid() && !queryValue.IsNil() {
			where0, args0 := BuildWhereClause(queryValue.Interface())
			sql += "SELECT id FROM " + FormatTable(fp.Path[i]) + where0 + "\nINTERSECT "
			args = append(args, args0...)
		}
		sql += "SELECT " + fp.JoinIds[i] + " FROM " + fp.JoinTables[i] + " WHERE " + fp.JoinIds[i+1] + " IN ("
	}
	where, args0 := BuildWhereClause(value.Interface())
	args = append(args, args0...)
	return sql + "SELECT " + fp.ForeignField + " FROM " + fp.TargetTable + where + closeParesis, args
}

func (fp *fpEntityPath) buildSql(columns string) string {
	return "SELECT " + columns +
		" FROM " + fp.TargetTable +
		" WHERE " + fp.LocalField +
		" IN (" +
		"SELECT " + fp.JoinIds[1] +
		" FROM " + fp.JoinTables[0] +
		" WHERE " + fp.JoinIds[0] + " = ?)"
}
