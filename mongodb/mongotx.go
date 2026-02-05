/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2026, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package mongodb

import (
	"context"

	. "github.com/doytowin/goooqo/core"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoTransactionManager struct {
	client *mongo.Client
}

func NewMongoTransactionManager(client *mongo.Client) TransactionManager {
	return &mongoTransactionManager{client: client}
}

func (tm *mongoTransactionManager) GetClient() any {
	return tm.client
}

func (tm *mongoTransactionManager) StartTransaction(ctx context.Context) (TransactionContext, error) {
	ssnCtx, err := tm.resolveCtx(ctx)
	if NoError(err) && !ssnCtx.active {
		err = ssnCtx.StartTransaction()
		if NoError(err) {
			ssnCtx.active = true
		}
	}
	return ssnCtx, err
}

func (tm *mongoTransactionManager) SubmitTransaction(ctx context.Context, callback func(tc TransactionContext) error) error {
	tc, err := tm.StartTransaction(ctx)
	if NoError(err) {
		err = TransactionCallback(tc, callback)
	}
	return err
}

func (tm *mongoTransactionManager) resolveCtx(ctx context.Context) (*mongoTransactionContext, error) {
	ssnCtx, ok := ctx.(*mongoTransactionContext)
	if !ok {
		sess, err := tm.client.StartSession()
		if err != nil {
			return nil, err
		}
		ssnCtx = &mongoTransactionContext{
			SessionContext: mongo.NewSessionContext(ctx, sess),
			active:         false,
			parent:         ctx,
		}
	}
	return ssnCtx, nil
}

type mongoTransactionContext struct {
	mongo.SessionContext
	active bool
	parent context.Context
}

func (t *mongoTransactionContext) Parent() context.Context {
	return t.parent
}

func (t *mongoTransactionContext) Commit() error {
	return t.CommitTransaction(t.SessionContext)
}

func (t *mongoTransactionContext) Rollback() error {
	return t.AbortTransaction(t.SessionContext)
}

func (t *mongoTransactionContext) SavePoint(name string) error {
	panic("not implemented")
}

func (t *mongoTransactionContext) RollbackTo(name string) error {
	panic("not implemented")
}
