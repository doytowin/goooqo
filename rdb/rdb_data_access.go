/*
 * The Clear BSD License
 *
 * Copyright (c) 2024, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

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
// connection by ConnectionCtx as return value.
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
	return da.doQuery(ctx, sqlStr, args, query.GetPageSize())
}

func (da *relationalDataAccess[E]) doQuery(ctx context.Context, sqlStr string, args []any, size int) ([]E, error) {
	logSqlWithArgs(sqlStr, args)

	result := make([]E, 0, size)

	entity := *new(E)
	elem := reflect.ValueOf(&entity).Elem()
	columnMetas := da.em.columnMetas
	pointers := make([]any, len(columnMetas))
	for i, cm := range columnMetas {
		pointers[i] = elem.FieldByName(cm.Field.Name).Addr().Interface()
	}

	stmt, err := da.getConn(ctx).PrepareContext(ctx, sqlStr)
	if NoError(err) {
		defer Close(stmt)
		var rows *sql.Rows
		rows, err = stmt.QueryContext(ctx, args...)
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

func (da *relationalDataAccess[E]) Patch(ctx context.Context, entity E) (int64, error) {
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
