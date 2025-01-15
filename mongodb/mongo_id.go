/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2025, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

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
	ID, err := ResolveId(id)
	self.(IdSetter).setId(ID)
	return err
}

type IdSetter interface {
	setId(id primitive.ObjectID)
}

func (e *MongoId) setId(id primitive.ObjectID) {
	e.Id = &id
}
