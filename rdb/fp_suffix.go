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

var opMap = CreateOpMap()
var escapeRgx = regexp.MustCompile("[\\\\_%]")

type operator struct {
	name, sign string
	process    func(value reflect.Value) (string, []any)
	isValid    func(value reflect.Value) bool
}

func ok(reflect.Value) bool {
	return true
}

func isNotBlank(value reflect.Value) bool {
	s := value.String()
	return strings.TrimSpace(s) != ""
}

func ReadValueToArray(value reflect.Value) (string, []any) {
	return "?", []any{ReadValue(value)}
}

func ReadValueForIn(value reflect.Value) []any {
	arg := reflect.Indirect(value)
	args := make([]any, 0, arg.Len())
	for i := 0; i < arg.Len(); i++ {
		args = append(args, arg.Index(i).Interface())
	}
	return args
}

func checkValueForIn(value reflect.Value) bool {
	args := ReadValueForIn(value)
	return len(args) > 0
}

func BuildArgsForIn(value reflect.Value) (string, []any) {
	args := ReadValueForIn(value)
	argsLen := len(args)
	ph := bytes.NewBuffer(make([]byte, 0, 3*argsLen))
	ph.WriteString("(")
	for i := 0; i < argsLen; i++ {
		ph.WriteString("?")
		if i < argsLen-1 {
			ph.WriteString(", ")
		}
	}
	ph.WriteString(")")
	return ph.String(), args
}

func ReadLikeValue(value reflect.Value) string {
	s := value.String()
	return escapeRgx.ReplaceAllString(s, "\\$0")
}

func CreateOpMap() map[string]operator {
	const Like = " LIKE "
	const NotLike = " NOT LIKE "
	opMap := make(map[string]operator)
	opMap["Gt"] = operator{"Gt", " > ", ReadValueToArray, ok}
	opMap["Ge"] = operator{"Ge", " >= ", ReadValueToArray, ok}
	opMap["Lt"] = operator{"Lt", " < ", ReadValueToArray, ok}
	opMap["Le"] = operator{"Le", " <= ", ReadValueToArray, ok}
	opMap["Ne"] = operator{"Ne", " <> ", ReadValueToArray, ok}
	opMap["Eq"] = operator{"Eq", " = ", ReadValueToArray, ok}
	opMap["Null"] = operator{"Null", "", func(rv reflect.Value) (string, []any) {
		if rv.Bool() == false {
			return " IS NOT NULL", []any{}
		}
		return " IS NULL", []any{}
	}, ok}
	opMap["In"] = operator{"In", " IN ", BuildArgsForIn, checkValueForIn}
	opMap["NotIn"] = operator{"NotIn", " NOT IN ", BuildArgsForIn, checkValueForIn}
	opMap["Like"] = operator{"Like", Like, func(value reflect.Value) (string, []any) {
		s := value.String()
		ph := resolvePlaceHolder(s)
		return ph, []any{s}
	}, isNotBlank}
	opMap["NotLike"] = operator{"NotLike", NotLike, func(value reflect.Value) (string, []any) {
		s := value.String()
		ph := resolvePlaceHolder(s)
		return ph, []any{s}
	}, isNotBlank}
	opMap["Contain"] = operator{"Contain", Like, func(value reflect.Value) (string, []any) {
		escape := ReadLikeValue(value)
		ph := resolvePlaceHolder(escape)
		return ph, []any{"%" + escape + "%"}
	}, isNotBlank}
	opMap["NotContain"] = operator{"NotContain", NotLike, func(value reflect.Value) (string, []any) {
		escape := ReadLikeValue(value)
		ph := resolvePlaceHolder(escape)
		return ph, []any{"%" + escape + "%"}
	}, isNotBlank}
	opMap["Start"] = operator{"Start", Like, func(value reflect.Value) (string, []any) {
		escape := ReadLikeValue(value)
		ph := resolvePlaceHolder(escape)
		return ph, []any{escape + "%"}
	}, isNotBlank}
	opMap["NotStart"] = operator{"NotStart", NotLike, func(value reflect.Value) (string, []any) {
		escape := ReadLikeValue(value)
		ph := resolvePlaceHolder(escape)
		return ph, []any{escape + "%"}
	}, isNotBlank}
	opMap["End"] = operator{"End", Like, func(value reflect.Value) (string, []any) {
		escape := ReadLikeValue(value)
		ph := resolvePlaceHolder(escape)
		return ph, []any{"%" + escape}
	}, isNotBlank}
	opMap["NotEnd"] = operator{"NotEnd", NotLike, func(value reflect.Value) (string, []any) {
		escape := ReadLikeValue(value)
		ph := resolvePlaceHolder(escape)
		return ph, []any{"%" + escape}
	}, isNotBlank}
	opMap["Rx"] = operator{"Rx", " REGEXP ", ReadValueToArray, isNotBlank}
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
	if match := SuffixRgx.FindStringSubmatch(fieldName); len(match) > 0 {
		op := opMap[match[1]]
		column := strings.TrimSuffix(fieldName, match[1])
		column = ConvertToColumnCase(column)
		return fpSuffix{column, op}
	}
	return fpSuffix{ConvertToColumnCase(fieldName), opMap["Eq"]}
}

func (fp fpSuffix) Process(value reflect.Value) (string, []any) {
	if !fp.op.isValid(value) {
		return "", []any{}
	}
	placeholder, args := fp.op.process(value)
	return fp.col + fp.op.sign + placeholder, args
}
