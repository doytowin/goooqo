package rdb

import (
	"context"
	"database/sql"
	. "github.com/doytowin/go-query/core"
)

type rdbTransactionManager struct {
	db *sql.DB
}

func (t *rdbTransactionManager) GetClient() any {
	return t.db
}

func NewTransactionManager(db *sql.DB) TransactionManager {
	return &rdbTransactionManager{db: db}
}

func getTx(ctx context.Context) *sql.Tx {
	if conn := ctx.Value(txKey); conn != nil {
		return conn.(*sql.Tx)
	}
	return nil
}

func (t *rdbTransactionManager) StartTransaction(ctx context.Context) TransactionContext {
	var err error
	tx := getTx(ctx)
	txCtx := ctx
	if tx == nil {
		tx, err = t.db.BeginTx(ctx, nil)
		if !NoError(err) {
			panic(err)
		}
		txCtx = context.WithValue(ctx, txKey, tx)
	} else {
		return ctx.(*rdbTransactionContext)
	}
	return &rdbTransactionContext{Context: txCtx, tx: tx}
}

type rdbTransactionContext struct {
	context.Context
	tx *sql.Tx
}

func (t rdbTransactionContext) Commit() error {
	return t.tx.Commit()
}

func (t rdbTransactionContext) Rollback() error {
	return t.tx.Rollback()
}
