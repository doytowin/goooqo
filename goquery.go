package goquery

import (
	"database/sql"
	fp "github.com/doytowin/doyto-query-go-sql/field"
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

type Database interface {
	Prepare(query string) (*sql.Stmt, error)
}

func IntPtr(o int) *int {
	return &o
}

func BuildEntityMetadata[E comparable](entity interface{}) EntityMetadata[E] {
	refType := reflect.TypeOf(entity)
	columns := make([]string, refType.NumField())
	for i := 0; i < refType.NumField(); i++ {
		field := refType.Field(i)
		columns[i] = fp.UnCapitalize(field.Name)
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
	whereClause, args := fp.BuildWhereClause(query)
	s := "SELECT " + em.ColStr + " FROM " + em.TableName + whereClause
	log.Info("SQL: " + s)
	return s, args
}

func (em *EntityMetadata[E]) buildSelectById() string {
	return "SELECT " + em.ColStr + " FROM " + em.TableName + " WHERE id = ?"
}

func (em *EntityMetadata[E]) Get(db Database, id interface{}) (E, error) {
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

func (em *EntityMetadata[E]) Query(db Database, query interface{}) ([]E, error) {
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

func (em *EntityMetadata[E]) buildDeleteById() string {
	return "DELETE FROM " + em.TableName + " WHERE id = ?"
}

func (em *EntityMetadata[E]) DeleteById(db Database, id interface{}) (int64, error) {
	sqlStr := em.buildDeleteById()
	stmt, _ := db.Prepare(sqlStr)
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	result, err := stmt.Exec(id)
	if err != nil {
		return 0, err
	}
	cnt, err := result.RowsAffected()
	return cnt, err
}

func (em *EntityMetadata[E]) buildDelete(query interface{}) (string, []any) {
	whereClause, args := fp.BuildWhereClause(query)
	s := "DELETE FROM " + em.TableName + whereClause
	log.Info("SQL: " + s)
	return s, args
}

func (em *EntityMetadata[E]) Delete(db Database, query interface{}) (int64, error) {
	sqlStr, args := em.buildDelete(query)
	stmt, _ := db.Prepare(sqlStr)
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	result, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	cnt, err := result.RowsAffected()
	return cnt, err
}
