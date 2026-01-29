/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2026, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package core

import (
	"io"
	"reflect"

	log "github.com/sirupsen/logrus"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func P[T any](t T) *T { return &t }

func ReadValue(value reflect.Value) any {
	typeStr := value.Type().String()
	log.Debug("Read value for type: ", typeStr)
	if value.Kind() == reflect.Ptr && !value.Elem().IsValid() {
		return nil
	}
	return reflect.Indirect(value).Interface()
}

func ConvertToColumnCase(fieldName string) string {
	return ToSnakeCase(fieldName)
}

func ToSnakeCase(fieldName string) string {
	col := make([]rune, 0, 2*len(fieldName))
	for i, letter := range fieldName {
		if letter >= 'A' && letter <= 'Z' {
			if i > 0 {
				col = append(col, '_')
			}
			col = append(col, letter|0x20)
		} else {
			col = append(col, letter)
		}
	}
	return string(col)
}

func ReadError(err error) *string {
	if err == nil {
		return nil
	}
	return P(err.Error())
}

func NoError(err error) bool {
	if err != nil {
		log.Error("Error occurred! ", err)
	}
	return err == nil
}

func HasError(err error) bool {
	return !NoError(err)
}

func Close(db io.Closer) {
	NoError(db.Close())
}

var capitalizer = cases.Title(language.English, cases.NoLower)

func Capitalize(str string) string {
	return capitalizer.String(str)
}

func Ternary[T any](test bool, r1, r2 T) T {
	if test {
		return r1
	}
	return r2
}
