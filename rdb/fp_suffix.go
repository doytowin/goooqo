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
	"bytes"
	. "github.com/doytowin/goooqo/core"
	"reflect"
	"regexp"
	"strings"
)

var (
	opMap     = CreateOpMap()
	suffixRgx = regexp.MustCompile(`(Gt|Ge|Lt|Le|Not|Ne|Eq|Null|NotIn|In|Like|NotLike|Contain|NotContain|Start|NotStart|End|NotEnd|Rx)$`)
	escapeRgx = regexp.MustCompile("[\\\\_%]")
)

type operator struct {
	name, sign string
	process    func(value reflect.Value) (string, []any)
}

func ReadValueToArray(value reflect.Value) (string, []any) {
	return "?", []any{ReadValue(value)}
}

func ReadValueForIn(value reflect.Value) (string, []any) {
	arg := reflect.Indirect(value)
	args := make([]any, 0, arg.Len())
	ph := bytes.NewBuffer(make([]byte, 0, 3*arg.Len()))
	ph.WriteString("(")
	for i := 0; i < arg.Len(); i++ {
		args = append(args, arg.Index(i).Interface())
		ph.WriteString("?")
		if i < arg.Len()-1 {
			ph.WriteString(", ")
		}
	}
	ph.WriteString(")")
	return ph.String(), args
}

func ReadLikeValue(value reflect.Value) string {
	s := ReadValue(value).(string)
	return escapeRgx.ReplaceAllString(s, "\\$0")
}

func CreateOpMap() map[string]operator {
	const Like = " LIKE "
	const NotLike = " NOT LIKE "
	opMap := make(map[string]operator)
	opMap["Gt"] = operator{"Gt", " > ", ReadValueToArray}
	opMap["Ge"] = operator{"Ge", " >= ", ReadValueToArray}
	opMap["Lt"] = operator{"Lt", " < ", ReadValueToArray}
	opMap["Le"] = operator{"Le", " <= ", ReadValueToArray}
	opMap["Not"] = operator{"Not", " != ", ReadValueToArray}
	opMap["Ne"] = operator{"Ne", " <> ", ReadValueToArray}
	opMap["Eq"] = operator{"Eq", " = ", ReadValueToArray}
	opMap["Null"] = operator{"Null", "", func(rv reflect.Value) (string, []any) {
		if rv.Bool() == false {
			return " IS NOT NULL", []any{}
		}
		return " IS NULL", []any{}
	}}
	opMap["In"] = operator{"In", " IN ", ReadValueForIn}
	opMap["NotIn"] = operator{"NotIn", " NOT IN ", ReadValueForIn}
	opMap["Like"] = operator{"Like", Like, func(value reflect.Value) (string, []any) {
		s := ReadValue(value).(string)
		ph := resolvePlaceHolder(s)
		return ph, []any{s}
	}}
	opMap["NotLike"] = operator{"NotLike", NotLike, func(value reflect.Value) (string, []any) {
		s := ReadValue(value).(string)
		ph := resolvePlaceHolder(s)
		return ph, []any{s}
	}}
	opMap["Contain"] = operator{"Contain", Like, func(value reflect.Value) (string, []any) {
		escape := ReadLikeValue(value)
		ph := resolvePlaceHolder(escape)
		return ph, []any{"%" + escape + "%"}
	}}
	opMap["NotContain"] = operator{"NotContain", NotLike, func(value reflect.Value) (string, []any) {
		escape := ReadLikeValue(value)
		ph := resolvePlaceHolder(escape)
		return ph, []any{"%" + escape + "%"}
	}}
	opMap["Start"] = operator{"Start", Like, func(value reflect.Value) (string, []any) {
		escape := ReadLikeValue(value)
		ph := resolvePlaceHolder(escape)
		return ph, []any{escape + "%"}
	}}
	opMap["NotStart"] = operator{"NotStart", NotLike, func(value reflect.Value) (string, []any) {
		escape := ReadLikeValue(value)
		ph := resolvePlaceHolder(escape)
		return ph, []any{escape + "%"}
	}}
	opMap["End"] = operator{"End", Like, func(value reflect.Value) (string, []any) {
		escape := ReadLikeValue(value)
		ph := resolvePlaceHolder(escape)
		return ph, []any{"%" + escape}
	}}
	opMap["NotEnd"] = operator{"NotEnd", NotLike, func(value reflect.Value) (string, []any) {
		escape := ReadLikeValue(value)
		ph := resolvePlaceHolder(escape)
		return ph, []any{"%" + escape}
	}}
	opMap["Rx"] = operator{"Rx", " REGEXP ", ReadValueToArray}
	return opMap
}

func resolvePlaceHolder(arg string) string {
	ph := "?"
	if strings.Contains(arg, "\\") {
		ph = ph + " ESCAPE '\\'"
	}
	return ph
}

type fpSuffix struct {
	col string
	op  operator
}

func buildFpSuffix(fieldName string) fpSuffix {
	if match := suffixRgx.FindStringSubmatch(fieldName); len(match) > 0 {
		op := opMap[match[1]]
		column := strings.TrimSuffix(fieldName, match[1])
		column = ConvertToColumnCase(column)
		return fpSuffix{column, op}
	}
	return fpSuffix{ConvertToColumnCase(fieldName), opMap["Eq"]}
}

func (fp fpSuffix) Process(value reflect.Value) (string, []any) {
	placeholder, args := fp.op.process(value)
	return fp.col + fp.op.sign + placeholder, args
}
