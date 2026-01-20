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
	"reflect"
	"strings"

	. "github.com/doytowin/goooqo/core"
)

type fpEntityPath struct {
	EntityPath
}

func buildFpEntityPath(field reflect.StructField) FieldProcessor {
	return &fpEntityPath{*BuildEntityPath(field)}
}

func BuildRelationEntityPath(field reflect.StructField) fpEntityPath {
	return fpEntityPath{*BuildEntityPath(field)}
}

func (fp *fpEntityPath) Process(value reflect.Value) (string, []any) {
	args := make([]any, 0)

	l := len(fp.Relations)
	sql := fp.Base.Fk1 + " IN ("
	closeParesis := strings.Repeat(")", l+1)
	for i := l - 1; i >= 0; i-- {
		relation := fp.Relations[i]
		sql += "SELECT " + relation.Fk2 + " FROM " + relation.At + " WHERE " + relation.Fk1 + " IN ("
		queryValue := value.FieldByName(Capitalize(fp.Path[i]) + "Query")
		if queryValue.IsValid() && !queryValue.IsNil() {
			where0, args0 := BuildWhereClause(queryValue.Interface())
			sql += "SELECT id FROM " + FormatTable(fp.Path[i]) + where0 + "\nINTERSECT "
			args = append(args, args0...)
		}
	}
	where, args0 := BuildWhereClause(value.Interface())
	args = append(args, args0...)
	e1, _, _ := strings.Cut(fp.Path[0], "->")
	return sql + "SELECT " + fp.Base.Fk2 + " FROM " + FormatTable(e1) + where + closeParesis, args
}

func buildColumns(fieldMetas []FieldMetadata) string {
	columns := make([]string, 0, len(fieldMetas))
	for _, md := range fieldMetas {
		if md.EntityPath == nil {
			columns = append(columns, md.ColumnName)
		}
	}
	return strings.Join(columns, ", ")
}

func (fp *fpEntityPath) buildQuery(query Query) (string, []any) {
	fieldMetas := BuildFieldMetas(fp.EntityType)
	columns := buildColumns(fieldMetas)

	s := fp.buildSql(columns)
	and, args := BuildConditions(query, " AND ", " AND ", "")
	s += and + BuildSortClause(query.GetSort())
	if query.NeedPaging() {
		s = BuildPageClause(&s, query.CalcOffset(), query.GetPageSize())
	}
	return s, args
}

func (fp *fpEntityPath) buildSql(columns string) string {
	l := len(fp.Relations)
	closeParesis := strings.Repeat(")", l)
	sql := "SELECT " + columns + " FROM " + fp.Base.At + " WHERE " + fp.Base.Fk2
	for i := l - 1; i >= 0; i-- {
		relation := fp.Relations[i]
		sql += " IN (" + "SELECT " + relation.Fk2 + " FROM " + relation.At + " WHERE " + relation.Fk1
	}
	return sql + " = ?" + closeParesis
}
