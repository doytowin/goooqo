package goquery

import (
	"context"
	"github.com/doytowin/goquery/core"
	"github.com/doytowin/goquery/web"
	"net/http"
)

type GoQuery = core.GoQuery

type PageQuery = core.PageQuery

type Entity = core.Entity

type DataAccess[C any, E any] core.DataAccess[C, E]

func BuildRestService[E any, Q GoQuery](
	prefix string,
	dataAccess DataAccess[context.Context, E],
	createEntity func() E,
	createQuery func() Q,
) {
	s := web.NewRestService[E, Q](prefix, dataAccess, createEntity, createQuery)
	http.Handle(prefix, s)
}
