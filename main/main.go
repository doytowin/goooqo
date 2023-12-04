package main

import (
	"database/sql"
	"github.com/doytowin/goquery"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./test.db")
	_, _ = db.Exec("drop table User")
	_, _ = db.Exec("create table User(id integer constraint user_pk primary key autoincrement, score int, memo varchar(255))")
	_, _ = db.Exec("insert into User(score, memo) values (85, 'Good'), (40, 'Bad'), (55, 'Bad'), (62, 'Well')")

	return db, err
}

func main() {
	log.SetLevel(log.DebugLevel)
	var db, _ = initDB()
	defer func() {
		_ = db.Close()
	}()

	rc := goquery.BuildController[UserEntity, *UserQuery](
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
