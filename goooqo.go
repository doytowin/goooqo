package goooqo

import (
	"context"
	"github.com/doytowin/goooqo/core"
	"github.com/doytowin/goooqo/web"
	"net/http"
)

type Query = core.Query

type PageQuery = core.PageQuery

type Entity = core.Entity

type DataAccess[C context.Context, E Entity] core.DataAccess[C, E]

type TransactionManager = core.TransactionManager

func BuildRestService[E Entity, Q Query](
	prefix string,
	dataAccess DataAccess[context.Context, E],
	createEntity func() E,
	createQuery func() Q,
) {
	s := web.NewRestService[E, Q](prefix, dataAccess, createEntity, createQuery)
	http.Handle(prefix, s)
}
