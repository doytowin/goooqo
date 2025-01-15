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
	"github.com/doytowin/goooqo/core"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strings"
)

var fpMap = make(map[string]FieldProcessor)
var fpTypeMap = make(map[reflect.Type]bool)

type FieldProcessor interface {
	Process(value reflect.Value) (string, []any)
}

func buildFpKey(queryType reflect.Type, field reflect.StructField) string {
	return queryType.PkgPath() + ":" + queryType.Name() + ":" + field.Name
}

func registerFpByType(queryType reflect.Type) {
	if fpTypeMap[queryType] == true {
		return
	}
	fpTypeMap[queryType] = true
	typeQuery := reflect.TypeOf((*core.Query)(nil)).Elem()

	for i := 0; i < queryType.NumField(); i++ {
		field := queryType.Field(i)
		// ignore PageQuery
		if field.Anonymous && field.Type.Implements(typeQuery) {
			continue
		}

		fpKey := buildFpKey(queryType, field)
		if field.Type.Kind() != reflect.Ptr {
			log.Warn("Type not supported: ", field.Type)
		} else if strings.HasSuffix(field.Name, "Or") {
			if field.Type.Elem().Kind() == reflect.Slice {
				if field.Type.Elem().Elem().Kind() == reflect.Struct {
					fpMap[fpKey] = buildFpStructArrayByOr()
				} else {
					fpMap[fpKey] = buildFpBasicArrayByOr(field.Name)
				}
			} else {
				fpMap[fpKey] = fpForOr
			}
		} else if strings.HasSuffix(field.Name, "And") {
			fpMap[fpKey] = fpForAnd
		} else if field.Type.Implements(typeQuery) {
			buildForQuery(field, fpKey)
		} else if _, ok := field.Tag.Lookup("condition"); ok {
			fpMap[fpKey] = buildFpCustom(field)
		} else {
			fpMap[fpKey] = buildFpSuffix(field.Name)
		}
	}
}

func buildForQuery(field reflect.StructField, fpKey string) {
	if _, ok := field.Tag.Lookup("entitypath"); ok {
		fpMap[fpKey] = buildFpEntityPath(field)
	} else if subqueryTag, ok := field.Tag.Lookup("subquery"); ok {
		fpMap[fpKey] = BuildBySubqueryTag(subqueryTag, field.Name)
	} else if _, ok := field.Tag.Lookup("select"); ok {
		fpMap[fpKey] = BuildBySelectTag(field.Tag, field.Name)
	} else if match := subOfRgx.FindStringSubmatch(field.Name); len(match) > 0 {
		fpMap[fpKey] = BuildByFieldName(match)
	} else {
		log.Debug("Not mapped by field processor : ", fpKey)
	}
}
