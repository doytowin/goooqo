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
}

type Entity interface {
	GetTableName() string
}

type DataAccess[C any, E any] interface {
	Get(conn C, id any) (*E, error)
	Delete(conn C, id any) (int64, error)
	Query(conn C, query GoQuery) ([]E, error)
	Count(conn C, query GoQuery) (int64, error)
	DeleteByQuery(conn C, query any) (int64, error)
	Page(conn C, query GoQuery) (PageList[E], error)
	Create(conn C, entity *E) (int64, error)
	CreateMulti(conn C, entities []E) (int64, error)
	Update(conn C, entity E) (int64, error)
	Patch(conn C, entity E) (int64, error)
	PatchByQuery(conn C, entity E, query GoQuery) (int64, error)
}

type TransactionManager interface {
	GetClient() any
	StartTransaction(ctx context.Context) TransactionContext
}

type TransactionContext interface {
	context.Context
	driver.Tx
}

type TxDataAccess[E any, Q GoQuery] interface {
	TransactionManager
	DataAccess[context.Context, E]
}
