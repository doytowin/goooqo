/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2025, DoytoWin, Inc.
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

var typeFmMap = make(map[reflect.Type][]FieldMetadata)

func BuildFieldMetas(structType reflect.Type) []FieldMetadata {
	fieldMetas := typeFmMap[structType]
	if fieldMetas == nil {
		fieldMetas = make([]FieldMetadata, 0, structType.NumField())
		for i := 0; i < structType.NumField(); i++ {
			fieldMetas = append(fieldMetas, buildFieldMetadata(structType.Field(i))...)
		}
		typeFmMap[structType] = fieldMetas
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
	Path       []string
	Base       Relation
	Relations  []Relation
	EntityType reflect.Type
}

type Relation struct {
	Fk1, Fk2, At string
}

func BuildEntityPath(field reflect.StructField) *EntityPath {
	entityPath := BuildEntityPathStr(field.Tag.Get("entitypath"))
	entityPath.EntityType = field.Type.Elem()
	return entityPath
}

func BuildEntityPathStr(aep string) *EntityPath {
	path := strings.Split(aep, ",")
	l := len(path)
	relations := make([]Relation, l-1)
	for i := 0; i < l-1; i++ {
		relations[i] = buildRelation(path[i], path[i+1])
	}
	targetTable := FormatTable(MapVirtualEntity(path[l-1]))
	base := Relation{"id", "id", targetTable}
	if strings.Contains(aep, "-") && l == 2 {
		base = relations[0]
		relations = []Relation{}
	}
	return &EntityPath{path, base, relations, nil}
}

// e1: left entity, e2: right entity
func buildRelation(e1 string, e2 string) Relation {
	if entity, fk, ok := strings.Cut(e1, "->"); ok {
		return Relation{"id", ConvertToColumnCase(fk), FormatTable(entity)}
	} else if fk, entity, ok := strings.Cut(e2, "<-"); ok {
		return Relation{ConvertToColumnCase(fk), "id", FormatTable(entity)}
	}
	if _, entity, ok := strings.Cut(e1, "<-"); ok {
		e1 = entity
	}
	e2, _, _ = strings.Cut(e2, "->")
	e1 = MapVirtualEntity(e1)
	return Relation{FormatJoinId(e1), FormatJoinId(e2), FormatJoinTable(e1, e2)}
}
