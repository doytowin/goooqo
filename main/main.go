package main

import (
	"github.com/doytowin/go-query"
	"github.com/doytowin/go-query/rdb"
	. "github.com/doytowin/go-query/test"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	log.SetLevel(log.DebugLevel)
	db := rdb.Connect("local.properties")
	InitDB(db)
	defer rdb.Disconnect(db)

	tm := rdb.NewTransactionManager(db)

	buildUserModule(tm)

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
