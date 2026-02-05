/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2026, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package rdb

import (
	"context"
	"database/sql"
	"sync/atomic"

	. "github.com/doytowin/goooqo/core"
	log "github.com/sirupsen/logrus"
)

type rdbTransactionManager struct {
	db *sql.DB
	sn *atomic.Value
}

func NewTransactionManager(db *sql.DB) TransactionManager {
	sn := &atomic.Value{}
	sn.Store(int64(0))
	return &rdbTransactionManager{db: db, sn: sn}
}

func (t *rdbTransactionManager) GetClient() any {
	return t.db
}

func (t *rdbTransactionManager) StartTransaction(ctx context.Context) (TransactionContext, error) {
	if tc, ok := ctx.(*rdbTransactionContext); ok {
		return tc, nil
	}
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	sn := t.fetchSn()
	log.Debug("Start Tx: ", sn)
	return &rdbTransactionContext{Context: ctx, tx: tx, sn: sn}, nil
}

func (t *rdbTransactionManager) SubmitTransaction(ctx context.Context, callback func(tc TransactionContext) error) error {
	tc, err := t.StartTransaction(ctx)
	if err == nil {
		err = TransactionCallback(tc, callback)
	}
	return err
}

func (t *rdbTransactionManager) fetchSn() int64 {
	var val = t.sn.Load().(int64)
	for !t.sn.CompareAndSwap(val, val+1) {
		val = t.sn.Load().(int64)
	}
	return val + 1
}

type rdbTransactionContext struct {
	context.Context
	tx *sql.Tx
	sn int64
}

func (t *rdbTransactionContext) Commit() error {
	log.Debug("Commit Tx: ", t.sn)
	return t.tx.Commit()
}

func (t *rdbTransactionContext) Rollback() error {
	log.Debug("Rollback Tx: ", t.sn)
	return t.tx.Rollback()
}

func (t *rdbTransactionContext) Parent() context.Context {
	return t.Context
}

func (t *rdbTransactionContext) SavePoint(name string) error {
	_, err := t.tx.ExecContext(t.Context, "SAVEPOINT "+name)
	return err
}

func (t *rdbTransactionContext) RollbackTo(name string) error {
	_, err := t.tx.ExecContext(t.Context, "ROLLBACK TO SAVEPOINT "+name)
	return err
}
