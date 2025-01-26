/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2025, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

import (
	"github.com/doytowin/goooqo/rdb"
	log "github.com/sirupsen/logrus"
	"go/ast"
	"reflect"
	"strings"
)

const format = "conditions = append(conditions, \"%s %s ?\")"

type SqlGenerator struct {
	*generator
}

func init() {
	sqlOpMap := make(map[string]operator)
	sqlOpMap["Eq"] = operator{name: "Eq", sign: "="}
	sqlOpMap["Ne"] = operator{name: "Ne", sign: "<>"}
	sqlOpMap["Gt"] = operator{name: "Gt", sign: ">"}
	sqlOpMap["Ge"] = operator{name: "Ge", sign: ">="}
	sqlOpMap["Lt"] = operator{name: "Lt", sign: "<"}
	sqlOpMap["Le"] = operator{name: "Le", sign: "<="}
	sqlOpMap["In"] = operator{name: "In", sign: "IN", format: "conditions = append(conditions, \"%s %s (\"+strings.Join(phs, \", \")+\")\")"}
	sqlOpMap["NotIn"] = operator{name: "NotIn", sign: "NOT IN", format: "conditions = append(conditions, \"%s %s (\"+strings.Join(phs, \", \")+\")\")"}
	sqlOpMap["Null"] = operator{name: "Null", sign: "IS NULL", format: "conditions = append(conditions, \"%s %s\")"}
	sqlOpMap["Like"] = operator{name: "Like", sign: "LIKE", format: format}
	sqlOpMap["NotLike"] = operator{name: "NotLike", sign: "NOT LIKE", format: format}
	sqlOpMap["Contain"] = operator{name: "Contain", sign: "LIKE", format: format}
	sqlOpMap["NotContain"] = operator{name: "NotContain", sign: "NOT LIKE", format: format}
	sqlOpMap["Start"] = operator{name: "Start", sign: "LIKE", format: format}
	sqlOpMap["NotStart"] = operator{name: "NotStart", sign: "NOT LIKE", format: format}
	sqlOpMap["End"] = operator{name: "End", sign: "LIKE", format: format}
	sqlOpMap["NotEnd"] = operator{name: "NotEnd", sign: "NOT LIKE", format: format}
	opMap["sql"] = sqlOpMap
}

func NewSqlGenerator() *SqlGenerator {
	return &SqlGenerator{newGenerator("sql",
		[]string{`. "github.com/doytowin/goooqo/rdb"`, `"strings"`}, format,
	)}
}

func (g *SqlGenerator) appendBuildMethod(ts *ast.TypeSpec) {
	g.WriteString(NewLine)
	g.writeInstruction("func (q %s) BuildConditions() ([]string, []any) {", ts.Name)
	g.appendFuncBody(ts)
	g.writeInstruction("}")
}

func (g *SqlGenerator) appendFuncBody(ts *ast.TypeSpec) {
	intent := g.incIntent()
	g.writeInstruction("conditions := make([]string, 0, 4)")
	g.writeInstruction("args := make([]any, 0, 4)")
	for _, field := range ts.Type.(*ast.StructType).Fields.List {
		if field.Names != nil {
			g.appendCondition(field, field.Names[0].Name)
		}
	}
	g.writeInstruction("return conditions, args")
	g.restoreIntent(intent)
}

func (g *SqlGenerator) appendStruct(stp *ast.StructType) {
	for _, field := range stp.Fields.List {
		if field.Names != nil {
			g.appendCondition(field, field.Names[0].Name)
		}
	}
}

func (g *SqlGenerator) appendCondition(field *ast.Field, fieldName string) {
	column, op := g.suffixMatch(fieldName)

	if field.Tag != nil {
		g.appendIfStartNil(fieldName)
		tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`"))
		if subqueryTag, ok := tag.Lookup("subquery"); ok {
			fpSubquery := rdb.BuildBySubqueryTag(subqueryTag, fieldName)
			subSelect := fpSubquery.Subquery()
			g.genSubquery(fieldName, subSelect)
		} else if _, ok = tag.Lookup("select"); ok {
			fpSubquery := rdb.BuildBySelectTag(tag, fieldName)
			subSelect := fpSubquery.Subquery()
			g.genSubquery(fieldName, subSelect)
		} else if conditionTag, ok := tag.Lookup("condition"); ok {
			g.appendIfBody("conditions = append(conditions, \"%s\")", conditionTag)
			for i := 0; i < strings.Count(conditionTag, "?"); i++ {
				g.appendArg(fieldName)
			}
		} else {
			log.Warn("Unsupported field: ", fieldName, " ", field.Type)
		}
	} else if strings.Contains(op.sign, "NULL") {
		g.appendIfStartNil(fieldName)
		intent := g.incIntent()
		g.writeInstruction("if *q.%s {", fieldName)
		g.appendIfBody(op.format, column, "IS NULL")
		g.writeInstruction("} else {")
		g.appendIfBody(op.format, column, "IS NOT NULL")
		g.appendIfEnd()
		g.restoreIntent(intent)
	} else if strings.Contains(op.sign, "IN") {
		g.appendIfStartNil(fieldName)
		g.appendIfBody("phs := make([]string, 0, len(*q.%s))", fieldName)
		g.appendIfBody("for _, arg := range *q.%s {", fieldName)
		g.appendIfBody("\targs = append(args, arg)")
		g.appendIfBody("\tphs = append(phs, \"?\")")
		g.appendIfBody("}")
		g.appendIfBody(op.format, column, op.sign)
	} else if strings.HasSuffix(fieldName, "Or") {
		g.appendIfStartNil(fieldName)
		g.appendIfBody("cond, args0 := BuildConditions(q.%s, \"(\", \" OR \", \")\")", fieldName)
		g.appendIfBody("conditions = append(conditions, cond)")
		g.appendIfBody("args = append(args, args0...)")
	} else if strings.HasSuffix(op.sign, "LIKE") {
		g.appendIfStartNil(fieldName)
		g.appendIfBody(op.format, column, op.sign)
		if strings.HasSuffix(op.name, "Contain") {
			g.appendIfBody("args = append(args, \"%%\"+*q.%s+\"%%\")", fieldName)
		} else if strings.HasSuffix(op.name, "Start") {
			g.appendIfBody("args = append(args, *q.%s+\"%%\")", fieldName)
		} else if strings.HasSuffix(op.name, "End") {
			g.appendIfBody("args = append(args, \"%%\"+*q.%s)", fieldName)
		} else {
			g.appendArg(fieldName)
		}
	} else {
		g.appendIfStartNil(fieldName)
		g.appendIfBody(op.format, column, op.sign)
		g.appendArg(fieldName)
	}
	g.appendIfEnd()
}

func (g *SqlGenerator) genSubquery(fieldName string, subSelect string) {
	g.appendIfBody("where, args1 := BuildWhereClause(q.%s)", fieldName)
	g.appendIfBody("conditions = append(conditions, \"" + subSelect + "\"+where+\")\")")
	g.appendIfBody("args = append(args, args1...)")
}

func (g *SqlGenerator) appendArg(fieldName string) {
	g.appendIfBody("args = append(args, *q.%s)", fieldName)
}
