/*
 * The Clear BSD License
 *
 * Copyright (c) 2024, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package core

import (
	"context"
	"database/sql/driver"
)

type PageList[D any] struct {
	List  []D   `json:"list"`
	Total int64 `json:"total"`
}

type Response struct {
	Data    any     `json:"data,omitempty"`
	Success bool    `json:"success"`
	Error   *string `json:"error,omitempty"`
}

type Query interface {
	GetPageNumber() int
	GetPageSize() int
	CalcOffset() int
	GetSort() *string
	NeedPaging() bool
}

type Entity interface {
	GetId() any

	// SetId set id to self.
	// self: the pointer point to the current entity.
	// id: type could be int64 or string so far.
	SetId(self any, id any) error
}

type DataAccess[E Entity] interface {
	Get(ctx context.Context, id any) (*E, error)
	Delete(ctx context.Context, id any) (int64, error)
	Query(ctx context.Context, query Query) ([]E, error)
	Count(ctx context.Context, query Query) (int64, error)
	DeleteByQuery(ctx context.Context, query Query) (int64, error)
	Page(ctx context.Context, query Query) (PageList[E], error)
	Create(ctx context.Context, entity *E) (int64, error)
	CreateMulti(ctx context.Context, entities []E) (int64, error)
	Update(ctx context.Context, entity E) (int64, error)
	Patch(ctx context.Context, entity E) (int64, error)
	PatchByQuery(ctx context.Context, entity E, query Query) (int64, error)
}

type TransactionManager interface {
	GetClient() any
	StartTransaction(ctx context.Context) (TransactionContext, error)
	SubmitTransaction(ctx context.Context, callback func(tc TransactionContext) error) error
}

func TransactionCallback(tc TransactionContext, callback func(tc TransactionContext) error) (err error) {
	err = callback(tc)
	if NoError(err) {
		err = tc.Commit()
	} else {
		NoError(tc.Rollback())
	}
	return
}

type TransactionContext interface {
	context.Context
	driver.Tx
	Parent() context.Context
	SavePoint(name string) error
	RollbackTo(name string) error
}

type TxDataAccess[E Entity] interface {
	TransactionManager
	DataAccess[E]
}
