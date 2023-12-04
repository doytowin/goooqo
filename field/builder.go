package field

import (
	"github.com/doytowin/goquery/util"
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
			} else if strings.HasSuffix(fieldName, "In") {
				conditions, args = resolveIn(conditions, args, fieldName, reflect.Indirect(value))
			} else {
				conditions = append(conditions, Process(fieldName))
				if !strings.HasSuffix(fieldName, "Null") {
					args = append(args, util.ReadValue(value))
				}
			}
		}
	}
	return conditions, args
}

func resolveIn(conditions []string, args []any, fieldName string, arg reflect.Value) ([]string, []any) {
	condition := Process(fieldName)
	ph := "("
	for i := 0; i < arg.Len(); i++ {
		args = append(args, arg.Index(i).Int())
		ph += "?"
		if i < arg.Len()-1 {
			ph += ", "
		}
	}
	ph += ")"
	condition += ph
	conditions = append(conditions, condition)
	return conditions, args
}
