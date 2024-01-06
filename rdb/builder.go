package rdb

import (
	log "github.com/sirupsen/logrus"
	"reflect"
	"strings"
)

func isValidValue(value reflect.Value) bool {
	typeName := value.Type().Name()
	if typeName == "bool" {
		return value.Bool()
	} else if typeName == "string" {
		return value.String() != ""
	} else if typeName == "flag" {
		return value.IsValid()
	} else if typeName == "PageQuery" {
		return false
	} else {
		log.Debug("Type:", typeName)
		return !value.IsNil()
	}
}

func BuildWhereClause(query any) (string, []any) {
	conditions, args := buildConditions(query)
	if len(conditions) == 0 {
		return "", []any{}
	}
	return " WHERE " + strings.Join(conditions, " AND "), args
}

func buildConditions(query any) ([]string, []any) {
	rtype := reflect.TypeOf(query)
	rvalue := reflect.ValueOf(query)
	if rtype.Kind() == reflect.Pointer {
		rtype = rtype.Elem()
		rvalue = rvalue.Elem()
	}
	args := make([]any, 0, rtype.NumField())
	conditions := make([]string, 0, rtype.NumField())

	registerFpByType(rtype)
	for i := 0; i < rtype.NumField(); i++ {
		field := rtype.Field(i)
		fieldName := field.Name
		value := rvalue.FieldByName(fieldName)
		if isValidValue(value) {
			processor := fpMap[buildFpKey(rtype, field)]
			condition, arr := processor.Process(value)
			conditions = append(conditions, condition)
			args = append(args, arr...)
		}
	}
	return conditions, args
}
