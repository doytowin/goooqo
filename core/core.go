package core

import (
	"context"
	"database/sql/driver"
)

type PageList[E any] struct {
	List  []E   `json:"list"`
	Total int64 `json:"total"`
}

type Response struct {
	Data    any     `json:"data,omitempty"`
	Success bool    `json:"success"`
	Error   *string `json:"error,omitempty"`
}

type GoQuery interface {
	NeedPaging() bool
	BuildPageClause() string
	BuildSortClause() string
}

type Entity interface {
	GetTableName() string
}

type DataAccess[C context.Context, E any] interface {
	Get(c C, id any) (*E, error)
	Delete(c C, id any) (int64, error)
	Query(c C, query GoQuery) ([]E, error)
	Count(c C, query GoQuery) (int64, error)
	DeleteByQuery(c C, query any) (int64, error)
	Page(c C, query GoQuery) (PageList[E], error)
	Create(c C, entity *E) (int64, error)
	CreateMulti(c C, entities []E) (int64, error)
	Update(c C, entity E) (int64, error)
	Patch(c C, entity E) (int64, error)
	PatchByQuery(c C, entity E, query GoQuery) (int64, error)
}

type TransactionManager interface {
	GetClient() any
	StartTransaction(ctx context.Context) (TransactionContext, error)
}

type TransactionContext interface {
	context.Context
	driver.Tx
	Parent() context.Context
	SavePoint(name string) error
	RollbackTo(name string) error
}

type TxDataAccess[E any, Q GoQuery] interface {
	TransactionManager
	DataAccess[context.Context, E]
}
