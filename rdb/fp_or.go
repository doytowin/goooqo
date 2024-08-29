package rdb

import (
	"reflect"
	"strings"
)

// build multiple conditions and connect them by the Connector
type fpMulti struct {
	Connector string
}

func buildFpOr() FieldProcessor {
	return &fpMulti{Connector: " OR "}
}

func (f *fpMulti) Process(value reflect.Value) (string, []any) {
	conditions, args := buildConditions(value.Elem().Interface())
	condition := strings.Join(conditions, f.Connector)
	if f.Connector == " OR " {
		condition = "(" + condition + ")"
	}
	return condition, args
}

type fpOr struct {
	fpSuffix fpSuffix
}

func buildFpOrForBasicArray(fieldName string) FieldProcessor {
	return &fpOr{fpSuffix: fpSuffix{fieldName: strings.TrimSuffix(fieldName, "Or")}}
}

func (fp *fpOr) Process(value reflect.Value) (condition string, args []any) {
	value = value.Elem()
	conditions := make([]string, value.Len())
	var arr []any
	for i := 0; i < value.Len(); i++ {
		conditions[i], arr = fp.fpSuffix.Process(value.Index(i))
		args = append(args, arr...)
	}
	condition = strings.Join(conditions, " OR ")
	return "(" + condition + ")", args
}
