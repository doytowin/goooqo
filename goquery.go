package goquery

import (
	"database/sql"
	"github.com/doytowin/goquery/core"
	"github.com/doytowin/goquery/rdb"
	"github.com/doytowin/goquery/web"
)

type GoQuery = core.GoQuery

type PageQuery = core.PageQuery

type DataAccess[E any] core.DataAccess[E]

func BuildController[E any, Q GoQuery](
	prefix string,
	db *sql.DB,
	createEntity func() E,
	createQuery func() Q,
) *web.RestService[E, Q] {
	return &web.RestService[E, Q]{
		Service: web.BuildService(prefix, db, createEntity, createQuery),
	}
}

func BuildDataAccess[E any](createEntity func() E) DataAccess[E] {
	return rdb.BuildDataAccess[E](createEntity)
}
