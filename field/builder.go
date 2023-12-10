package field

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
	var (
		args       []any
		conditions []string
		rtype      = reflect.TypeOf(query)
		rvalue     = reflect.ValueOf(query)
	)
	if rtype.Kind() == reflect.Pointer {
		rtype = rtype.Elem()
		rvalue = rvalue.Elem()
	}

	for i := 0; i < rtype.NumField(); i++ {
		field := rtype.Field(i)
		fieldName := field.Name
		value := rvalue.FieldByName(fieldName)
		if isValidValue(value) {
			if strings.HasSuffix(fieldName, "Or") {
				condition, arr := ProcessOr(value.Elem().Interface())
				conditions = append(conditions, condition)
				args = append(args, arr...)
			} else {
				condition, arg := Process(fieldName, value)
				conditions = append(conditions, condition)
				args = append(args, arg...)
			}
		}
	}
	return conditions, args
}
