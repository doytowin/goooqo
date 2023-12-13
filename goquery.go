package goquery

import (
	"database/sql"
	. "github.com/doytowin/goquery/util"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strings"
)

type RelationalDataAccess[E comparable] struct {
	em   EntityMetadata[E]
	Type reflect.Type
	zero E
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

func buildRelationalDataAccess[E comparable](entity any) RelationalDataAccess[E] {
	em := buildEntityMetadata[E](entity)
	return RelationalDataAccess[E]{
		em:   em,
		Type: reflect.TypeOf(entity),
		zero: reflect.New(reflect.TypeOf(entity)).Elem().Interface().(E),
	}
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
		TableName:       tableName,
		ColStr:          strings.Join(columns, ", "),
		fieldsWithoutId: fieldsWithoutId,
		createStr:       createStr,
		placeholders:    placeholders,
		updateStr:       updateStr,
	}
}

func (da *RelationalDataAccess[E]) Get(conn connection, id any) (E, error) {
	sqlStr := da.em.buildSelectById()
	rows, err := da.doQuery(conn, sqlStr, []any{id})
	if len(rows) == 1 {
		return rows[0], err
	}
	return da.zero, err
}

func (da *RelationalDataAccess[E]) Query(conn connection, query GoQuery) ([]E, error) {
	sqlStr, args := da.em.buildSelect(query)
	return da.doQuery(conn, sqlStr, args)
}

func (da *RelationalDataAccess[E]) doQuery(conn connection, sqlStr string, args []any) ([]E, error) {
	result := []E{}

	entity := reflect.New(da.Type).Elem().Interface().(E)
	elem := reflect.ValueOf(&entity).Elem()
	length := elem.NumField()
	pointers := make([]any, length)
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

func (da *RelationalDataAccess[E]) Count(conn connection, query GoQuery) (int, error) {
	cnt := 0
	sqlStr, args := da.em.buildCount(query)
	stmt, err := conn.Prepare(sqlStr)
	if noError(err) {
		row := stmt.QueryRow(args...)
		err = row.Scan(&cnt)
		_ = stmt.Close()
	}
	return cnt, err
}

func (da *RelationalDataAccess[E]) Page(conn connection, query GoQuery) (PageList[E], error) {
	var count int
	data, err := da.Query(conn, query)
	if noError(err) {
		count, err = da.Count(conn, query)
	}
	return PageList[E]{data, count}, err
}

func (da *RelationalDataAccess[E]) IsZero(entity E) bool {
	return da.zero == entity
}

func (da *RelationalDataAccess[E]) Delete(conn connection, id any) (int64, error) {
	sqlStr := da.em.buildDeleteById()
	result, err := da.doUpdate(conn, sqlStr, []any{id})
	if noError(err) {
		return result.RowsAffected()
	}
	return 0, err
}

func (da *RelationalDataAccess[E]) DeleteByQuery(conn connection, query any) (int64, error) {
	sqlStr, args := da.em.buildDelete(query)
	result, err := da.doUpdate(conn, sqlStr, args)
	if noError(err) {
		return result.RowsAffected()
	}
	return 0, err
}

func (da *RelationalDataAccess[E]) doUpdate(conn connection, sqlStr string, args []any) (sql.Result, error) {
	stmt, err := conn.Prepare(sqlStr)
	if noError(err) {
		defer func() {
			noError(stmt.Close())
		}()
		return stmt.Exec(args...)
	}
	return nil, err
}

func (da *RelationalDataAccess[E]) Create(conn connection, entity *E) (int64, error) {
	sqlStr, args := da.em.buildCreate(*entity)
	result, err := da.doUpdate(conn, sqlStr, args)
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

func (da *RelationalDataAccess[E]) CreateMulti(conn connection, entities []E) (int64, error) {
	if len(entities) == 0 {
		return 0, nil
	}
	sqlStr, args := da.em.buildCreateMulti(entities)
	log.Debug("CREATE SQL: ", sqlStr)
	result, err := da.doUpdate(conn, sqlStr, args)
	if noError(err) {
		return result.RowsAffected()
	}
	return 0, err
}

func (da *RelationalDataAccess[E]) Update(conn connection, entity E) (int64, error) {
	sqlStr, args := da.em.buildUpdate(entity)
	result, err := da.doUpdate(conn, sqlStr, args)
	if noError(err) {
		return result.RowsAffected()
	}
	return 0, err
}

func (da *RelationalDataAccess[E]) Patch(conn connection, entity E) (int64, error) {
	sqlStr, args := da.em.buildPatchById(entity)
	result, err := da.doUpdate(conn, sqlStr, args)
	if noError(err) {
		return result.RowsAffected()
	}
	return 0, err
}

func (da *RelationalDataAccess[E]) PatchByQuery(conn connection, entity E, query GoQuery) (int64, error) {
	args, sqlStr := da.em.buildPatchByQuery(entity, query)
	result, err := da.doUpdate(conn, sqlStr, args)
	if noError(err) {
		return result.RowsAffected()
	}
	return 0, err
}
