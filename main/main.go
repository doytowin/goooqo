package main

import (
	"context"
	"database/sql"
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

	buildUserModule(db)

	ctx := context.Background()
	var client = mongodb.Connect(ctx, "local.properties")
	defer mongodb.Disconnect(client, ctx)

	buildInventoryModule(client, ctx)

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func buildUserModule(db *sql.DB) {
	createUserEntity := func() UserEntity { return UserEntity{} }
	userDataAccess := rdb.BuildRelationalDataAccess[UserEntity](createUserEntity)
	userController := goquery.BuildController[rdb.Connection, UserEntity, UserQuery](
		"/user/", db, userDataAccess, createUserEntity,
		func() UserQuery { return UserQuery{} },
	)
	http.Handle(userController.Prefix, userController)
}

func buildInventoryModule(client *mongo.Client, ctx context.Context) {
	mongoDataAccess := mongodb.BuildMongoDataAccess[context.Context, InventoryEntity](client, func() InventoryEntity { return InventoryEntity{} })
	createInventoryEntity := func() InventoryEntity { return InventoryEntity{} }
	roleController := goquery.BuildController[context.Context, InventoryEntity, InventoryQuery](
		"/inventory/", ctx, mongoDataAccess, createInventoryEntity,
		func() InventoryQuery { return InventoryQuery{} },
	)
	http.Handle(roleController.Prefix, roleController)
}
