/*
 * The Clear BSD License
 *
 * Copyright (c) 2024, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package web

import (
	"github.com/doytowin/goooqo/core"
	log "github.com/sirupsen/logrus"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

var converterMap = map[reflect.Type]func(v []string) (any, error){}

func RegisterConverter(typeName reflect.Type, converter func(v []string) (any, error)) {
	converterMap[typeName] = converter
}

func init() {
	RegisterConverter(reflect.TypeOf(true), func(v []string) (any, error) {
		return strings.EqualFold(v[0], "TRue"), nil
	})
	RegisterConverter(reflect.PointerTo(reflect.TypeOf(true)), func(v []string) (any, error) {
		ok := strings.EqualFold(v[0], "TRue")
		return &ok, nil
	})

	RegisterConverter(reflect.TypeOf(0), func(v []string) (any, error) {
		return strconv.Atoi(v[0])
	})

	RegisterConverter(reflect.PointerTo(reflect.TypeOf(0)), func(v []string) (any, error) {
		v0, err := strconv.Atoi(v[0])
		return &v0, err
	})

	RegisterConverter(reflect.PointerTo(reflect.TypeOf(0.1)), func(v []string) (any, error) {
		v0, err := strconv.ParseFloat(v[0], 64)
		return &v0, err
	})

	RegisterConverter(reflect.PointerTo(reflect.TypeOf([]int{0})), func(params []string) (any, error) {
		if len(params) == 1 {
			params = strings.Split(params[0], ",")
		}
		v := make([]int, 0, len(params))
		for _, s := range params {
			num, err := strconv.Atoi(s)
			if core.NoError(err) {
				v = append(v, num)
			}
		}
		return &v, nil
	})

	RegisterConverter(reflect.PointerTo(reflect.TypeOf("")), func(v []string) (any, error) {
		return &v[0], nil
	})
}

func ResolveQuery(queryMap url.Values, query any) {
	elem := reflect.ValueOf(query).Elem()
	for name, v := range queryMap {
		path := strings.Split(name, ".")
		field := resolveParam(elem, path[0])
		for i := 1; i < len(path); i++ {
			if !field.IsValid() {
				break
			}
			if field.IsNil() {
				fieldType := field.Type().Elem()
				field.Set(reflect.New(fieldType))
			}
			field = resolveParam(field.Elem(), path[i])
		}

		if field.IsValid() {
			convertAndSet(field, v)
		}
	}
}

func resolveParam(elem reflect.Value, fieldName string) reflect.Value {
	field := elem.FieldByName(fieldName)
	if !field.IsValid() {
		title := core.Capitalize(fieldName)
		field = elem.FieldByName(title)
	}
	return field
}

func convertAndSet(field reflect.Value, v []string) {
	log.Debug("field.Type: ", field.Type())
	fieldType := field.Type()
	v0, err := converterMap[fieldType](v)
	if core.NoError(err) || v0 != nil {
		field.Set(reflect.ValueOf(v0))
	}
}
