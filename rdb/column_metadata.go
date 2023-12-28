package rdb

import (
	"github.com/doytowin/goooqo/core"
	"reflect"
)

type columnMetadata struct {
	field      reflect.StructField
	isId       bool
	columnName string
}

func buildColumnMetas(structType reflect.Type) []columnMetadata {
	var columnMetas []columnMetadata
	for i := 0; i < structType.NumField(); i++ {
		columnMetas = append(columnMetas, buildColumnMetadata(structType.Field(i))...)
	}
	return columnMetas
}

func buildColumnMetadata(field reflect.StructField) []columnMetadata {
	if field.Type.Kind() == reflect.Struct {
		return buildColumnMetas(field.Type)
	}
	return []columnMetadata{{
		field:      field,
		isId:       field.Name == "Id",
		columnName: core.ConvertToColumnCase(field.Name),
	}}
}
