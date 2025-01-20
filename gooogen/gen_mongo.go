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
	"github.com/doytowin/goooqo/core"
	log "github.com/sirupsen/logrus"
	"go/ast"
	"reflect"
	"strings"
)

type MongoGenerator struct {
	*generator
}

func NewMongoGenerator() *MongoGenerator {
	return &MongoGenerator{newGenerator("mongo",
		[]string{`. "go.mongodb.org/mongo-driver/bson/primitive"`, `. "github.com/doytowin/goooqo/mongodb"`},
		"d = append(d, D{{\"%s\", D{{\"%s\", q.%s}}}})",
	)}
}

const regexSign = "$regex"

func init() {
	mongoOpMap := make(map[string]operator)
	mongoOpMap["Eq"] = operator{name: "Eq", sign: "$eq"}
	mongoOpMap["Ne"] = operator{name: "Ne", sign: "$ne"}
	mongoOpMap["Gt"] = operator{name: "Gt", sign: "$gt"}
	mongoOpMap["Ge"] = operator{name: "Ge", sign: "$gte"}
	mongoOpMap["Lt"] = operator{name: "Lt", sign: "$lt"}
	mongoOpMap["Le"] = operator{name: "Le", sign: "$lte"}
	mongoOpMap["In"] = operator{name: "In", sign: "$in"}
	mongoOpMap["NotIn"] = operator{name: "NotIn", sign: "$nin"}
	mongoOpMap["Null"] = operator{name: "Null", sign: "$type"}
	mongoOpMap["Contain"] = operator{
		name:   "Contain",
		sign:   regexSign,
		format: "d = append(d, D{{\"%s\", D{{\"%s\", q.%s}}}})",
	}
	mongoOpMap["NotContain"] = operator{
		name:   "NotContain",
		sign:   regexSign,
		format: "d = append(d, D{{\"%s\", D{{\"$not\", D{{\"%s\", q.%s}}}}}})",
	}
	mongoOpMap["Start"] = operator{
		name:   "Start",
		sign:   regexSign,
		format: "d = append(d, D{{\"%s\", D{{\"%s\", \"^\" + *q.%s}}}})",
	}
	mongoOpMap["NotStart"] = operator{
		name:   "NotStart",
		sign:   regexSign,
		format: "d = append(d, D{{\"%s\", D{{\"$not\", D{{\"%s\", \"^\" + *q.%s}}}}}})",
	}
	mongoOpMap["End"] = operator{
		name:   "End",
		sign:   regexSign,
		format: "d = append(d, D{{\"%s\", D{{\"%s\", *q.%s + \"$\"}}}})",
	}
	mongoOpMap["NotEnd"] = operator{
		name:   "NotEnd",
		sign:   regexSign,
		format: "d = append(d, D{{\"%s\", D{{\"$not\", D{{\"%s\", *q.%s + \"$\"}}}}}})",
	}
	opMap["mongo"] = mongoOpMap
}

func (g *MongoGenerator) appendBuildMethod(ts *ast.TypeSpec) {
	g.WriteString(NewLine)
	g.writeInstruction("func (q %s) BuildFilter(connector string) D {", ts.Name)
	g.appendFuncBody(ts)
	g.writeInstruction("}")
}

func (g *MongoGenerator) appendFuncBody(ts *ast.TypeSpec) {
	g.appendIfBody("d := make(A, 0, 4)")
	g.appendStruct(ts, []string{})
	g.appendIfBody("return CombineConditions(connector, d)")
}

func (g *MongoGenerator) appendStruct(ts *ast.TypeSpec, path []string) {
	stp := ts.Type.(*ast.StructType)
	g.intent = strings.Repeat("\t", len(path)+1)
	for _, field := range stp.Fields.List {
		if field.Names != nil {
			g.appendCondition(field, path, field.Names[0].Name)
		} else if tn := resolveTypeName(field.Type); strings.HasSuffix(tn, "Or") {
			g.appendCondition(field, path, strings.TrimPrefix(tn, "*"))
		} else {
			log.Info("[MongoGenerator#appendStruct] Unresolved TypeName: ", tn)
		}
	}
	g.intent = strings.Repeat("\t", len(path))
}

