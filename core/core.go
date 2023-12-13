package core

import "database/sql"

type GoQuery interface {
	GetPageQuery() *PageQuery
}

type Entity interface {
	GetTableName() string
}

type PageList[E any] struct {
	List  []E
	Total int
}

type connection interface {
	Prepare(query string) (*sql.Stmt, error)
}

type Connection = connection

type DataAccess[E any] interface {
	Get(conn connection, id any) (*E, error)
	Delete(conn connection, id any) (int64, error)
	Query(conn connection, query GoQuery) ([]E, error)
	Count(conn connection, query GoQuery) (int, error)
	DeleteByQuery(conn connection, query any) (int64, error)
	Page(conn connection, query GoQuery) (PageList[E], error)
	Create(conn connection, entity *E) (int64, error)
	CreateMulti(conn connection, entities []E) (int64, error)
	Update(conn connection, entity E) (int64, error)
	Patch(conn connection, entity E) (int64, error)
	PatchByQuery(conn connection, entity E, query GoQuery) (int64, error)
}

type Response struct {
	Data    any
	Success bool
	Error   *string
}
