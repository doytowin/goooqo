package field

import (
	"reflect"
	"strings"
)

func isValidValue(value reflect.Value) bool {
	if value.Type().Name() == "bool" {
		return value.Bool()
	} else if value.Type().Name() == "string" {
		return value.String() != ""
	} else {
		return !value.IsNil()
	}
}

func BuildWhereClause(query interface{}) (string, []any) {
	conditions, args := buildConditions(query)
	if len(conditions) == 0 {
		return "", []any{}
	}
	return " WHERE " + strings.Join(conditions, " AND "), args
}

func buildConditions(query interface{}) ([]string, []any) {
	refType := reflect.TypeOf(query)
	rv := reflect.ValueOf(query)
	cnt, argCnt := 0, 0
	conditions := make([]string, refType.NumField())
	args := make([]any, refType.NumField(), 2*refType.NumField())
	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)
		value := rv.FieldByName(field.Name)
		if isValidValue(value) {
			conditions[cnt] = Process(field.Name)
			cnt++
			if value.Type().String() == "*int" {
				args[argCnt] = reflect.Indirect(value).Int()
				argCnt++
			}
		}
	}
	return conditions[0:cnt], args[0:argCnt]
}
