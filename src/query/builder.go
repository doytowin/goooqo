package query

import (
	suffix "../field"
	"reflect"
	"strings"
)

func BuildConditions(query interface{}) (string, []any) {
	refType := reflect.TypeOf(query)
	rv := reflect.ValueOf(query)
	cnt, argCnt := 0, 0
	conditions := make([]string, refType.NumField())
	args := make([]any, refType.NumField(), 2*refType.NumField())
	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)
		value := rv.FieldByName(field.Name)
		if IsValidValue(value) {
			conditions[cnt] = suffix.Process(field.Name)
			cnt++
			if value.Type().String() == "*int" {
				args[argCnt] = reflect.Indirect(value).Int()
				argCnt++
			}
		}
	}
	return strings.Join(conditions[0:cnt], " AND "), args[0:argCnt]
}

func IsValidValue(value reflect.Value) bool {
	return (value.Type().Name() == "bool" && value.Bool()) || !value.IsNil()
}
