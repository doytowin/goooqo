package rdb

import (
	"database/sql"
	. "github.com/doytowin/goquery/core"
	log "github.com/sirupsen/logrus"
	"reflect"
)

type RelationalDataAccess[E any] struct {
	em     EntityMetadata[E]
	create func() E
}

type connection = Connection

func noError(err error) bool {
	if err == nil {
		return true
	}
	log.Error("Error occurred! ", err)
	return false
}
func BuildDataAccess[E any](createEntity func() E) DataAccess[E] {
	e := buildRelationalDataAccess[E](createEntity)
	return &e
}

func buildRelationalDataAccess[E any](createEntity func() E) RelationalDataAccess[E] {
	em := buildEntityMetadata[E](createEntity())
	return RelationalDataAccess[E]{
		em:     em,
		create: createEntity,
	}
}

func (da *RelationalDataAccess[E]) Get(conn connection, id any) (*E, error) {
	sqlStr := da.em.buildSelectById()
	rows, err := da.doQuery(conn, sqlStr, []any{id})
	if len(rows) == 1 {
		return &rows[0], err
	}
	return nil, err
}

func (da *RelationalDataAccess[E]) Query(conn connection, query GoQuery) ([]E, error) {
	sqlStr, args := da.em.buildSelect(query)
	return da.doQuery(conn, sqlStr, args)
}

func (da *RelationalDataAccess[E]) doQuery(conn connection, sqlStr string, args []any) ([]E, error) {
	result := []E{}

	entity := da.create()
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
