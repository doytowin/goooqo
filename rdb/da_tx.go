package rdb

import (
	"context"
	. "github.com/doytowin/goooqo/core"
)

type txDataAccess[C context.Context, E any] struct {
	TransactionManager
	conn     Connection
	delegate DataAccess[ConnectionCtx, E]
}

type connectionCtx struct {
	context.Context
	Connection
}

func NewTxDataAccess[E Entity](tm TransactionManager, createEntity func() E) DataAccess[context.Context, E] {
	return &txDataAccess[context.Context, E]{
		TransactionManager: tm,
		conn:               tm.GetClient().(Connection),
		delegate:           newRelationalDataAccess[E](createEntity),
	}
}

// getConnCtx get connection from ctx, wrap the ctx and
// connection by ConnectionCtx as return value.
// ctx could be a TransactionContext with an active tx.
func (t *txDataAccess[C, E]) getConnCtx(ctx context.Context) ConnectionCtx {
	connection := t.conn
	if tc, ok := ctx.(*rdbTransactionContext); ok {
		connection = tc.tx
	}
	return &connectionCtx{ctx, connection}
}

func (t *txDataAccess[C, E]) Get(ctx C, id any) (*E, error) {
	return t.delegate.Get(t.getConnCtx(ctx), id)
}

func (t *txDataAccess[C, E]) Delete(ctx C, id any) (int64, error) {
	return t.delegate.Delete(t.getConnCtx(ctx), id)
}

func (t *txDataAccess[C, E]) Query(ctx C, query Query) ([]E, error) {
	return t.delegate.Query(t.getConnCtx(ctx), query)
}

func (t *txDataAccess[C, E]) Count(ctx C, query Query) (int64, error) {
	return t.delegate.Count(t.getConnCtx(ctx), query)
}

func (t *txDataAccess[C, E]) DeleteByQuery(ctx C, query Query) (int64, error) {
	return t.delegate.DeleteByQuery(t.getConnCtx(ctx), query)
}

func (t *txDataAccess[C, E]) Page(ctx C, query Query) (PageList[E], error) {
	return t.delegate.Page(t.getConnCtx(ctx), query)
}

func (t *txDataAccess[C, E]) Create(ctx C, entity *E) (int64, error) {
	return t.delegate.Create(t.getConnCtx(ctx), entity)
}

func (t *txDataAccess[C, E]) CreateMulti(ctx C, entities []E) (int64, error) {
	return t.delegate.CreateMulti(t.getConnCtx(ctx), entities)
}

func (t *txDataAccess[C, E]) Update(ctx C, entity E) (int64, error) {
	return t.delegate.Update(t.getConnCtx(ctx), entity)
}

func (t *txDataAccess[C, E]) Patch(ctx C, entity E) (int64, error) {
	return t.delegate.Patch(t.getConnCtx(ctx), entity)
}

func (t *txDataAccess[C, E]) PatchByQuery(ctx C, entity E, query Query) (int64, error) {
	return t.delegate.PatchByQuery(t.getConnCtx(ctx), entity, query)
}
