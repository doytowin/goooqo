package main

import (
	"github.com/doytowin/goooqo"
	"github.com/doytowin/goooqo/rdb"
	. "github.com/doytowin/goooqo/test"
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

func buildUserModule(tm goooqo.TransactionManager) {
	userDataAccess := rdb.NewTxDataAccess[UserEntity](tm)
	goooqo.BuildRestService[UserEntity, UserQuery]("/user/", userDataAccess)
}
