package goquery

import (
	"github.com/doytowin/goquery/core"
	"github.com/doytowin/goquery/web"
)

type GoQuery = core.GoQuery

type PageQuery = core.PageQuery

type Entity = core.Entity

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
