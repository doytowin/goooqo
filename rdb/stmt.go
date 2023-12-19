package rdb

import (
	. "github.com/doytowin/go-query/core"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strings"
)

var whereId = " WHERE id = ?"

type EntityMetadata[E any] struct {
	TableName       string
	ColStr          string
	fieldsWithoutId []string
	createStr       string
	placeholders    string
	updateStr       string
}

func readId(entity any) any {
	rv := reflect.ValueOf(entity)
	value := rv.FieldByName("Id")
	return ReadValue(value)
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

func (em *EntityMetadata[E]) buildSelect(query GoQuery) (string, []any) {
	whereClause, args := BuildWhereClause(query)
	s := "SELECT " + em.ColStr + " FROM " + em.TableName + whereClause
	if query.NeedPaging() {
		s += query.BuildPageClause()
	}
	return s, args
}

func (em *EntityMetadata[E]) buildSelectById() string {
	return "SELECT " + em.ColStr + " FROM " + em.TableName + whereId
}

func (em *EntityMetadata[E]) buildCount(query GoQuery) (string, []any) {
	whereClause, args := BuildWhereClause(query)
	sqlStr := "SELECT count(0) FROM " + em.TableName + whereClause

	log.Debug("SQL: ", sqlStr)
	log.Debug("ARG: ", args)
	return sqlStr, args
}

func (em *EntityMetadata[E]) buildDeleteById() string {
	return "DELETE FROM " + em.TableName + whereId
}

func (em *EntityMetadata[E]) buildDelete(query any) (string, []any) {
	whereClause, args := BuildWhereClause(query)
	sqlStr := "DELETE FROM " + em.TableName + whereClause
	log.Debug("SQL: ", sqlStr)
	log.Debug("ARG: ", args)
	return sqlStr, args
}

func (em *EntityMetadata[E]) buildCreate(entity E) (string, []any) {
	args := em.buildArgs(entity)
	log.Debug("SQL: ", em.createStr)
	log.Debug("ARG: ", args)
	return em.createStr, args
}

func (em *EntityMetadata[E]) buildCreateMulti(entities []E) (string, []any) {
	var args []any
	for _, entity := range entities {
		args = append(args, em.buildArgs(entity)...)
	}
	createStr := em.createStr
	for i := 1; i < len(entities); i++ {
		createStr += ", " + em.placeholders
	}
	log.Debug("SQL: ", createStr)
	log.Debug("ARG: ", args)
	return createStr, args
}

func (em *EntityMetadata[E]) buildUpdate(entity E) (string, []any) {
	args := em.buildArgs(entity)
	args = append(args, readId(entity))
	log.Debug("SQL: ", em.updateStr)
	log.Debug("ARG: ", args)
	return em.updateStr, args
}

func (em *EntityMetadata[E]) buildPatch(entity E) (string, []any) {
	var args []any
	sqlStr := "UPDATE " + em.TableName + " SET "

	rv := reflect.ValueOf(entity)
	for _, col := range em.fieldsWithoutId {
		value := rv.FieldByName(col)
		v := ReadValue(value)
		if v != nil {
			sqlStr += UnCapitalize(col) + " = ?, "
			args = append(args, v)
		}
	}
	return sqlStr[0 : len(sqlStr)-2], args
}

func (em *EntityMetadata[E]) buildPatchById(entity E) (string, []any) {
	sqlStr, args := em.buildPatch(entity)
	sqlStr = sqlStr + whereId
	args = append(args, readId(entity))
	log.Debug("SQL: ", sqlStr)
	log.Debug("ARG: ", args)
	return sqlStr, args
}

func (em *EntityMetadata[E]) buildPatchByQuery(entity E, query GoQuery) ([]any, string) {
	patchClause, argsE := em.buildPatch(entity)
	whereClause, argsQ := BuildWhereClause(query)

	args := append(argsE, argsQ...)
	sqlStr := patchClause + whereClause

	log.Debug("SQL: ", sqlStr)
	log.Debug("ARG: ", args)
	return args, sqlStr
}

func buildEntityMetadata[E Entity](entity E) EntityMetadata[E] {
	refType := reflect.TypeOf(entity)
	columns := make([]string, refType.NumField())
	var columnsWithoutId []string
	var fieldsWithoutId []string
	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)
		columns[i] = UnCapitalize(field.Name)
		if field.Name != "Id" {
			fieldsWithoutId = append(fieldsWithoutId, field.Name)
			columnsWithoutId = append(columnsWithoutId, UnCapitalize(field.Name))
		}
	}

	var tableName = entity.GetTableName()

	placeholders := "(?"
	for i := 1; i < len(columnsWithoutId); i++ {
		placeholders += ", ?"
	}
	placeholders += ")"
	createStr := "INSERT INTO " + tableName +
		" (" + strings.Join(columnsWithoutId, ", ") + ") " +
		"VALUES " + placeholders

	set := make([]string, len(columnsWithoutId))
	for i, col := range columnsWithoutId {
		set[i] = col + " = ?"
	}
	updateStr := "UPDATE " + tableName + " SET " + strings.Join(set, ", ") + whereId

	return EntityMetadata[E]{
		TableName:       tableName,
		ColStr:          strings.Join(columns, ", "),
		fieldsWithoutId: fieldsWithoutId,
		createStr:       createStr,
		placeholders:    placeholders,
		updateStr:       updateStr,
	}
}
