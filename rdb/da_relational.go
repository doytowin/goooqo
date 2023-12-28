package rdb

import (
	"context"
	"database/sql"
	. "github.com/doytowin/goooqo/core"
	log "github.com/sirupsen/logrus"
	"reflect"
)

type Connection interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

type ConnectionCtx interface {
	context.Context
	Connection
}

type relationalDataAccess[C ConnectionCtx, E Entity] struct {
	em     EntityMetadata[E]
	create func() E
}

func newRelationalDataAccess[E RdbEntity](createEntity func() E) DataAccess[ConnectionCtx, E] {
	return &relationalDataAccess[ConnectionCtx, E]{
		em:     buildEntityMetadata[E](createEntity()),
		create: createEntity,
	}
}

func (da *relationalDataAccess[C, E]) Get(connCtx C, id any) (*E, error) {
	sqlStr := da.em.buildSelectById()
	rows, err := da.doQuery(connCtx, sqlStr, []any{id})
	if len(rows) == 1 {
		return &rows[0], err
	}
	return nil, err
}

func (da *relationalDataAccess[C, E]) Query(connCtx C, query Query) ([]E, error) {
	sqlStr, args := da.em.buildSelect(query)
	return da.doQuery(connCtx, sqlStr, args)
}

func (da *relationalDataAccess[C, E]) doQuery(connCtx C, sqlStr string, args []any) ([]E, error) {
	log.Debug("SQL: ", sqlStr)
	log.Debug("ARG: ", args)

	result := []E{}

	entity := da.create()
	elem := reflect.ValueOf(&entity).Elem()
	columnMetas := da.em.columnMetas
	pointers := make([]any, len(columnMetas))
	for i, cm := range columnMetas {
		pointers[i] = elem.FieldByName(cm.field.Name).Addr().Interface()
	}

	stmt, err := connCtx.PrepareContext(connCtx, sqlStr)
	if NoError(err) {
		defer Close(stmt)
		var rows *sql.Rows
		rows, err = stmt.QueryContext(connCtx, args...)
		if NoError(err) {
			for rows.Next() {
				err = rows.Scan(pointers...)
				if NoError(err) {
					result = append(result, entity)
				}
			}
		}
	}

	return result, err
}

func (da *relationalDataAccess[C, E]) Count(connCtx C, query Query) (int64, error) {
	var cnt int64
	sqlStr, args := da.em.buildCount(query)
	stmt, err := connCtx.PrepareContext(connCtx, sqlStr)
	if NoError(err) {
		defer Close(stmt)
		row := stmt.QueryRowContext(connCtx, args...)
		err = row.Scan(&cnt)
	}
	return cnt, err
}

func (da *relationalDataAccess[C, E]) Page(connCtx C, query Query) (PageList[E], error) {
	var cnt int64
	data, err := da.Query(connCtx, query)
	if NoError(err) {
		cnt, err = da.Count(connCtx, query)
	}
	return PageList[E]{List: data, Total: cnt}, err
}

func (da *relationalDataAccess[C, E]) Delete(connCtx C, id any) (int64, error) {
	sqlStr := da.em.buildDeleteById()
	return parse(da.doUpdate(connCtx, sqlStr, []any{id}))
}

func (da *relationalDataAccess[C, E]) DeleteByQuery(connCtx C, query Query) (int64, error) {
	sqlStr, args := da.em.buildDelete(query)
	return parse(da.doUpdate(connCtx, sqlStr, args))
}

func (da *relationalDataAccess[C, E]) doUpdate(connCtx C, sqlStr string, args []any) (sql.Result, error) {
	stmt, err := connCtx.PrepareContext(connCtx, sqlStr)
	if NoError(err) {
		defer Close(stmt)
		return stmt.ExecContext(connCtx, args...)
	}
	return nil, err
}

func (da *relationalDataAccess[C, E]) Create(connCtx C, entity *E) (int64, error) {
	sqlStr, args := da.em.buildCreate(*entity)
	result, err := da.doUpdate(connCtx, sqlStr, args)
	var id int64
	if NoError(err) {
		id, err = result.LastInsertId()
		if NoError(err) {
			(*entity).SetId(entity, id)
		}
	}
	return id, err
}

func (da *relationalDataAccess[C, E]) CreateMulti(connCtx C, entities []E) (int64, error) {
	if len(entities) == 0 {
		return 0, nil
	}
	sqlStr, args := da.em.buildCreateMulti(entities)
	return parse(da.doUpdate(connCtx, sqlStr, args))
}

func (da *relationalDataAccess[C, E]) Update(connCtx C, entity E) (int64, error) {
	sqlStr, args := da.em.buildUpdate(entity)
	return parse(da.doUpdate(connCtx, sqlStr, args))
}

func (da *relationalDataAccess[C, E]) Patch(connCtx C, entity E) (int64, error) {
	sqlStr, args := da.em.buildPatchById(entity)
	return parse(da.doUpdate(connCtx, sqlStr, args))
}

func (da *relationalDataAccess[C, E]) PatchByQuery(connCtx C, entity E, query Query) (int64, error) {
	args, sqlStr := da.em.buildPatchByQuery(entity, query)
	return parse(da.doUpdate(connCtx, sqlStr, args))
}

func parse(result sql.Result, err error) (int64, error) {
	if NoError(err) {
		return result.RowsAffected()
	}
	return 0, err
}
