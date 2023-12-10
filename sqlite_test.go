package goquery

import (
	"database/sql"
	. "github.com/doytowin/goquery/util"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestSQLite(t *testing.T) {
	db := initDB()
	defer func() {
		_ = db.Close()
	}()

	userDataAccess := BuildDataAccess[UserEntity](UserEntity{})

	t.Run("Query Entities", func(t *testing.T) {
		userQuery := UserQuery{ScoreLt: PInt(80)}
		users, err := userDataAccess.Query(db, &userQuery)

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
		if !(user.Id == 3 && *user.Score == 55) {
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
		cnt, err := userDataAccess.Delete(tx, 3)
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
		cnt, err := userDataAccess.DeleteByQuery(tx, userQuery)
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
		cnt, err := userDataAccess.Count(db, &userQuery)
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
		page, err := userDataAccess.Page(db, &userQuery)
		if err != nil {
			t.Error("Error", err)
			return
		}
		if !(page.Total == 3 && page.List[0].Id == 2) {
			t.Errorf("Got : %v", page)
		}
	})

	t.Run("Create Entity", func(t *testing.T) {
		tx, err := db.Begin()
		entity := UserEntity{Score: PInt(90), Memo: PStr("Great")}
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

	t.Run("Create Entities", func(t *testing.T) {
		tx, err := db.Begin()
		entities := []UserEntity{{Score: PInt(90), Memo: PStr("Great")}, {Score: PInt(55), Memo: PStr("Bad")}}
		cnt, err := userDataAccess.CreateMulti(tx, entities)
		if err != nil {
			t.Error("Error", err)
			return
		}
		if !(cnt == 2) {
			t.Errorf("\nExpected: %d\nBut got : %d", 2, cnt)
		}
		_ = tx.Rollback()
	})

	t.Run("Create 0 Entity", func(t *testing.T) {
		tx, err := db.Begin()
		var entities []UserEntity
		cnt, err := userDataAccess.CreateMulti(tx, entities)
		if err != nil {
			t.Error("Error", err)
			return
		}
		if cnt != 0 {
			t.Errorf("\nExpected: %d\nBut got : %d", 0, cnt)
		}
		_ = tx.Rollback()
	})

	t.Run("Update Entity", func(t *testing.T) {
		tx, err := db.Begin()
		entity := UserEntity{Id: 2, Score: PInt(90), Memo: PStr("Great")}
		cnt, err := userDataAccess.Update(tx, entity)
		if err != nil {
			t.Error("Error", err)
			return
		}
		userEntity, err := userDataAccess.Get(tx, 2)

		if !(cnt == 1 && *userEntity.Score == 90) {
			t.Errorf("\nExpected: %d\n     Got: %d", 1, cnt)
			t.Errorf("\nExpected: %d\n     Got: %d", 90, *userEntity.Score)
		}
		_ = tx.Rollback()
	})

	t.Run("Patch Entity", func(t *testing.T) {
		tx, err := db.Begin()
		entity := UserEntity{Id: 2, Score: PInt(90)}
		cnt, err := userDataAccess.Patch(tx, entity)
		if err != nil {
			t.Error("Error", err)
			return
		}
		userEntity, err := userDataAccess.Get(tx, 2)

		if !(cnt == 1 && *userEntity.Score == 90 && *userEntity.Memo == "Bad") {
			t.Errorf("\nExpected: %d %d %s\nBut got : %d %d %s",
				2, 90, "Bad", userEntity.Id, *userEntity.Score, *userEntity.Memo)
		}
		_ = tx.Rollback()
	})

	t.Run("Patch Entity By Query", func(t *testing.T) {
		tx, err := db.Begin()
		entity := UserEntity{Memo: PStr("Add Memo")}
		query := UserQuery{MemoNull: true}
		cnt, err := userDataAccess.PatchByQuery(tx, entity, &query)

		if cnt != 1 {
			t.Errorf("\nExpected: %d\nBut got : %d", 1, err)
		}
		count, err := userDataAccess.Count(tx, &query)

		if count != 0 {
			t.Errorf("\nExpected: %d\nBut got : %d", 0, count)
		}
		_ = tx.Rollback()
	})
}

func initDB() *sql.DB {
	db, _ := sql.Open("sqlite3", "./test.db")
	_, _ = db.Exec("drop table User")
	_, _ = db.Exec("create table User(id integer constraint user_pk primary key autoincrement, score int, memo varchar(255))")
	_, _ = db.Exec("insert into User(score, memo) values (85, 'Good'), (40, 'Bad'), (55, null), (62, 'Well')")
	return db
}
