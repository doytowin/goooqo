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

func (tm *mongoTransactionManager) resolveCtx(ctx context.Context) (*mongoTransactionContext, error) {
	ssnCtx, ok := ctx.(*mongoTransactionContext)
	if !ok {
		sess, err := tm.client.StartSession()
		if HasError(err) {
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
