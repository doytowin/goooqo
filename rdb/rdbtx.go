package rdb

import (
	"context"
	"database/sql"
	. "github.com/doytowin/go-query/core"
	log "github.com/sirupsen/logrus"
	"sync/atomic"
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
	if !NoError(err) {
		return nil, err
	}
	sn := t.fetchSn()
	log.Debug("Start Tx: ", sn)
	return &rdbTransactionContext{Context: ctx, tx: tx, sn: sn}, nil
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
