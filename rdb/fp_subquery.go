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
	"regexp"
	"strings"
)

type fpSubquery struct {
	column, sign  string
	select_, from string
}

func (fp *fpSubquery) Process(value reflect.Value) (string, []any) {
	where, args := BuildWhereClause(value.Interface())
	return fp.buildCondition(where), args
}

func (fp *fpSubquery) Subquery() string {
	if em := emMap[fp.from]; em != nil {
		fp.from = em.TableName
	}
	return fp.column + fp.sign + "(SELECT " + fp.select_ + " FROM " + fp.from
}

func (fp *fpSubquery) buildCondition(where string) string {
	return fp.Subquery() + where + ")"
}

var sqRegx = regexp.MustCompile(`(select|from):([\w()]+)`)

func buildFpSubquery(field reflect.StructField) *fpSubquery {
	subqueryStr := field.Tag.Get("subquery")
	return BuildSubquery(subqueryStr, field.Name)
}

func BuildSubquery(subqueryStr string, fieldName string) (fp *fpSubquery) {
	fp = &fpSubquery{}
	submatch := sqRegx.FindAllStringSubmatch(subqueryStr, -1)
	for _, group := range submatch {
		if group[1] == "select" {
			fp.select_ = group[2]
		} else if group[1] == "from" {
			fp.from = group[2]
		}
	}
	anyAll := ""
	if strings.HasSuffix(fieldName, "Any") {
		anyAll = "ANY"
		fieldName = strings.TrimSuffix(fieldName, "Any")
	} else if strings.HasSuffix(fieldName, "All") {
		anyAll = "ALL"
		fieldName = strings.TrimSuffix(fieldName, "All")
	}
	fieldName = trimFieldName(fieldName)
	FpSuffix := buildFpSuffix(fieldName)
	fp.column, fp.sign = FpSuffix.col, FpSuffix.op.sign+anyAll
	return
}

var fieldRgx = regexp.MustCompile("(\\w+" + suffixStr + ")[A-Z\\d]")

// Trim tailing string after the predicate suffix.
// Examples:
// `ScoreLtAvg` -> `ScoreLt`
// `ScoreGeAvg` -> `ScoreGe`
func trimFieldName(fieldName string) string {
	if match := fieldRgx.FindStringSubmatch(fieldName); len(match) > 0 {
		fieldName = match[1]
	}
	return fieldName
}
