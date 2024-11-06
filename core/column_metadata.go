/*
 * The Clear BSD License
 *
 * Copyright (c) 2024, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package core

import (
	"reflect"
	"strings"
)

type FieldMetadata struct {
	Field      reflect.StructField
	IsId       bool
	ColumnName string
	EntityPath *EntityPath
}

func BuildFieldMetas(structType reflect.Type) []FieldMetadata {
	fieldMetas := make([]FieldMetadata, 0, structType.NumField())
	for i := 0; i < structType.NumField(); i++ {
		fieldMetas = append(fieldMetas, buildFieldMetadata(structType.Field(i))...)
	}
	return fieldMetas
}

func buildFieldMetadata(field reflect.StructField) []FieldMetadata {
	if field.Type.Kind() == reflect.Struct {
		return BuildFieldMetas(field.Type)
	}
	cm := FieldMetadata{
		Field:      field,
		IsId:       field.Name == "Id",
		ColumnName: ConvertToColumnCase(field.Name),
	}
	if _, ok := field.Tag.Lookup("entitypath"); ok {
		cm.EntityPath = BuildEntityPath(field)
	}
	return []FieldMetadata{cm}
}

type EntityPath struct {
	Path         []string
	JoinTables   []string
	JoinIds      []string
	TargetTable  string
	LocalField   string
	ForeignField string
}

func BuildEntityPath(field reflect.StructField) *EntityPath {
	path := strings.Split(field.Tag.Get("entitypath"), ",")
	l := len(path)
	joinIds := make([]string, l)
	for i, domain := range path {
		joinIds[i] = FormatJoinId(domain)
	}
	joinTables := make([]string, l-1)
	for i := 0; i < l-1; i++ {
		joinTables[i] = FormatJoinTable(path[i], path[i+1])
	}
	targetTable := FormatTable(path[l-1])
	localFieldColumn := ConvertToColumnCase(field.Tag.Get("localField"))
	if localFieldColumn == "" {
		localFieldColumn = "id"
	}
	foreignFieldColumn := ConvertToColumnCase(field.Tag.Get("foreignField"))
	if foreignFieldColumn == "" {
		foreignFieldColumn = "id"
	}
	return &EntityPath{path, joinTables, joinIds, targetTable, localFieldColumn, foreignFieldColumn}
}
