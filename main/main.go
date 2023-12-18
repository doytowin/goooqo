package main

import (
	"github.com/doytowin/goquery"
	"github.com/doytowin/goquery/core"
	"github.com/doytowin/goquery/rdb"
	. "github.com/doytowin/goquery/test"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	log.SetLevel(log.DebugLevel)
	db := rdb.Connect("local.properties")
	InitDB(db)
	defer rdb.Disconnect(db)

	tm := rdb.NewTransactionManager(db)

	buildUserRestService(tm)

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func buildUserRestService(tm core.TransactionManager) {
	createUserEntity := func() UserEntity { return UserEntity{} }
	userDataAccess := rdb.NewTxDataAccess[UserEntity](tm, createUserEntity)
	goquery.BuildRestService[UserEntity, UserQuery](
		"/user/",
		userDataAccess,
		createUserEntity,
		func() UserQuery { return UserQuery{} },
	)
}
