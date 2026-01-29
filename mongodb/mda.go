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
	"errors"
	"reflect"
	"strings"

	. "github.com/doytowin/goooqo/core"
	log "github.com/sirupsen/logrus"
	. "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MID = "_id"

type MongoEntity interface {
	Entity
	Database() string
	Collection() string
}

type QueryBuilder interface {
	BuildFilter(connector string) D
}

func CombineConditions(connector string, a A) D {
	if len(a) > 1 {
		return D{{connector, a}}
	} else if len(a) == 1 {
		return a[0].(D)
	} else {
		return D{}
	}
}

type mongoDataAccess[E MongoEntity] struct {
	TransactionManager
	collection *mongo.Collection
}

func NewMongoDataAccess[E MongoEntity](tm TransactionManager) TxDataAccess[E] {
	entity := *new(E)
	client := tm.GetClient().(*mongo.Client)
	collection := client.Database(entity.Database()).Collection(entity.Collection())
	entityType := reflect.TypeOf(entity)
	createIndex(entityType, collection)
	return &mongoDataAccess[E]{
		TransactionManager: tm,
		collection:         collection,
	}
}

func createIndex(entityType reflect.Type, collection *mongo.Collection) {
	indexModel := createIndexModel(entityType)
	if len(indexModel) > 0 {
		indexName, err := collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{Keys: indexModel})
		if NoError(err) {
			log.Infof("Index created: %s:%s", collection.Name(), indexName)
		}
	}
}

func createIndexModel(entityType reflect.Type) D {
	ret := make(D, 0, 4)
	for i := 0; i < entityType.NumField(); i++ {
		field := entityType.Field(i)
		if columnTag, ok := field.Tag.Lookup("column"); ok {
			column, indexed := resolveColumnTag(columnTag, field.Name)
			if indexed {
				ret = append(ret, E{column, "text"})
			}
		}
	}
	return ret
}

func resolveColumnTag(columnTag string, fieldName string) (string, bool) {
	values := strings.Split(columnTag, ",")
	column := values[0]
	if column == "" {
		column = ConvertToColumnCase(fieldName)
	}
	indexed := strings.Contains(columnTag, ",index")
	return column, indexed
}

func (m *mongoDataAccess[E]) Get(ctx context.Context, id any) (*E, error) {
	ID, err := ResolveId(id)
	if NoError(err) {
		e := *new(E)
		err = m.collection.FindOne(ctx, buildIdFilter(ID)).Decode(&e)
		if NoError(err) {
			return &e, err
		}
	}
	return nil, err
}

func (m *mongoDataAccess[E]) Delete(ctx context.Context, id any) (int64, error) {
	ID, err := ResolveId(id)
	if NoError(err) {
		return unwrap(m.collection.DeleteOne(ctx, buildIdFilter(ID)))
	}
	return 0, err
}

func buildIdFilter(objectID any) D {
	return D{{MID, objectID}}
}

func ResolveId(id any) (ObjectID, error) {
	switch x := id.(type) {
	case ObjectID:
		return x, nil
	case string:
		return ObjectIDFromHex(x)
	}
	return NilObjectID, errors.New("unknown type of id: " + reflect.TypeOf(id).String())
}

func (m *mongoDataAccess[E]) Query(ctx context.Context, query Query) ([]E, error) {
	filter := buildFilter(query)
	return m.doQuery(ctx, query, filter)
}

func (m *mongoDataAccess[E]) doQuery(ctx context.Context, query Query, filter D) ([]E, error) {
	result := make([]E, 0, query.GetPageSize())
	cursor, err := m.collection.Find(ctx, filter, buildPageOpt(query))
	if NoError(err) {
		err = cursor.All(ctx, &result)
	}
	return result, err
}

func buildPageOpt(query Query) *options.FindOptions {
	pageOpt := &options.FindOptions{}
	if query.NeedPaging() {
		pageOpt.Limit = PInt64(query.GetPageSize())
		pageOpt.Skip = PInt64(query.CalcOffset())
	}
	if query.GetSort() != "" {
		pageOpt.Sort = buildSort(query.GetSort())
	}
	return pageOpt
}

func PInt64(i int) *int64 {
	i64 := int64(i)
	return &i64
}

func buildFilter(query Query) D {
	if qb, ok := query.(QueryBuilder); ok {
		return qb.BuildFilter("$and")
	}
	panic(errors.New("Query object should be type of QueryBuilder"))
}

func (m *mongoDataAccess[E]) Count(ctx context.Context, query Query) (int64, error) {
	filter := buildFilter(query)
	return m.doCount(ctx, filter)
}

