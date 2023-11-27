package goquery

import (
	"database/sql"
	suffix "github.com/doytowin/doyto-query-go-sql/field"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strings"
)

type EntityMetadata[E comparable] struct {
	Type      reflect.Type
	TableName string
	Columns   []string
	ColStr    string
	zero      E
}

func IntPtr(o int) *int {
	return &o
}

func buildWhereClause(query interface{}) (string, []any) {
	refType := reflect.TypeOf(query)
	rv := reflect.ValueOf(query)
	cnt, argCnt := 0, 0
	conditions := make([]string, refType.NumField())
	args := make([]any, refType.NumField(), 2*refType.NumField())
	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)
		value := rv.FieldByName(field.Name)
		if isValidValue(value) {
			conditions[cnt] = suffix.Process(field.Name)
			cnt++
			if value.Type().String() == "*int" {
				args[argCnt] = reflect.Indirect(value).Int()
				argCnt++
			}
		}
	}
	return strings.Join(conditions[0:cnt], " AND "), args[0:argCnt]
}

func isValidValue(value reflect.Value) bool {
	if value.Type().Name() == "bool" {
		return value.Bool()
	} else {
		return !value.IsNil()
	}
}

func BuildEntityMetadata[E comparable](entity interface{}) EntityMetadata[E] {
	refType := reflect.TypeOf(entity)
	columns := make([]string, refType.NumField())
	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)
		columns[i] = suffix.UnCapitalize(field.Name)
	}
	return EntityMetadata[E]{
		Type:      refType,
		TableName: strings.TrimSuffix(refType.Name(), "Entity"),
		Columns:   columns,
		ColStr:    strings.Join(columns, ", "),
		zero:      reflect.New(refType).Elem().Interface().(E),
	}
}

func (em *EntityMetadata[E]) buildSelect(query interface{}) (string, []any) {
	conditions, args := buildWhereClause(query)
	s := "SELECT " + em.ColStr + " FROM " + em.TableName
	if len(conditions) > 0 {
		s += " WHERE " + conditions
	}
	log.Info("SQL: " + s)
	return s, args
}

func (em *EntityMetadata[E]) buildSelectById() string {
	return "SELECT " + em.ColStr + " FROM " + em.TableName + " WHERE id = ?"
}

func (em *EntityMetadata[E]) Get(db *sql.DB, id interface{}) (E, error) {
	sqlStr := em.buildSelectById()
	stmt, err := db.Prepare(sqlStr)
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	rows, err := em.doQuery(stmt, []any{id})
	if len(rows) == 1 {
		return rows[0], err
	}
	return em.zero, err
}

func (em *EntityMetadata[E]) Query(db *sql.DB, query interface{}) ([]E, error) {
	sqlStr, args := em.buildSelect(query)
	stmt, _ := db.Prepare(sqlStr)
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	return em.doQuery(stmt, args)
}

func (em *EntityMetadata[E]) doQuery(stmt *sql.Stmt, args []any) ([]E, error) {
	var result []E

	length := len(em.Columns)
	pointers := make([]interface{}, length)
	entity := reflect.New(em.Type).Elem().Interface().(E)
	elem := reflect.ValueOf(&entity).Elem()
	for i := range pointers {
		pointers[i] = elem.Field(i).Addr().Interface()
	}

	rows, err := stmt.Query(args...)
	for rows.Next() && err == nil {
		err := rows.Scan(pointers...)
		if err == nil {
			result = append(result, entity)
		}
	}

	return result, err
}

func (em *EntityMetadata[E]) IsZero(entity E) bool {
	return em.zero == entity
}