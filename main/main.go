package main

import (
	. "github.com/doytowin/goquery/test"
	"github.com/doytowin/goquery/web"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	log.SetLevel(log.DebugLevel)
	db := InitDB()
	defer func() {
		_ = db.Close()
	}()

	rc := web.BuildController[UserEntity, *UserQuery](
		"/user/", db,
		func() UserEntity { return UserEntity{} },
		func() *UserQuery { return &UserQuery{} },
	)
	http.Handle("/user/", rc)

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
