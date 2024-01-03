package mongodb

import (
	"context"
	"errors"
	. "github.com/doytowin/goooqo/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const msg = "implement me"

type MongoEntity interface {
	Entity
	Database() string
	Collection() string
}

type QueryBuilder interface {
	BuildFilter() []bson.D
}

type mongoDataAccess[C context.Context, E MongoEntity] struct {
	TransactionManager
	collection *mongo.Collection
	create     func() E
}

func NewMongoDataAccess[E MongoEntity](tm TransactionManager, createEntity func() E) TxDataAccess[E, Query] {
	entity := createEntity()
	client := tm.GetClient().(*mongo.Client)
	collection := client.Database(entity.Database()).Collection(entity.Collection())
	return &mongoDataAccess[context.Context, E]{
		TransactionManager: tm,
		collection:         collection,
		create:             createEntity,
	}
}

func (m *mongoDataAccess[C, E]) Get(c C, id any) (*E, error) {
	objectID, err := primitive.ObjectIDFromHex(id.(string))
	if NoError(err) {
		e := m.create()
		err = m.collection.FindOne(c, bson.M{"_id": objectID}).Decode(&e)
		if NoError(err) {
			return &e, err
		}
	}
	return nil, err
}

func (m *mongoDataAccess[C, E]) Delete(ctx C, id any) (int64, error) {
	panic(msg)
}

func (m *mongoDataAccess[C, E]) Query(ctx C, query Query) ([]E, error) {
	var result []E
	cursor, err := m.collection.Find(ctx, buildFilter(query))
	if NoError(err) {
		err = cursor.All(ctx, &result)
	}
	return result, err
}

func buildFilter(query Query) bson.D {
	if qb, ok := query.(QueryBuilder); ok {
		d := qb.BuildFilter()
		ret := bson.D{}
		if len(d) > 0 {
			ret = bson.D{{"$and", d}}
		}
		return ret
	}
	panic(errors.New("Query object should be type of QueryBuilder"))
}

func (m *mongoDataAccess[C, E]) Count(ctx C, query Query) (int64, error) {
	return m.collection.CountDocuments(ctx, buildFilter(query))
}

func (m *mongoDataAccess[C, E]) DeleteByQuery(ctx C, query Query) (int64, error) {
	panic(msg)
}

func (m *mongoDataAccess[C, E]) Page(ctx C, query Query) (PageList[E], error) {
	var count int64
	data, err := m.Query(ctx, query)
	if NoError(err) {
		count, err = m.Count(ctx, query)
	}
	return PageList[E]{List: data, Total: count}, err
}

func (m *mongoDataAccess[C, E]) Create(ctx C, entity *E) (int64, error) {
	panic(msg)
}

func (m *mongoDataAccess[C, E]) CreateMulti(ctx C, entities []E) (int64, error) {
	panic(msg)
}

func (m *mongoDataAccess[C, E]) Update(ctx C, entity E) (int64, error) {
	panic(msg)
}

func (m *mongoDataAccess[C, E]) Patch(ctx C, entity E) (int64, error) {
	panic(msg)
}

func (m *mongoDataAccess[C, E]) PatchByQuery(ctx C, entity E, query Query) (int64, error) {
	panic(msg)
}
