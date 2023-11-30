package goquery

import (
	"github.com/doytowin/goquery/field"
	. "github.com/doytowin/goquery/util"
	"github.com/sirupsen/logrus"
	"reflect"
)

var whereId = " WHERE id = ?"

func readId(entity any) any {
	rv := reflect.ValueOf(entity)
	value := rv.FieldByName("Id")
	readValue := ReadValue(value)
	return readValue
}

func (em *EntityMetadata[E]) buildArgs(entity E) []any {
	var args []any

	rv := reflect.ValueOf(entity)
	for _, col := range em.fieldsWithoutId {
		value := rv.FieldByName(col)
		args = append(args, ReadValue(value))
	}
	return args
}

func (em *EntityMetadata[E]) buildSelect(query GoQuery) (string, []any) {
	whereClause, args := field.BuildWhereClause(query)
	s := "SELECT " + em.ColStr + " FROM " + em.TableName + whereClause
	logrus.Debug("SQL: " + s)
	pageQuery := query.GetPageQuery()
	if pageQuery.needPaging() {
		s += pageQuery.buildPageClause()
	}
	return s, args
}

func (em *EntityMetadata[E]) buildSelectById() string {
	return "SELECT " + em.ColStr + " FROM " + em.TableName + whereId
}

func (em *EntityMetadata[E]) buildCount(query GoQuery) (string, []any) {
	whereClause, args := field.BuildWhereClause(query)
	s := "SELECT count(0) FROM " + em.TableName + whereClause

	logrus.Debug("SQL: ", s)
	pageQuery := query.GetPageQuery()
	if pageQuery.needPaging() {
		s += pageQuery.buildPageClause()
	}
	return s, args
}

func (em *EntityMetadata[E]) buildDeleteById() string {
	return "DELETE FROM " + em.TableName + whereId
}

func (em *EntityMetadata[E]) buildDelete(query any) (string, []any) {
	whereClause, args := field.BuildWhereClause(query)
	s := "DELETE FROM " + em.TableName + whereClause
	logrus.Debug("SQL: " + s)
	return s, args
}

func (em *EntityMetadata[E]) buildCreate(entity E) (string, []any) {
	return em.createStr, em.buildArgs(entity)
}

func (em *EntityMetadata[E]) buildUpdate(entity E) (string, []any) {
	args := em.buildArgs(entity)
	args = append(args, readId(entity))
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
	sqlStr = sqlStr[0:len(sqlStr)-2] + whereId
	args = append(args, readId(entity))
	logrus.Info("UPDATE SQL: ", sqlStr)
	return sqlStr, args
}
