/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2026, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package rdb

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	. "github.com/doytowin/goooqo/core"
)

var whereId = " WHERE id = ?"
var emMap = make(map[string]*metadata)

type metadata struct {
	TableName string
}

type EntityMetadata[E Entity] struct {
	metadata
	columnMetas     []FieldMetadata
	relationMetas   []FieldMetadata
	ColStr          string
	fieldsWithoutId []string
	createStr       string
	placeholders    string
	updateStr       string
	Type            reflect.Type
}

func RegisterEntity(entityName string, tableName string) {
	emMap[entityName] = &metadata{TableName: tableName}
}

func (em *EntityMetadata[E]) buildArgs(entity E) []any {
	args := make([]any, len(em.fieldsWithoutId))
	rv := reflect.ValueOf(entity)
	for i, col := range em.fieldsWithoutId {
		value := rv.FieldByName(col)
		args[i] = ReadValue(value)
	}
	return args
}

func (em *EntityMetadata[E]) buildSelect(query Query) (string, []any) {
	whereClause, args := BuildWhereClause(query)
	s := "SELECT " + em.ColStr + " FROM " + em.TableName + whereClause
	s += BuildSortClause(query.GetSort())
	if query.NeedPaging() {
		s = Dialect.BuildPageClause(s, query.CalcOffset(), query.GetPageSize())
	}
	return s, args
}

func (em *EntityMetadata[E]) buildSelectById() string {
	return "SELECT " + em.ColStr + " FROM " + em.TableName + whereId
}

func (em *EntityMetadata[E]) buildCount(query Query) (string, []any) {
	whereClause, args := BuildWhereClause(query)
	sqlStr := "SELECT count(0) FROM " + em.TableName + whereClause
	return sqlStr, args
}

func (em *EntityMetadata[E]) buildDeleteById() string {
	return "DELETE FROM " + em.TableName + whereId
}

func (em *EntityMetadata[E]) buildDelete(query any) (string, []any) {
	whereClause, args := BuildWhereClause(query)
	sqlStr := "DELETE FROM " + em.TableName + whereClause
	return sqlStr, args
}

func (em *EntityMetadata[E]) buildCreate(entity E) (string, []any) {
	return em.createStr, em.buildArgs(entity)
}

func (em *EntityMetadata[E]) buildCreateMulti(entities []E) (string, []any) {
	args := make([]any, 0, len(entities)*len(em.fieldsWithoutId))
	for _, entity := range entities {
		args = append(args, em.buildArgs(entity)...)
	}
	createStr := em.createStr + strings.Repeat(", "+em.placeholders, len(entities)-1)
	return createStr, args
}

func (em *EntityMetadata[E]) buildUpdate(entity E) (string, []any) {
	args := em.buildArgs(entity)
	args = append(args, entity.GetId())
	return em.updateStr, args
}

func (em *EntityMetadata[E]) buildPatch(entity Entity, extra int) (string, []any) {
	rv := reflect.ValueOf(entity)
	var patchFields = em.fieldsWithoutId

	if structType := rv.Type(); structType != em.Type {
		for i := 0; i < structType.NumField(); i++ {
			if structType.Field(i).Type.Kind() != reflect.Struct {
				fieldname := structType.Field(i).Name
				patchFields = append(em.fieldsWithoutId, fieldname)
			}
		}
	}

	args := make([]any, 0, len(patchFields)+extra)
	sqlStr := "UPDATE " + em.TableName + " SET "
	setClauses := make([]string, 0)

	for _, col := range patchFields {
		value := rv.FieldByName(col)
		v := ReadValue(value)
		if v != nil {
			setClauses = append(setClauses, resolveSetClause(col))
			args = append(args, v)
		}
	}
	return sqlStr + strings.Join(setClauses, ", "), args
}

func resolveSetClause(fieldname string) string {
	if strings.HasSuffix(fieldname, "Ae") {
		column := ConvertToColumnCase(strings.TrimSuffix(fieldname, "Ae"))
		sign := " + "
		return column + " = " + column + sign + "?"
	}
	return ConvertToColumnCase(fieldname) + " = ?"
}

func (em *EntityMetadata[E]) buildPatchById(entity Entity) (string, []any) {
	sqlStr, args := em.buildPatch(entity, 1)
	sqlStr = sqlStr + whereId
	args = append(args, entity.GetId())
	return sqlStr, args
}

func (em *EntityMetadata[E]) buildPatchByQuery(entity E, query Query) (string, []any, error) {
	whereClause, argsQ := BuildWhereClause(query)
	patchClause, argsE := em.buildPatch(entity, len(argsQ))

	if strings.HasSuffix(patchClause, "SET ") {
		return "", nil, errors.New("at least one field should be updated")
	}

	args := append(argsE, argsQ...)
	sqlStr := patchClause + whereClause

	return sqlStr, args, nil
}

func FormatTableByEntity(entity any) string {
	if rdbEntity, ok := entity.(RdbEntity); ok {
		return rdbEntity.GetTableName()
	}
	name := reflect.ValueOf(entity).Type().Name()
	name = strings.ToLower(strings.TrimSuffix(name, "Entity"))
	return fmt.Sprintf(Config.TableFormat, name)
}

func buildEntityMetadata[E Entity]() EntityMetadata[E] {
	entity := *new(E)
	entityType := reflect.TypeOf(entity)
	fieldMetas := BuildFieldMetas(entityType)

	columnMetas := make([]FieldMetadata, 0, len(fieldMetas))
	relationMetas := make([]FieldMetadata, 0, len(fieldMetas))

	for _, md := range fieldMetas {
		if md.EntityPath == nil {
			columnMetas = append(columnMetas, md)
		} else {
			relationMetas = append(relationMetas, md)
		}
	}

	columns := make([]string, len(columnMetas))
	columnsWithoutId := make([]string, 0, len(columnMetas))
	fieldsWithoutId := make([]string, 0, len(columnMetas))

	for i, md := range columnMetas {
		columns[i] = md.ColumnName
		if !md.IsId {
			fieldsWithoutId = append(fieldsWithoutId, md.Field.Name)
			columnsWithoutId = append(columnsWithoutId, md.ColumnName)
		}
	}

	tableName := FormatTableByEntity(entity)

	placeholders := "(?" + strings.Repeat(", ?", len(columnsWithoutId)-1) + ")"
	createStr := "INSERT INTO " + tableName +
		" (" + strings.Join(columnsWithoutId, ", ") + ") " +
		"VALUES " + placeholders

	set := make([]string, len(columnsWithoutId))
	for i, col := range columnsWithoutId {
		set[i] = col + " = ?"
	}
	updateStr := "UPDATE " + tableName + " SET " + strings.Join(set, ", ") + whereId

	RegisterEntity(entityType.Name(), tableName)
	return EntityMetadata[E]{
		metadata:        *emMap[entityType.Name()],
		columnMetas:     columnMetas,
		relationMetas:   relationMetas,
		ColStr:          strings.Join(columns, ", "),
		fieldsWithoutId: fieldsWithoutId,
		createStr:       createStr,
		placeholders:    placeholders,
		updateStr:       updateStr,
		Type:            reflect.TypeOf(*new(E)),
	}
}