func buildNestedFieldName(path []string, fieldName string) string {
	return strings.Join(append(path, fieldName), ".")
}

func (g *MongoGenerator) buildNestedProperty(path []string, column string) string {
	props := make([]string, 0, len(path)+1)
	currentPrefix := g.prefix[g.structIdx-1]
	if currentPrefix != "" {
		props = append(props, currentPrefix)
	}
	for _, fieldName := range path {
		if strings.HasSuffix(fieldName, "Or") {
			continue
		}
		// LATER: determine property by tag
		props = append(props, core.ConvertToColumnCase(fieldName))
	}
	return strings.Join(append(props, column), ".")
}

func (g *MongoGenerator) appendCondition(field *ast.Field, path []string, fieldName string) {
	column, op := g.suffixMatch(fieldName)
	if column == "id" {
		column = "_id"
	}

	structName := buildNestedFieldName(path, fieldName)
	column = g.buildNestedProperty(path, column)

	if field.Tag != nil {
		if columnTag, ok := reflect.StructTag(strings.Trim(field.Tag.Value, "`")).Lookup("column"); ok {
			column = columnTag
		}
	}

	if ts := toTypePointer(field); ts != nil {
		g.appendSubStruct(ts, structName, fieldName, column)
	} else if op.sign == "$type" {
		g.appendIfStartNil(structName)
		g.writeInstruction("\tif *q.%s {", structName)
		g.writeInstruction(g.replaceIns("\t\td = append(d, D{{\"%s\", D{{\"$type\", 10}}}})"), column)
		g.writeInstruction("\t} else {")
		g.writeInstruction(g.replaceIns("\t\td = append(d, D{{\"%s\", D{{\"$not\", D{{\"$type\", 10}}}}}})"), column)
		g.writeInstruction("\t}")
	} else if op.sign == regexSign {
		g.writeInstruction("if q.%s != nil && *q.%s != \"\" {", structName, structName)
		g.appendIfBody(op.format, column, op.sign, structName)
	} else if fieldName == "Search" {
		g.appendIfStartNil(structName)
		g.appendIfBody("d = append(d, D{{\"$text\", D{{\"$search\", *q.%s}}}})", structName)
	} else {
		g.appendIfStartNil(structName)
		if resolveTypeName(field.Type) == "*M" {
			g.appendIfBody("d = append(d, *q.%s)", structName)
		} else {
			g.appendIfBody(op.format, column, op.sign, structName)
		}
	}
	g.appendIfEnd()
}

func (g *MongoGenerator) appendSubStruct(ts *ast.TypeSpec, structName string, fieldName string, column string) {
	g.appendIfStartNil(structName)
	if strings.HasSuffix(fieldName, "Or") {
		g.addStruct("", ts) // generate for type recursively
		g.writeInstruction("\td = append(d, q.%s.BuildFilter(\"$or\"))", structName)
	} else {
		g.addStruct(column, ts) // generate for type recursively
		g.writeInstruction("\td = append(d, q.%s.BuildFilter(\"$and\"))", structName)
	}
}

func resolveTypeName(expr ast.Expr) string {
	var stack []string
loop:
	for {
		switch x := expr.(type) {
		case *ast.StarExpr:
			expr = x.X
			stack = append(stack, "*")
		case *ast.ArrayType:
			expr = x.Elt
			stack = append(stack, "[]")
		case *ast.Ident:
			stack = append(stack, x.Name)
			break loop
		case *ast.SelectorExpr:
			stack = append(stack, x.Sel.Name)
			break loop
		default:
			break loop
		}
	}
	return strings.Join(stack, "")
}
