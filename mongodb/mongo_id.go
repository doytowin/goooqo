package mongodb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	self.(IdSetter).setId(ID)
	return err
}

type IdSetter interface {
	setId(id primitive.ObjectID)
}

func (e *MongoId) setId(id primitive.ObjectID) {
	e.Id = &id
}
