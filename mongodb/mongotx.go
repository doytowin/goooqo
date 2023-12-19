package mongodb

import (
	"context"
	. "github.com/doytowin/go-query/core"
	"go.mongodb.org/mongo-driver/mongo"
)

var sessionKey = struct{}{}

type mongoTransactionManager struct {
	client *mongo.Client
}

func NewMongoTransactionManager(client *mongo.Client) TransactionManager {
	return &mongoTransactionManager{client: client}
}

func (tm *mongoTransactionManager) GetClient() any {
	return tm.client
}

func (tm *mongoTransactionManager) StartTransaction(ctx context.Context) TransactionContext {
	ssnCtx := tm.resolveCtx(ctx)
	if !ssnCtx.active {
		err := ssnCtx.StartTransaction()
		ssnCtx.active = true
		if !NoError(err) {
			panic(err)
		}
	}
	return ssnCtx
}

func (tm *mongoTransactionManager) resolveCtx(ctx context.Context) *mongoTransactionController {
	ssnCtx, ok := ctx.(*mongoTransactionController)
	if !ok {
		sess, err := tm.client.StartSession()
		if NoError(err) {
			ssnCtx = &mongoTransactionController{
				SessionContext: mongo.NewSessionContext(ctx, sess),
				active:         false,
			}
		} else {
			panic(err)
		}
	}
	return ssnCtx
}

type mongoTransactionController struct {
	mongo.SessionContext
	active bool
}

func (t *mongoTransactionController) Commit() error {
	return t.CommitTransaction(t.SessionContext)
}

func (t *mongoTransactionController) Rollback() error {
	return t.AbortTransaction(t.SessionContext)
}
