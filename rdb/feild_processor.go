package rdb

import (
	log "github.com/sirupsen/logrus"
	"reflect"
	"strings"
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

		fpKey := buildFpKey(queryType, field)
		if strings.HasSuffix(field.Name, "Or") {
			typeName := field.Type.Elem().String()
			if strings.Contains(typeName, "[]") {
				fpMap[fpKey] = buildFpOrForBasicArray(field.Name)
			} else {
				fpMap[fpKey] = buildFpOr()
			}
		} else if strings.HasSuffix(field.Name, "And") {
			fpMap[fpKey] = buildFpAnd()
		} else if _, ok := field.Tag.Lookup("subquery"); ok {
			fpMap[fpKey] = buildFpSubquery(field)
		} else if _, ok := field.Tag.Lookup("condition"); ok {
			fpMap[fpKey] = buildFpCustom(field)
		} else if field.Type.Kind() == reflect.Ptr &&
			field.Type.Elem().Kind() == reflect.Struct {
			log.Info("[registerFpByType] field: ", field.Type.Elem().Name(), " ", field.Type.Elem().Kind())
			registerFpByType(field.Type.Elem())
		} else if field.Type.Name() != "PageQuery" {
			fpMap[fpKey] = &fpSuffix{field.Name}
		}
	}
}

func buildFpAnd() FieldProcessor {
	return &fpMulti{Connector: " AND "}
}
