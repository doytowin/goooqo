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

var sqRegx = regexp.MustCompile(`(?i)(select|from)[\s:]([\w()]+)`)

func BuildBySubqueryTag(subqueryStr string, fieldName string) *fpSubquery {
	fp := &fpSubquery{}
	submatch := sqRegx.FindAllStringSubmatch(subqueryStr, -1)
	for _, group := range submatch {
		if strings.EqualFold(group[1], "select") {
			fp.select_ = group[2]
		} else if strings.EqualFold(group[1], "from") {
			fp.from = group[2]
		}
	}
	fp.buildComp(fieldName)
	return fp
}

func BuildBySelectTag(tag reflect.StructTag, fieldName string) *fpSubquery {
	fp := &fpSubquery{}
	fp.select_ = tag.Get("select")
	fp.from = tag.Get("from")
	fp.buildComp(fieldName)
	return fp
}

var subOfRgx = regexp.MustCompile("(\\w+(Any|All|" + core.SuffixStr + "))(([A-Z]\\w+)Of([A-Z]\\w+))")
var aggregateRgx = regexp.MustCompile("(Avg|Max|Min|Sum|First|Last|Push)(\\w+)")

func BuildByFieldName(match []string) *fpSubquery {
	fp := &fpSubquery{}
	fp.select_ = convertForAggColumn(match[4])
	fp.from = match[5]
	fp.buildComp(match[1])
	return fp
}

func convertForAggColumn(col string) string {
	if match := aggregateRgx.FindStringSubmatch(col); len(match) > 0 {
		return strings.ToUpper(match[1]) + "(" + core.ConvertToColumnCase(match[2]) + ")"
	}
	return core.ConvertToColumnCase(col)
}

func (fp *fpSubquery) buildComp(fieldName string) {
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

var fieldRgx = regexp.MustCompile("(\\w+(" + core.SuffixStr + "))[A-Z\\d]")

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
