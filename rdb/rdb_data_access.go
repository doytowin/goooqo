/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2025, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package rdb

import (
	"context"
	"database/sql"
	"reflect"

	. "github.com/doytowin/goooqo/core"
	log "github.com/sirupsen/logrus"
)

type Connection interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type relationalDataAccess[E Entity] struct {
	TransactionManager
	conn Connection
	em   EntityMetadata[E]
}

func logSqlWithArgs(sqlStr string, args []any) (string, []any) {
	log.WithFields(log.Fields{"SQL": sqlStr, "args": args}).Info("Executing")
	return sqlStr, args
}

func NewTxDataAccess[E Entity](tm TransactionManager) TxDataAccess[E] {
	return &relationalDataAccess[E]{
		TransactionManager: tm,
		conn:               tm.GetClient().(Connection),
		em:                 buildEntityMetadata[E](),
	}
}

// getConn get connection from ctx, wrap the ctx and
// connection by Connection as return value.
// ctx could be a TransactionContext with an active tx.
func (da *relationalDataAccess[E]) getConn(ctx context.Context) Connection {
	if tc, ok := ctx.(*rdbTransactionContext); ok {
		return tc.tx
	}
	return da.conn
}

func (da *relationalDataAccess[E]) Get(ctx context.Context, id any) (*E, error) {
	sqlStr := da.em.buildSelectById()
	rows, err := da.doQuery(ctx, sqlStr, []any{id}, 1)
	if len(rows) == 1 {
		return &rows[0], err
	}
	return nil, err
}

func (da *relationalDataAccess[E]) Query(ctx context.Context, query Query) ([]E, error) {
	sqlStr, args := da.em.buildSelect(query)
	entities, err := da.doQuery(ctx, sqlStr, args, query.GetPageSize())
	if NoError(err) && len(da.em.relationMetas) > 0 {
		da.queryRelationEntities(ctx, entities, query)
	}
	return entities, err
}

