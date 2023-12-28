package rdb

import (
	log "github.com/sirupsen/logrus"
	"reflect"
)

var fpMap = make(map[string]FieldProcessor)
var fpTypeMap = make(map[reflect.Type]bool)

type FieldProcessor interface {
	Process(value reflect.Value) (string, []any)
}

func buildFpKey(queryType reflect.Type, field reflect.StructField) string {
	return queryType.PkgPath() + ":" + queryType.Name() + ":" + field.Name
}

func registerFpByType(queryType reflect.Type) {
	if fpTypeMap[queryType] == true {
		return
	}
	fpTypeMap[queryType] = true

	for i := 0; i < queryType.NumField(); i++ {
		field := queryType.Field(i)
		if _, ok := field.Tag.Lookup("subquery"); ok {
			fpMap[buildFpKey(queryType, field)] = buildFpSubquery(field)
		} else if field.Type.Kind() == reflect.Ptr &&
			field.Type.Elem().Kind() == reflect.Struct {
			log.Info("[registerFpByType] field: ", field.Type.Elem().Name(), " ", field.Type.Elem().Kind())
			registerFpByType(field.Type.Elem())
		} else if field.Type.Name() != "PageQuery" {
			fpMap[buildFpKey(queryType, field)] = &fpSuffix{field}
		}
	}
}
