package test

import (
	"database/sql"
	. "github.com/doytowin/goquery"
	. "github.com/doytowin/goquery/util"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestSQLite(t *testing.T) {
	db, _ := sql.Open("sqlite3", "./test.db")
	_, _ = db.Exec("drop table User")
	_, _ = db.Exec("create table User(id integer constraint user_pk primary key autoincrement, score int, memo varchar(255))")
	_, _ = db.Exec("insert into User(score, memo) values (85, 'Good'), (40, 'Bad'), (55, 'Bad'), (62, 'Well')")
	defer func() {
		_ = db.Close()
	}()

	userDataAccess := BuildDataAccess[UserEntity](UserEntity{})

	t.Run("Query Entities", func(t *testing.T) {
		userQuery := UserQuery{ScoreLt: PInt(80)}
		users, err := userDataAccess.Query(db, userQuery)

		if err != nil {
			t.Error("Error", err)
		}
		if !(len(users) == 3 && users[0].Id == 2) {
			t.Errorf("Data is not expected: %v", users)
		}
	})

	t.Run("Query By Id", func(t *testing.T) {
		user, err := userDataAccess.Get(db, 3)

		if err != nil {
			t.Error("Error", err)
		}
		if !(user.Id == 3 && user.Score == 55) {
			t.Errorf("Data is not expected: %v", user)
		}
	})

	t.Run("Query By Non-Existent Id", func(t *testing.T) {
		user, err := userDataAccess.Get(db, -1)

		if err != nil {
			t.Error("Error", err)
		}
		if !userDataAccess.IsZero(user) {
			t.Errorf("Data is not expected: %v", user)
		}
	})

	t.Run("Delete By Id", func(t *testing.T) {
		tx, err := db.Begin()
		cnt, err := userDataAccess.DeleteById(tx, 3)
		if err != nil {
			t.Error("Error", err)
		}
		if cnt != 1 {
			t.Errorf("Delete failed. Deleted: %v", cnt)
		}
		_ = tx.Rollback()
	})

	t.Run("Delete By Query", func(t *testing.T) {
		tx, err := db.Begin()
		userQuery := UserQuery{ScoreLt: PInt(80)}
		cnt, err := userDataAccess.Delete(tx, userQuery)
		if err != nil {
			t.Error("Error", err)
		}
		if cnt != 3 {
			t.Errorf("Delete failed. Deleted: %v", cnt)
		}
		_ = tx.Rollback()
	})

	t.Run("Count By Query", func(t *testing.T) {
		userQuery := UserQuery{ScoreLt: PInt(60)}
		cnt, err := userDataAccess.Count(db, userQuery)
		if err != nil {
			t.Error("Error", err)
		}
		if cnt != 2 {
			t.Errorf("\nExpected: %d\nBut got : %d", 2, cnt)
		}
	})

	t.Run("Page By Query", func(t *testing.T) {
		userQuery := UserQuery{
			PageQuery: PageQuery{PageSize: PInt(2)},
			ScoreLt:   PInt(80),
		}
		page, err := userDataAccess.Page(db, userQuery)
		if err != nil {
			t.Error("Error", err)
			return
		}
		if !(page.Total == 3 && page.Data[0].Id == 2) {
			t.Errorf("Got : %v", page)
		}
	})

	t.Run("Create Entity", func(t *testing.T) {
		tx, err := db.Begin()
		entity := UserEntity{Score: 90, Memo: "Great"}
		id, err := userDataAccess.Create(tx, &entity)
		if err != nil {
			t.Error("Error", err)
			return
		}
		if !(id == 5 && entity.Id == 5) {
			t.Errorf("\nExpected: %d\nBut got : %d", 5, id)
		}
		_ = tx.Rollback()
	})
}
