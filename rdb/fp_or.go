package rdb

import (
	"reflect"
	"strings"
)

type fpOr struct {
}

func buildFpOr() FieldProcessor {
	return &fpOr{}
}

func (f *fpOr) Process(value reflect.Value) (string, []any) {
	return ProcessOr(value.Elem().Interface())
}

func ProcessOr(or any) (string, []any) {
	conditions, args := buildConditions(or)
	return "(" + strings.Join(conditions, " OR ") + ")", args
}
