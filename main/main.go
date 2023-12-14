package main

import (
	"github.com/doytowin/goquery"
	"github.com/doytowin/goquery/rdb"
	. "github.com/doytowin/goquery/test"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	log.SetLevel(log.DebugLevel)
	db := InitDB()
	defer func() {
		_ = db.Close()
	}()

	createUserEntity := func() UserEntity { return UserEntity{} }
	userDataAccess := rdb.BuildRelationalDataAccess[UserEntity](createUserEntity)
	userController := goquery.BuildController[rdb.Connection, UserEntity, *UserQuery](
		"/user/", db, userDataAccess, createUserEntity,
		func() *UserQuery { return &UserQuery{} },
	)
	http.Handle(userController.Prefix, userController)

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
