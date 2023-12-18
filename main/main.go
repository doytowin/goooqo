package main

import (
	"context"
	"github.com/doytowin/goquery"
	"github.com/doytowin/goquery/mongodb"
	"github.com/doytowin/goquery/rdb"
	. "github.com/doytowin/goquery/test"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

func main() {
	log.SetLevel(log.DebugLevel)
	db := rdb.Connect("local.properties")
	InitDB(db)
	defer rdb.Disconnect(db)
	tm := rdb.NewTransactionManager(db)

	buildUserModule(tm)

	ctx := context.Background()
	var client = mongodb.Connect(ctx, "local.properties")
	defer mongodb.Disconnect(client, ctx)

	buildInventoryModule(client)

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func buildUserModule(tm goquery.TransactionManager) {
	createUserEntity := func() UserEntity { return UserEntity{} }
	userDataAccess := rdb.NewTxDataAccess[UserEntity](tm, createUserEntity)
	goquery.BuildRestService[UserEntity, UserQuery](
		"/user/",
		userDataAccess,
		createUserEntity,
		func() UserQuery { return UserQuery{} },
	)
}

func buildInventoryModule(client *mongo.Client) {
	mongoDataAccess := mongodb.BuildMongoDataAccess[context.Context, InventoryEntity](client, func() InventoryEntity { return InventoryEntity{} })
	createInventoryEntity := func() InventoryEntity { return InventoryEntity{} }
	goquery.BuildRestService[InventoryEntity, InventoryQuery](
		"/inventory/", mongoDataAccess, createInventoryEntity,
		func() InventoryQuery { return InventoryQuery{} },
	)
}
