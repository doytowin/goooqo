package mongodb

import (
	"github.com/doytowin/goooqo/core"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
)

type MongoId struct {
	Id *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
}

func NewMongoId(Id *primitive.ObjectID) MongoId {
	return MongoId{Id: Id}
}

func (e MongoId) GetId() any {
	return *e.Id
}

func (e MongoId) SetId(self any, id any) error {
	ID, err := resolveId(id)
	if core.NoError(err) {
		elem := reflect.ValueOf(self).Elem()
		elem.FieldByName("Id").Set(reflect.ValueOf(&ID))
	}
	return err
}
