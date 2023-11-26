package query

import (
	suffix "../field"
	"reflect"
	"strings"
)

func BuildConditions(query interface{}) string {
	refType := reflect.TypeOf(query)
	rv := reflect.ValueOf(query)
	cnt := 0
	conditions := make([]string, refType.NumField())
	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)
		if !rv.FieldByName(field.Name).IsNil() {
			conditions[cnt] = suffix.Process(field.Name)
			cnt++
		}
	}
	return strings.Join(conditions[0:cnt], " AND ")
}
