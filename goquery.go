package goquery

import (
	"github.com/doytowin/goquery/core"
	"github.com/doytowin/goquery/rdb"
	"github.com/doytowin/goquery/web"
)

type GoQuery = core.GoQuery

type PageQuery = core.PageQuery

type Connection = core.Connection

type DataAccess[C any, E any] core.DataAccess[C, E]

func BuildController[C any, E any, Q GoQuery](
	prefix string, c C,
	dataAccess DataAccess[C, E],
	createEntity func() E,
	createQuery func() Q,
) *web.RestService[C, E, Q] {
	return &web.RestService[C, E, Q]{
		Service: web.BuildService[C, E, Q](prefix, c, dataAccess, createEntity, createQuery),
		Prefix:  prefix,
	}
}

func BuildRelationalDataAccess[E any](createEntity func() E) DataAccess[core.Connection, E] {
	return rdb.BuildRelationalDataAccess[E](createEntity)
}
