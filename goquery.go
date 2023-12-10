package goquery

import (
	"database/sql"
	. "github.com/doytowin/goquery/util"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strings"
)

type EntityMetadata[E comparable] struct {
	Type            reflect.Type
	TableName       string
	Columns         []string
	ColStr          string
	fieldsWithoutId []string
	placeholders    string
	createStr       string
	updateStr       string
	zero            E
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
	var tableName string
	v, ok := entity.(Entity)
	if ok {
		tableName = v.GetTableName()
	} else {
		tableName = strings.TrimSuffix(refType.Name(), "Entity")
	}

	placeholders := "(?"
	for i := 1; i < len(columnsWithoutId); i++ {
		placeholders += ", ?"
	}
	placeholders += ")"
	createStr := "INSERT INTO " + tableName +
		" (" + strings.Join(columnsWithoutId, ", ") + ") " +
		"VALUES " + placeholders
	log.Debug("CREATE SQL: ", createStr)

	set := make([]string, len(columnsWithoutId))
	for i, col := range columnsWithoutId {
		set[i] = col + " = ?"
	}
	updateStr := "UPDATE " + tableName + " SET " + strings.Join(set, ", ") + whereId
	log.Debug("UPDATE SQL: ", updateStr)

	return EntityMetadata[E]{
		Type:            refType,
		TableName:       tableName,
		Columns:         columns,
		ColStr:          strings.Join(columns, ", "),
		fieldsWithoutId: fieldsWithoutId,
		createStr:       createStr,
		placeholders:    placeholders,
		updateStr:       updateStr,
		zero:            reflect.New(refType).Elem().Interface().(E),
	}
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
	result := []E{}

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

func (em *EntityMetadata[E]) Delete(conn connection, id any) (int64, error) {
	sqlStr := em.buildDeleteById()
	result, err := em.doUpdate(conn, sqlStr, []any{id})
	if noError(err) {
		return result.RowsAffected()
	}
	return 0, err
}

func (em *EntityMetadata[E]) DeleteByQuery(conn connection, query any) (int64, error) {
	sqlStr, args := em.buildDelete(query)
	result, err := em.doUpdate(conn, sqlStr, args)
	if noError(err) {
		return result.RowsAffected()
	}
	return 0, err
}

func (em *EntityMetadata[E]) doUpdate(conn connection, sqlStr string, args []any) (sql.Result, error) {
	stmt, err := conn.Prepare(sqlStr)
	if noError(err) {
		defer func() {
			noError(stmt.Close())
		}()
		return stmt.Exec(args...)
	}
	return nil, err
}

func (em *EntityMetadata[E]) Create(conn connection, entity *E) (int64, error) {
	sqlStr, args := em.buildCreate(*entity)
	result, err := em.doUpdate(conn, sqlStr, args)
	if noError(err) {
		id, err := result.LastInsertId()
		if noError(err) {
			elem := reflect.ValueOf(entity).Elem()
			elem.FieldByName("Id").SetInt(id)
		}
		return id, err
	}
	return 0, err
}

func (em *EntityMetadata[E]) CreateMulti(conn connection, entities []E) (int64, error) {
	sqlStr, args := em.buildCreateMulti(entities)
	log.Debug("CREATE SQL: ", sqlStr)
	result, err := em.doUpdate(conn, sqlStr, args)
	if noError(err) {
		return result.RowsAffected()
	}
	return 0, err
}

func (em *EntityMetadata[E]) Update(conn connection, entity E) (int64, error) {
	sqlStr, args := em.buildUpdate(entity)
	result, err := em.doUpdate(conn, sqlStr, args)
	if noError(err) {
		return result.RowsAffected()
	}
	return 0, err
}

func (em *EntityMetadata[E]) Patch(conn connection, entity E) (int64, error) {
	sqlStr, args := em.buildPatchById(entity)
	result, err := em.doUpdate(conn, sqlStr, args)
	if noError(err) {
		return result.RowsAffected()
	}
	return 0, err
}

func (em *EntityMetadata[E]) PatchByQuery(conn connection, entity E, query GoQuery) (int64, error) {
	args, sqlStr := em.buildPatchByQuery(entity, query)
	result, err := em.doUpdate(conn, sqlStr, args)
	if noError(err) {
		return result.RowsAffected()
	}
	return 0, err
}
