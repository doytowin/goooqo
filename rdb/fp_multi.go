package rdb

import (
	"reflect"
	"strings"
)

// build multiple conditions mapped by a struct and connect them
type fpMultiConditions struct {
	connect func([]string) string
}

var fpForOr = &fpMultiConditions{connect: func(conditions []string) string {
	return "(" + strings.Join(conditions, " OR ") + ")"
}}

var fpForAnd = &fpMultiConditions{connect: func(conditions []string) string {
	return strings.Join(conditions, " AND ")
}}

func (fp *fpMultiConditions) Process(value reflect.Value) (string, []any) {
	conditions, args := buildConditions(value.Elem().Interface())
	return fp.connect(conditions), args
}
