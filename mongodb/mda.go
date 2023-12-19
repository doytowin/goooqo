package mongodb

import (
	"context"
	. "github.com/doytowin/goquery/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const msg = "implement me"

type MongoEntity interface {
	Database() string
	Collection() string
}

type mongoDataAccess[C context.Context, E MongoEntity] struct {
	TransactionManager
	collection *mongo.Collection
	create     func() E
}

func NewMongoDataAccess[E MongoEntity](tm TransactionManager, createEntity func() E) TxDataAccess[E, GoQuery] {
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

func (m *mongoDataAccess[C, E]) Query(ctx C, query GoQuery) ([]E, error) {
	var result []E
	cursor, err := m.collection.Find(ctx, bson.D{{}})
	if NoError(err) {
		err = cursor.All(ctx, &result)
	}
	return result, err
}

func (m *mongoDataAccess[C, E]) Count(ctx C, query GoQuery) (int64, error) {
	return m.collection.CountDocuments(ctx, bson.M{})
}

func (m *mongoDataAccess[C, E]) DeleteByQuery(ctx C, query any) (int64, error) {
	panic(msg)
}

func (m *mongoDataAccess[C, E]) Page(ctx C, query GoQuery) (PageList[E], error) {
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

func (m *mongoDataAccess[C, E]) PatchByQuery(ctx C, entity E, query GoQuery) (int64, error) {
	panic(msg)
}
