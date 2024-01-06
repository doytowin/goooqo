package mongodb

import (
	"context"
	"errors"
	. "github.com/doytowin/goooqo/core"
	. "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"strings"
)

const MID = "_id"

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
}

func NewMongoDataAccess[E MongoEntity](tm TransactionManager) TxDataAccess[E, Query] {
	entity := *new(E)
	client := tm.GetClient().(*mongo.Client)
	collection := client.Database(entity.Database()).Collection(entity.Collection())
	return &mongoDataAccess[context.Context, E]{
		TransactionManager: tm,
		collection:         collection,
	}
}

func (m *mongoDataAccess[C, E]) Get(c C, id any) (*E, error) {
	ID, err := resolveId(id)
	if NoError(err) {
		e := *new(E)
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

func buildIdFilter(objectID any) D {
	return D{{MID, objectID}}
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
	result := make([]E, 0, query.GetPageSize())
	cursor, err := m.collection.Find(ctx, buildFilter(query), buildPageOpt(query))
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
	if query.GetSort() != nil {
		pageOpt.Sort = buildSort(*query.GetSort())
	}
	return pageOpt
}

func PInt64(i int) *int64 {
	i64 := int64(i)
	return &i64
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

func (m *mongoDataAccess[C, E]) QueryIds(ctx C, query Query) ([]any, error) {
	filter := buildFilter(query)
	return m.doQueryIds(ctx, query, filter)
}

func (m *mongoDataAccess[C, E]) doQueryIds(ctx C, query Query, filter any) ([]any, error) {
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

func (m *mongoDataAccess[C, E]) PatchByQuery(ctx C, entity E, query Query) (int64, error) {
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
