package field

import (
	log "github.com/sirupsen/logrus"
	"reflect"
	"strings"
)

func isValidValue(value reflect.Value) bool {
	if value.Type().Name() == "bool" {
		return value.Bool()
	} else if value.Type().Name() == "string" {
		return value.String() != ""
	} else if value.Type().Name() == "flag" {
		return value.IsValid()
	} else if value.Type().Name() == "PageQuery" {
		return false
	} else {
		log.Info("Type:", value.Type().Name())
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
	cnt := 0
	conditions := make([]string, refType.NumField())
	var args []any

	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)
		value := rv.FieldByName(field.Name)
		if isValidValue(value) {
			if strings.HasSuffix(field.Name, "Or") {
				var arr []any
				conditions[cnt], arr = ProcessOr(value.Elem().Interface())
				cnt++
				args = append(args, arr...)
			} else {
				conditions[cnt] = Process(field.Name)
				cnt++
				typeStr := value.Type().String()
				switch typeStr {
				case "bool", "*bool":
					args = append(args, reflect.Indirect(value).Bool())
				case "*int":
					args = append(args, reflect.Indirect(value).Int())
				case "*string":
					args = append(args, reflect.Indirect(value).String())
				default:
					log.Warn("Type not support: ", typeStr)
				}
			}
		}
	}
	return conditions[0:cnt], args
}
