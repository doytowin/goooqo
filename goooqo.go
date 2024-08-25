package goooqo

import (
	"github.com/doytowin/goooqo/core"
	"github.com/doytowin/goooqo/web"
	"net/http"
)

type Query = core.Query

type PageQuery = core.PageQuery

type Entity = core.Entity

type Int64Id = core.Int64Id

type DataAccess[E Entity] core.DataAccess[E]

type TransactionManager = core.TransactionManager

func BuildRestService[E Entity, Q Query](
	prefix string,
	dataAccess DataAccess[E],
) {
	s := web.NewRestService[E, Q](prefix, dataAccess)
	http.Handle(prefix, s)
}
