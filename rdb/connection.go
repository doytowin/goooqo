package rdb

import (
	"database/sql"
	"github.com/joho/godotenv"
	"io"
	"os"
)

func Connect(filenames ...string) *sql.DB {
	err := godotenv.Load(filenames...)
	if err != nil {
		db, _ := sql.Open("sqlite3", "./test.db")
		return db
	}
	driver := os.Getenv("driver")
	if driver == "sqlite3" {
		dataSourceName := os.Getenv("sqlite3_data_source")
		db, err := sql.Open(driver, dataSourceName)
		if err == nil {
			return db
		}
	}
	panic(err)

}
func Disconnect(db io.Closer) {
	_ = db.Close()
}
