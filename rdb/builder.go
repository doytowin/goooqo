/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2025, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package rdb

import (
	"reflect"
	"strings"
)

type QueryBuilder interface {
	BuildConditions() ([]string, []any)
}

func isValidValue(value reflect.Value) bool {
	return !value.IsNil()
}

func BuildWhereClause(query any) (string, []any) {
	return BuildConditions(query, " WHERE ", " AND ", "")
}

func BuildConditions(query any, prefix string, delimiter string, suffix string) (a string, args []any) {
	var conditions []string
	if qb, ok := query.(QueryBuilder); ok {
		conditions, args = qb.BuildConditions()
	} else {
		conditions, args = buildConditions(query)
	}
	if len(conditions) == 0 {
		return "", []any{}
	}
	return prefix + strings.Join(conditions, delimiter) + suffix, args
}

func buildConditions(query any) ([]string, []any) {
	rtype := reflect.TypeOf(query)
	rvalue := reflect.ValueOf(query)
	if rtype.Kind() == reflect.Pointer {
		rtype = rtype.Elem()
		rvalue = rvalue.Elem()
	}
	args := make([]any, 0, rtype.NumField())
	conditions := make([]string, 0, rtype.NumField())

	registerFpByType(rtype)
	for i := 0; i < rtype.NumField(); i++ {
		field := rtype.Field(i)
		fpKey := buildFpKey(rtype, field)
		processor := fpMap[fpKey]
		if processor != nil {
			value := rvalue.FieldByName(field.Name)
			if isValidValue(value) {
				condition, arr := processor.Process(value.Elem())
				if condition != "" {
					conditions = append(conditions, condition)
					args = append(args, arr...)
				}
			}
		}
	}
	return conditions, args
}
