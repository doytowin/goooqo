package test

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB(db *sql.DB) {
	_, _ = db.Exec("drop table User")
	_, _ = db.Exec("create table User(id integer constraint user_pk primary key autoincrement, score int, memo varchar(255))")
	_, _ = db.Exec("insert into User(score, memo) values (85, 'Good'), (40, 'Bad'), (55, null), (62, 'Well')")
}