func (da *relationalDataAccess[E]) doQuery(ctx context.Context, sqlStr string, args []any, size int) ([]E, error) {
	logSqlWithArgs(sqlStr, args)

	result := make([]E, 0, size)

	entity := *new(E)

	stmt, err := da.getConn(ctx).PrepareContext(ctx, sqlStr)
	if NoError(err) {
		defer Close(stmt)
		var rows *sql.Rows
		rows, err = stmt.QueryContext(ctx, args...)
		if NoError(err) {
			var pointers []any
			if mapper, ok := any(entity).(EntityMapper); ok {
				pointers = mapper.FieldsAddr()
			} else {
				pointers = preparePointers(reflect.ValueOf(&entity), da.em.columnMetas)
			}
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

func preparePointers(p reflect.Value, fieldMetas []FieldMetadata) []any {
	pointers := make([]any, len(fieldMetas))
	for i, fm := range fieldMetas {
		pointers[i] = p.Elem().FieldByName(fm.Field.Name).Addr().Interface()
	}
	return pointers
}

func (da *relationalDataAccess[E]) queryRelationEntities(ctx context.Context, entities []E, query Query) {
	elem := reflect.ValueOf(query)
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}
	for _, rm := range da.em.relationMetas {
		queryName := "With" + rm.Field.Name
		entityQueryVal := elem.FieldByName(queryName)
		if !entityQueryVal.IsNil() {
			ep := fpEntityPath{*rm.EntityPath}
			sqlStr, args := ep.buildQuery(entityQueryVal.Interface().(Query))

			for i, entity := range entities {
				relatedEntities, err := QueryRelated(ctx, da.getConn(ctx), sqlStr,
					append([]any{entity.GetId()}, args...), ep.EntityType)
				if NoError(err) {
					reflect.ValueOf(&entities[i]).Elem().FieldByName(rm.Field.Name).Set(relatedEntities)
				}
			}
		}
	}
}

func QueryRelated(ctx context.Context, conn Connection, sqlStr string, args []any, entityType reflect.Type) (reflect.Value, error) {
	logSqlWithArgs(sqlStr, args)

	stmt, err := conn.PrepareContext(ctx, sqlStr)
	result := reflect.MakeSlice(reflect.SliceOf(entityType), 0, 10)
	if NoError(err) {
		defer Close(stmt)
		var rows *sql.Rows
		rows, err = stmt.QueryContext(ctx, args...)
		if NoError(err) {
			var pointers []any
			pEntity := reflect.New(entityType)
			if mapper, ok := pEntity.Interface().(EntityMapper); ok {
				pointers = mapper.FieldsAddr()
			} else {
				fmArr := retainColumns(BuildFieldMetas(entityType))
				pointers = preparePointers(pEntity, fmArr)
			}
			for rows.Next() {
				err = rows.Scan(pointers...)
				if NoError(err) {
					result = reflect.Append(result, pEntity.Elem())
				}
			}
		}
	}

	return result, err
}

func retainColumns(fieldMetas []FieldMetadata) []FieldMetadata {
	columnMetas := make([]FieldMetadata, 0, len(fieldMetas))
	for _, md := range fieldMetas {
		if md.EntityPath == nil {
			columnMetas = append(columnMetas, md)
		}
	}
	return columnMetas
}

func (da *relationalDataAccess[E]) Count(ctx context.Context, query Query) (int64, error) {
	var cnt int64
	sqlStr, args := da.em.buildCount(query)
	logSqlWithArgs(sqlStr, args)
	stmt, err := da.getConn(ctx).PrepareContext(ctx, sqlStr)
	if NoError(err) {
		defer Close(stmt)
		row := stmt.QueryRowContext(ctx, args...)
		err = row.Scan(&cnt)
	}
	return cnt, err
}

func (da *relationalDataAccess[E]) Page(ctx context.Context, query Query) (PageList[E], error) {
	var cnt int64
	data, err := da.Query(ctx, query)
	if NoError(err) {
		cnt, err = da.Count(ctx, query)
	}
	return PageList[E]{List: data, Total: cnt}, err
}

func (da *relationalDataAccess[E]) Delete(ctx context.Context, id any) (int64, error) {
	sqlStr := da.em.buildDeleteById()
	return parse(da.doUpdate(ctx, sqlStr, []any{id}))
}

func (da *relationalDataAccess[E]) DeleteByQuery(ctx context.Context, query Query) (int64, error) {
	sqlStr, args := da.em.buildDelete(query)
	return parse(da.doUpdate(ctx, sqlStr, args))
}

func (da *relationalDataAccess[E]) doUpdate(ctx context.Context, sqlStr string, args []any) (sql.Result, error) {
	logSqlWithArgs(sqlStr, args)
	stmt, err := da.getConn(ctx).PrepareContext(ctx, sqlStr)
	if NoError(err) {
		defer Close(stmt)
		return stmt.ExecContext(ctx, args...)
	}
	return nil, err
}

func (da *relationalDataAccess[E]) Create(ctx context.Context, entity *E) (int64, error) {
	sqlStr, args := da.em.buildCreate(*entity)
	result, err := da.doUpdate(ctx, sqlStr, args)
	var id int64
	if NoError(err) {
		id, err = result.LastInsertId()
		if NoError(err) {
			err = (*entity).SetId(entity, id)
		}
	}
	return id, err
}

func (da *relationalDataAccess[E]) CreateMulti(ctx context.Context, entities []E) (int64, error) {
	if len(entities) == 0 {
		return 0, nil
	}
	sqlStr, args := da.em.buildCreateMulti(entities)
	return parse(da.doUpdate(ctx, sqlStr, args))
}

func (da *relationalDataAccess[E]) Update(ctx context.Context, entity E) (int64, error) {
	sqlStr, args := da.em.buildUpdate(entity)
	return parse(da.doUpdate(ctx, sqlStr, args))
}

func (da *relationalDataAccess[E]) Patch(ctx context.Context, entity Entity) (int64, error) {
	sqlStr, args := da.em.buildPatchById(entity)
	return parse(da.doUpdate(ctx, sqlStr, args))
}

func (da *relationalDataAccess[E]) PatchByQuery(ctx context.Context, entity E, query Query) (int64, error) {
	sqlStr, args := da.em.buildPatchByQuery(entity, query)
	return parse(da.doUpdate(ctx, sqlStr, args))
}

func parse(result sql.Result, err error) (int64, error) {
	if NoError(err) {
		return result.RowsAffected()
	}
	return 0, err
}
