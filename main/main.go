package main

import (
	"context"
	"github.com/doytowin/goooqo"
	"github.com/doytowin/goooqo/mongodb"
	"github.com/doytowin/goooqo/rdb"
	. "github.com/doytowin/goooqo/test"
	"github.com/doytowin/goooqo/web"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"reflect"
)

func main() {
	web.RegisterConverter(reflect.PointerTo(reflect.TypeOf(primitive.NilObjectID)), func(v []string) (any, error) {
		objectID, err := mongodb.ResolveId(v[0])
		return &objectID, err
	})

	log.SetLevel(log.DebugLevel)
	db := rdb.Connect("local.properties")
	InitDB(db)
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
