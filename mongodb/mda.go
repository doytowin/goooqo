package mongodb

import (
	"context"
	"errors"
	. "github.com/doytowin/goooqo/core"
	. "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

const msg = "implement me"

type MongoEntity interface {
	Entity
	Database() string
	Collection() string
}

type QueryBuilder interface {
	BuildFilter() []D
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
	ID, err := resolveId(id)
	if NoError(err) {
		e := m.create()
		err = m.collection.FindOne(c, buildIdFilter(ID)).Decode(&e)
		if NoError(err) {
			return &e, err
		}
	}
	return nil, err
}

func (m *mongoDataAccess[C, E]) Delete(ctx C, id any) (int64, error) {
	ID, err := resolveId(id)
	if NoError(err) {
		return unwrap(m.collection.DeleteOne(ctx, buildIdFilter(ID)))
	}
	return 0, err
}

func buildIdFilter(objectID any) M {
	return M{"_id": objectID}
}

func resolveId(id any) (ObjectID, error) {
	switch x := id.(type) {
	case ObjectID:
		return x, nil
	case string:
		return ObjectIDFromHex(x)
	}
	return NilObjectID, errors.New("unknown type of id: " + reflect.TypeOf(id).String())
}

func (m *mongoDataAccess[C, E]) Query(ctx C, query Query) ([]E, error) {
	var result []E
	cursor, err := m.collection.Find(ctx, buildFilter(query))
	if NoError(err) {
		err = cursor.All(ctx, &result)
	}
	return result, err
}

func buildFilter(query Query) D {
	if qb, ok := query.(QueryBuilder); ok {
		d := qb.BuildFilter()
		ret := D{}
		if len(d) > 0 {
			ret = D{{"$and", d}}
		}
		return ret
	}
	panic(errors.New("Query object should be type of QueryBuilder"))
}

func (m *mongoDataAccess[C, E]) Count(ctx C, query Query) (int64, error) {
	return m.collection.CountDocuments(ctx, buildFilter(query))
}

func (m *mongoDataAccess[C, E]) DeleteByQuery(ctx C, query Query) (int64, error) {
	return unwrap(m.collection.DeleteMany(ctx, buildFilter(query)))
}

func unwrap(result *mongo.DeleteResult, err error) (int64, error) {
	if NoError(err) {
		return result.DeletedCount, nil
	}
	return 0, err
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
	result, err := m.collection.InsertOne(ctx, entity)
	if NoError(err) {
		err = (*entity).SetId(entity, result.InsertedID)
	}
	return 0, err
}

func (m *mongoDataAccess[C, E]) CreateMulti(ctx C, entities []E) (int64, error) {
	docs := make([]any, len(entities))
	for i := range entities {
		docs[i] = entities[i]
	}

	result, err := m.collection.InsertMany(ctx, docs)
	if NoError(err) {
		for i, ID := range result.InsertedIDs {
			err = entities[i].SetId(&entities[i], ID)
		}
		return int64(len(result.InsertedIDs)), err
	}
	return 0, err
}

func (m *mongoDataAccess[C, E]) Update(ctx C, entity E) (int64, error) {
	result, err := m.collection.ReplaceOne(ctx, buildIdFilter(entity.GetId()), entity)
	if NoError(err) {
		return result.MatchedCount, err
	}
	return 0, err
}

func (m *mongoDataAccess[C, E]) Patch(ctx C, entity E) (int64, error) {
	panic(msg)
}

func (m *mongoDataAccess[C, E]) PatchByQuery(ctx C, entity E, query Query) (int64, error) {
	panic(msg)
}
