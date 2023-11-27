package test

import (
	"database/sql"
	q "github.com/doytowin/doyto-query-go-sql/query"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestSQLite(t *testing.T) {
	db, _ := sql.Open("sqlite3", "./test.db")
	_, _ = db.Exec("create table if not exists User(id integer constraint user_pk primary key autoincrement, score int, memo varchar(255))")
	_, _ = db.Exec("insert into User(score, memo) values (85, 'Good'), (40, 'Bad'), (55, 'Bad'), (62, 'Well')")
	defer func() {
		_, _ = db.Exec("drop table User")
		_ = db.Close()
	}()

	em := q.BuildEntityMetadata[UserEntity](UserEntity{})

	t.Run("Query Entities", func(t *testing.T) {
		userQuery := UserQuery{ScoreLt: q.IntPtr(80)}
		users, err := em.Query(db, userQuery)

		if err != nil {
			t.Error("Error", err)
		}
		if !(len(users) == 3 && users[0].Id == 2) {
			t.Errorf("Data is not expected: %v", users)
		}
	})
}
