package core

import (
	"reflect"
)

type ColumnMetadata struct {
	Field      reflect.StructField
	IsId       bool
	ColumnName string
}

func BuildColumnMetas(structType reflect.Type) (columnMetas []ColumnMetadata) {
	for i := 0; i < structType.NumField(); i++ {
		columnMetas = append(columnMetas, buildColumnMetadata(structType.Field(i))...)
	}
	return
}

func buildColumnMetadata(field reflect.StructField) []ColumnMetadata {
	if field.Type.Kind() == reflect.Struct {
		return BuildColumnMetas(field.Type)
	}
	return []ColumnMetadata{{
		Field:      field,
		IsId:       field.Name == "Id",
		ColumnName: ConvertToColumnCase(field.Name),
	}}
}
