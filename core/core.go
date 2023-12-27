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

type Query interface {
	GetPageNumber() int
	GetPageSize() int
	CalcOffset() int
	GetSort() *string
	NeedPaging() bool
}

type Entity interface {
	GetTableName() string
}

type DataAccess[C context.Context, E any] interface {
	Get(ctx C, id any) (*E, error)
	Delete(ctx C, id any) (int64, error)
	Query(ctx C, query Query) ([]E, error)
	Count(ctx C, query Query) (int64, error)
	DeleteByQuery(ctx C, query Query) (int64, error)
	Page(ctx C, query Query) (PageList[E], error)
	Create(ctx C, entity *E) (int64, error)
	CreateMulti(ctx C, entities []E) (int64, error)
	Update(ctx C, entity E) (int64, error)
	Patch(ctx C, entity E) (int64, error)
	PatchByQuery(ctx C, entity E, query Query) (int64, error)
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

type TxDataAccess[E any, Q Query] interface {
	TransactionManager
	DataAccess[context.Context, E]
}
