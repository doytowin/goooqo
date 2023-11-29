package goquery

import (
	"database/sql"
	fp "github.com/doytowin/goquery/field"
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

type connection interface {
	Prepare(query string) (*sql.Stmt, error)
}

func noError(err error) bool {
	if err == nil {
		return true
	}
	log.Error("Error occurred! ", err)
	return false
}

func buildEntityMetadata[E comparable](entity any) EntityMetadata[E] {
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

func (em *EntityMetadata[E]) buildSelect(query GoQuery) (string, []any) {
	whereClause, args := fp.BuildWhereClause(query)
	s := "SELECT " + em.ColStr + " FROM " + em.TableName + whereClause
	log.Debug("SQL: " + s)
	pageQuery := query.GetPageQuery()
	if pageQuery.needPaging() {
		s += pageQuery.buildPageClause()
	}
	return s, args
}

func (em *EntityMetadata[E]) buildSelectById() string {
	return "SELECT " + em.ColStr + " FROM " + em.TableName + " WHERE id = ?"
}

func (em *EntityMetadata[E]) Get(conn connection, id any) (E, error) {
	sqlStr := em.buildSelectById()
	rows, err := em.doQuery(conn, sqlStr, []any{id})
	if len(rows) == 1 {
		return rows[0], err
	}
	return em.zero, err
}

func (em *EntityMetadata[E]) Query(conn connection, query GoQuery) ([]E, error) {
	sqlStr, args := em.buildSelect(query)
	return em.doQuery(conn, sqlStr, args)
}

func (em *EntityMetadata[E]) doQuery(conn connection, sqlStr string, args []any) ([]E, error) {
	var result []E

	length := len(em.Columns)
	pointers := make([]any, length)
	entity := reflect.New(em.Type).Elem().Interface().(E)
	elem := reflect.ValueOf(&entity).Elem()
	for i := range pointers {
		pointers[i] = elem.Field(i).Addr().Interface()
	}

	stmt, err := conn.Prepare(sqlStr)
	if noError(err) {
		rows, err := stmt.Query(args...)
		for noError(err) && rows.Next() {
			err := rows.Scan(pointers...)
			if noError(err) {
				result = append(result, entity)
			}
		}
		_ = rows.Close()
		_ = stmt.Close()
	}

	return result, err
}

func (em *EntityMetadata[E]) buildCount(query GoQuery) (string, []any) {
	whereClause, args := fp.BuildWhereClause(query)
	s := "SELECT count(0) FROM " + em.TableName + whereClause

	log.Debug("SQL: ", s)
	pageQuery := query.GetPageQuery()
	if pageQuery.needPaging() {
		s += pageQuery.buildPageClause()
	}
	return s, args
}

func (em *EntityMetadata[E]) Count(conn connection, query GoQuery) (int, error) {
	cnt := 0
	sqlStr, args := em.buildCount(query)
	stmt, err := conn.Prepare(sqlStr)
	if noError(err) {
		row := stmt.QueryRow(args...)
		err = row.Scan(&cnt)
		_ = stmt.Close()
	}
	return cnt, err
}

func (em *EntityMetadata[E]) Page(conn connection, query GoQuery) (PageList[E], error) {
	var count int
	data, err := em.Query(conn, query)
	if noError(err) {
		count, err = em.Count(conn, query)
	}
	return PageList[E]{data, count}, err
}

func (em *EntityMetadata[E]) IsZero(entity E) bool {
	return em.zero == entity
}

func (em *EntityMetadata[E]) buildDeleteById() string {
	return "DELETE FROM " + em.TableName + " WHERE id = ?"
}

func (em *EntityMetadata[E]) DeleteById(conn connection, id any) (int64, error) {
	sqlStr := em.buildDeleteById()
	return em.doUpdate(conn, sqlStr, []any{id})
}

func (em *EntityMetadata[E]) buildDelete(query any) (string, []any) {
	whereClause, args := fp.BuildWhereClause(query)
	s := "DELETE FROM " + em.TableName + whereClause
	log.Debug("SQL: " + s)
	return s, args
}

func (em *EntityMetadata[E]) Delete(conn connection, query any) (int64, error) {
	sqlStr, args := em.buildDelete(query)
	return em.doUpdate(conn, sqlStr, args)
}

func (em *EntityMetadata[E]) doUpdate(conn connection, sqlStr string, args []any) (int64, error) {
	stmt, err := conn.Prepare(sqlStr)
	if noError(err) {
		defer func() {
			noError(stmt.Close())
		}()

		result, err := stmt.Exec(args...)
		if noError(err) {
			return result.RowsAffected()
		}
	}
	return 0, err
}
