package rdb

import (
	"context"
	. "github.com/doytowin/go-query/core"
)

var txKey = struct{ string }{"tx-key"}

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

// getConnCtx get connection from ctx, wrap the ctx
// and connection by ConnectionCtx as return value.
// ctx could be a TransactionContext with key txKey
// and value tx for transaction operation.
func (t *txDataAccess[C, E]) getConnCtx(ctx C) ConnectionCtx {
	connection := t.conn
	if conn := ctx.Value(txKey); conn != nil {
		connection = conn.(Connection)
	}
	return &connectionCtx{ctx, connection}
}

func (t *txDataAccess[C, E]) Get(ctx C, id any) (*E, error) {
	return t.delegate.Get(t.getConnCtx(ctx), id)
}

func (t *txDataAccess[C, E]) Delete(ctx C, id any) (int64, error) {
	return t.delegate.Delete(t.getConnCtx(ctx), id)
}

func (t *txDataAccess[C, E]) Query(ctx C, query GoQuery) ([]E, error) {
	return t.delegate.Query(t.getConnCtx(ctx), query)
}

func (t *txDataAccess[C, E]) Count(ctx C, query GoQuery) (int64, error) {
	return t.delegate.Count(t.getConnCtx(ctx), query)
}

func (t *txDataAccess[C, E]) DeleteByQuery(ctx C, query any) (int64, error) {
	return t.delegate.DeleteByQuery(t.getConnCtx(ctx), query)
}

func (t *txDataAccess[C, E]) Page(ctx C, query GoQuery) (PageList[E], error) {
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

func (t *txDataAccess[C, E]) PatchByQuery(ctx C, entity E, query GoQuery) (int64, error) {
	return t.delegate.PatchByQuery(t.getConnCtx(ctx), entity, query)
}
