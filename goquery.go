package goquery

import (
	"context"
	"github.com/doytowin/goquery/core"
	"github.com/doytowin/goquery/web"
)

type GoQuery = core.GoQuery

type PageQuery = core.PageQuery

type Entity = core.Entity

type DataAccess[C any, E any] core.DataAccess[C, E]

func BuildController[E any, Q GoQuery](
	prefix string,
	dataAccess DataAccess[context.Context, E],
	createEntity func() E,
	createQuery func() Q,
) *web.RestService[E, Q] {
	return &web.RestService[E, Q]{
		Service: web.BuildService[E, Q](prefix, dataAccess, createEntity, createQuery),
		Prefix:  prefix,
	}
}
