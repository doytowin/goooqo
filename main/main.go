/*
 * The Clear BSD License
 *
 * Copyright (c) 2024, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

import (
	"context"
	"encoding/json"
	"github.com/doytowin/goooqo"
	"github.com/doytowin/goooqo/core"
	"github.com/doytowin/goooqo/mongodb"
	"github.com/doytowin/goooqo/rdb"
	"github.com/doytowin/goooqo/test"
	"github.com/doytowin/goooqo/web"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"reflect"
	"strings"
)

func init() {
	web.RegisterConverter(reflect.PointerTo(reflect.TypeOf(primitive.NilObjectID)), func(v []string) (any, error) {
		objectID, err := mongodb.ResolveId(v[0])
		return &objectID, err
	})
	web.RegisterConverter(reflect.PointerTo(reflect.TypeOf([]primitive.ObjectID{})), func(params []string) (any, error) {
		if len(params) == 1 {
			params = strings.Split(params[0], ",")
		}
		v := make([]primitive.ObjectID, 0, len(params))
		for _, s := range params {
			objectID, err := mongodb.ResolveId(s)
			if core.NoError(err) {
				v = append(v, objectID)
			}
		}
		return &v, nil
	})
	web.RegisterConverter(reflect.PointerTo(reflect.TypeOf(primitive.M{})), func(v []string) (any, error) {
		d := primitive.M{}
		err := json.Unmarshal([]byte(v[0]), &d)
		return &d, err
	})
}

func main() {
	log.SetLevel(log.DebugLevel)
	db := rdb.Connect("local.properties")
	test.InitDB(db)
	defer rdb.Disconnect(db)
	tm := rdb.NewTransactionManager(db)

	buildUserModule(tm)

	ctx := context.Background()
	var client = mongodb.Connect(ctx, "local.properties")
	defer mongodb.Disconnect(client, ctx)

	mtm := mongodb.NewMongoTransactionManager(client)
	buildInventoryModule(mtm)

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func buildUserModule(tm goooqo.TransactionManager) {
	userDataAccess := rdb.NewTxDataAccess[UserEntity](tm)
	goooqo.BuildRestService[UserEntity, UserQuery]("/user/", userDataAccess)
}

func buildInventoryModule(tm goooqo.TransactionManager) {
	mongoDataAccess := mongodb.NewMongoDataAccess[InventoryEntity](tm)
	goooqo.BuildRestService[InventoryEntity, InventoryQuery]("/inventory/", mongoDataAccess)
}
