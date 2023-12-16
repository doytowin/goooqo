package rdb

import (
	"database/sql"
	"github.com/doytowin/goquery/core"
	"github.com/joho/godotenv"
	"io"
	"os"
)

func Connect(filenames ...string) *sql.DB {
	var db *sql.DB
	var dataSourceName string
	err := godotenv.Load(filenames...)
	if err != nil {
		db, _ = sql.Open("sqlite3", "./test.db")
		return db
	}
	driver := os.Getenv("driver")
	if driver == "mysql" {
		username := os.Getenv("mysql_username")
		password := os.Getenv("mysql_password")
		url := os.Getenv("mysql_url")
		dataSourceName = username + ":" + password + "@" + url
	} else {
		dataSourceName = os.Getenv("data_source")
	}
	db, err = sql.Open(driver, dataSourceName)
	if err != nil {
		panic(err)
	}
	return db

}

func Disconnect(c io.Closer) {
	core.Close(c)
}
