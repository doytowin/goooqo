package rdb

import (
	"context"
	. "github.com/doytowin/goquery/core"
)

type Ck struct {
}

var connKey = Ck{}

type txDataAccess[C context.Context, E any] struct {
	conn     Connection
	delegate DataAccess[Connection, E]
}

func NewTxDataAccess[E Entity](conn Connection, createEntity func() E) DataAccess[context.Context, E] {
	return &txDataAccess[context.Context, E]{
		conn:     conn,
		delegate: BuildRelationalDataAccess[E](createEntity),
	}
}

func (t *txDataAccess[C, E]) getConn(ctx C) Connection {
	if conn := ctx.Value(connKey); conn != nil {
		return conn.(Connection)
	}
	return t.conn
}

func (t *txDataAccess[C, E]) Get(ctx C, id any) (*E, error) {
	return t.delegate.Get(t.getConn(ctx), id)
}

func (t *txDataAccess[C, E]) Delete(ctx C, id any) (int64, error) {
	return t.delegate.Delete(t.getConn(ctx), id)
}

func (t *txDataAccess[C, E]) Query(ctx C, query GoQuery) ([]E, error) {
	return t.delegate.Query(t.getConn(ctx), query)
}

func (t *txDataAccess[C, E]) Count(ctx C, query GoQuery) (int64, error) {
	return t.delegate.Count(t.getConn(ctx), query)
}

func (t *txDataAccess[C, E]) DeleteByQuery(ctx C, query any) (int64, error) {
	return t.delegate.DeleteByQuery(t.getConn(ctx), query)
}

func (t *txDataAccess[C, E]) Page(ctx C, query GoQuery) (PageList[E], error) {
	return t.delegate.Page(t.getConn(ctx), query)
}

func (t *txDataAccess[C, E]) Create(ctx C, entity *E) (int64, error) {
	return t.delegate.Create(t.getConn(ctx), entity)
}

func (t *txDataAccess[C, E]) CreateMulti(ctx C, entities []E) (int64, error) {
	return t.delegate.CreateMulti(t.getConn(ctx), entities)
}

func (t *txDataAccess[C, E]) Update(ctx C, entity E) (int64, error) {
	return t.delegate.Update(t.getConn(ctx), entity)
}

func (t *txDataAccess[C, E]) Patch(ctx C, entity E) (int64, error) {
	return t.delegate.Patch(t.getConn(ctx), entity)
}

func (t *txDataAccess[C, E]) PatchByQuery(ctx C, entity E, query GoQuery) (int64, error) {
	return t.delegate.PatchByQuery(t.getConn(ctx), entity, query)
}
