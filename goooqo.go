/*
 * The Clear BSD License
 *
 * Copyright (c) 2024, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package goooqo

import (
	"github.com/doytowin/goooqo/core"
	"github.com/doytowin/goooqo/web"
	"net/http"
)

type Query = core.Query

type PageQuery = core.PageQuery

type Entity = core.Entity

type IntId = core.IntId

var NewIntId = core.NewIntId

type Int64Id = core.Int64Id

var NewInt64Id = core.NewInt64Id

func P[T any](t T) *T { return &t }

type DataAccess[E Entity] core.DataAccess[E]

type TransactionManager = core.TransactionManager

func BuildRestService[E Entity, Q Query](
	prefix string,
	dataAccess DataAccess[E],
) {
	s := web.NewRestService[E, Q](prefix, dataAccess)
	http.Handle(prefix, s)
}
