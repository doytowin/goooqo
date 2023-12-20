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

func (t *rdbTransactionManager) StartTransaction(ctx context.Context) (TransactionContext, error) {
	if tc, ok := ctx.(*rdbTransactionContext); ok {
		return tc, nil
	}
	tx, err := t.db.BeginTx(ctx, nil)
	if !NoError(err) {
		return nil, err
	}
	txCtx := context.WithValue(ctx, txKey, tx)
	return &rdbTransactionContext{Context: txCtx, tx: tx}, nil
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