func (m *mongoDataAccess[E]) doCount(ctx context.Context, filter D) (int64, error) {
	return m.collection.CountDocuments(ctx, filter)
}

func (m *mongoDataAccess[E]) DeleteByQuery(ctx context.Context, query Query) (int64, error) {
	filter := buildFilter(query)
	if query.NeedPaging() {
		IDs, err := m.doQueryIds(ctx, query, filter)
		if HasError(err) {
			return 0, err
		}
		filter = D{{MID, D{{"$in", IDs}}}}
	}
	return unwrap(m.collection.DeleteMany(ctx, filter))
}

func (m *mongoDataAccess[E]) QueryIds(ctx context.Context, query Query) ([]any, error) {
	filter := buildFilter(query)
	return m.doQueryIds(ctx, query, filter)
}

func (m *mongoDataAccess[E]) doQueryIds(ctx context.Context, query Query, filter any) ([]any, error) {
	pageOpt := buildPageOpt(query).SetProjection(M{MID: 1})
	cursor, err := m.collection.Find(ctx, filter, pageOpt)
	if NoError(err) {
		var result []M
		err = cursor.All(ctx, &result)
		if NoError(err) {
			IDs := make([]any, 0, len(result))
			for _, e := range result {
				IDs = append(IDs, e[MID])
			}
			return IDs, err
		}
	}
	return nil, err
}

func unwrap(result *mongo.DeleteResult, err error) (int64, error) {
	if NoError(err) {
		return result.DeletedCount, nil
	}
	return 0, err
}

func (m *mongoDataAccess[E]) Page(ctx context.Context, query Query) (PageList[E], error) {
	var count int64
	filter := buildFilter(query)
	data, err := m.doQuery(ctx, query, filter)
	if NoError(err) {
		count, err = m.doCount(ctx, filter)
	}
	return PageList[E]{List: data, Total: count}, err
}

func (m *mongoDataAccess[E]) Create(ctx context.Context, entity *E) (int64, error) {
	result, err := m.collection.InsertOne(ctx, entity)
	if NoError(err) {
		err = (*entity).SetId(entity, result.InsertedID)
	}
	return 0, err
}

func (m *mongoDataAccess[E]) CreateMulti(ctx context.Context, entities []E) (int64, error) {
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

func (m *mongoDataAccess[E]) Update(ctx context.Context, entity E) (int64, error) {
	result, err := m.collection.ReplaceOne(ctx, buildIdFilter(entity.GetId()), entity)
	if NoError(err) {
		return result.MatchedCount, err
	}
	return 0, err
}

func (m *mongoDataAccess[E]) Patch(ctx context.Context, entity Entity) (int64, error) {
	doc := buildPatch(entity)
	idFilter := buildIdFilter(entity.GetId())
	return unwrapPatch(m.collection.UpdateMany(ctx, idFilter, doc))
}

func buildPatch(entity any) M {
	dst := M{}
	flattenDoc(dst, "", entity)
	return M{"$set": dst}
}

func flattenDoc(dst M, path string, value any) {
	vType := reflect.TypeOf(value)
	if vType.Kind() == reflect.Struct {
		if path != "" {
			path += "."
		}
		rv := reflect.ValueOf(value)
		for i := 0; i < vType.NumField(); i++ {
			value = ReadValue(rv.Field(i))
			if value != nil {
				name := readFieldName(vType.Field(i))
				flattenDoc(dst, path+name, value)
			}
		}
	} else if path != MID {
		dst[path] = value
	}
}

func readFieldName(field reflect.StructField) string {
	if bsonTag, ok := field.Tag.Lookup("bson"); ok {
		if strings.Index(bsonTag, "inline") > 0 {
			return ""
		}
		if name := strings.Split(bsonTag, ",")[0]; name != "" {
			return name
		}
	}
	return ConvertToColumnCase(field.Name)
}

func (m *mongoDataAccess[E]) PatchByQuery(ctx context.Context, entity E, query Query) (int64, error) {
	doc := buildPatch(entity)
	filter := buildFilter(query)
	if query.NeedPaging() {
		IDs, err := m.doQueryIds(ctx, query, filter)
		if HasError(err) {
			return 0, err
		}
		filter = D{{MID, D{{"$in", IDs}}}}
	}
	return unwrapPatch(m.collection.UpdateMany(ctx, filter, doc))
}

func unwrapPatch(result *mongo.UpdateResult, err error) (int64, error) {
	if NoError(err) {
		return result.MatchedCount, err
	}
	return 0, err
}
