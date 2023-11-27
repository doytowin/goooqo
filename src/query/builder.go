package query

import (
	suffix "../field"
	"reflect"
	"strings"
)

type EntityMetadata struct {
	TableName string
	Columns   []string
	ColStr    string
}

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
	if value.Type().Name() == "bool" {
		return value.Bool()
	} else {
		return !value.IsNil()
	}
}

func BuildEntityMetadata(entity interface{}) EntityMetadata {
	refType := reflect.TypeOf(entity)
	conditions := make([]string, refType.NumField())
	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)
		conditions[i] = field.Name
	}
	return EntityMetadata{
		TableName: strings.TrimSuffix(refType.Name(), "Entity"),
		Columns:   conditions,
		ColStr:    strings.Join(conditions, ", "),
	}
}

func (em *EntityMetadata) BuildSelect(query interface{}) (string, []any) {
	conditions, args := BuildConditions(query)
	return "SELECT " + em.ColStr +
		" FROM " + em.TableName +
		" WHERE " + conditions, args
}
