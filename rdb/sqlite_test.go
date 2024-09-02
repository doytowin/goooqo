/*
 * The Clear BSD License
 *
 * Copyright (c) 2024, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package rdb

import (
	"context"
	"errors"
	. "github.com/doytowin/goooqo/core"
	. "github.com/doytowin/goooqo/test"
	"testing"
)

func TestSQLite(t *testing.T) {
	db := Connect()
	InitDB(db)
	defer Disconnect(db)
	ctx := context.Background()
	tm := NewTransactionManager(db)

	userDataAccess := newRelationalDataAccess[UserEntity](tm)

	t.Run("Query Entities", func(t *testing.T) {
		userQuery := UserQuery{ScoreLt: P(80)}
		users, err := userDataAccess.Query(ctx, &userQuery)

		if err != nil {
			t.Error("Error", err)
		}
		if !(len(users) == 3 && users[0].Id == 2) {
			t.Errorf("Data is not expected: %v", users)
		}
	})

	t.Run("Query By Id", func(t *testing.T) {
		user, err := userDataAccess.Get(ctx, 3)

		if err != nil {
			t.Error("Error", err)
		}
		if !(user.Id == 3 && *user.Score == 55) {
			t.Errorf("Data is not expected: %v", user)
		}
	})

	t.Run("Query By Non-Existent Id", func(t *testing.T) {
		user, err := userDataAccess.Get(ctx, -1)

		if err != nil {
			t.Error("Error", err)
		}
		if user != nil {
			t.Errorf("Data is not expected: %v", &user)
		}
	})

	t.Run("Delete By Id", func(t *testing.T) {
		tc, err := tm.StartTransaction(ctx)
		cnt, err := userDataAccess.Delete(tc, 3)
		if err != nil {
			t.Error("Error", err)
		}
		if cnt != 1 {
			t.Errorf("Delete failed. Deleted: %v", cnt)
		}
		_ = tc.Rollback()
	})

	t.Run("Delete By Query", func(t *testing.T) {
		tc, err := tm.StartTransaction(ctx)
		userQuery := UserQuery{ScoreLt: P(80)}
		cnt, err := userDataAccess.DeleteByQuery(tc, userQuery)
		if err != nil {
			t.Error("Error", err)
		}
		if cnt != 3 {
			t.Errorf("Delete failed. Deleted: %v", cnt)
		}
		_ = tc.Rollback()
	})

	t.Run("Rollback: Delete By Query", func(t *testing.T) {
		var userQuery Query
		err := tm.SubmitTransaction(ctx, func(tc TransactionContext) error {
			userQuery = UserQuery{ScoreLt: P(80)}
			userDataAccess.DeleteByQuery(tc, userQuery)
			return errors.New("Rollback")
		})
		count, err := userDataAccess.Count(ctx, userQuery)
		if count != 3 {
			t.Error("Error: rollback failed: ", err)
		}
	})

	t.Run("Rollback for panic: Delete By Query", func(t *testing.T) {
		var userQuery Query
		err := tm.SubmitTransaction(ctx, func(tc TransactionContext) error {
			userQuery = UserQuery{ScoreLt: P(80)}
			userDataAccess.DeleteByQuery(tc, userQuery)
			panic("Recover")
		})
		if err.Error() != "Recover" {
			t.Error("Error", err)
		}
		count, err := userDataAccess.Count(ctx, userQuery)
		if count != 3 {
			t.Error("Error: rollback failed: ", err)
		}
	})

	t.Run("Count By Query", func(t *testing.T) {
		userQuery := UserQuery{ScoreLt: P(60)}
		cnt, err := userDataAccess.Count(ctx, &userQuery)
		if err != nil {
			t.Error("Error", err)
		}
		if cnt != 2 {
			t.Errorf("\nExpected: %d\nBut got : %d", 2, cnt)
		}
	})

	t.Run("Page By Query", func(t *testing.T) {
		userQuery := UserQuery{
			PageQuery: PageQuery{PageSize: P(2)},
			ScoreLt:   P(80),
		}
		page, err := userDataAccess.Page(ctx, &userQuery)
		if err != nil {
			t.Error("Error", err)
			return
		}
		if !(page.Total == 3 && page.List[0].Id == 2) {
			t.Errorf("Got : %v", page)
		}
	})

	t.Run("Create Entity", func(t *testing.T) {
		tc, err := tm.StartTransaction(ctx)
		entity := UserEntity{Score: P(90), Memo: P("Great")}
		id, err := userDataAccess.Create(tc, &entity)
		if err != nil {
			t.Error("Error", err)
			return
		}
		if !(id == 5 && entity.Id == 5) {
			t.Errorf("\nExpected: %d\nBut got : %d", 5, id)
		}
		_ = tc.Rollback()
	})

	t.Run("Create Entities", func(t *testing.T) {
		tc, err := tm.StartTransaction(ctx)
		entities := []UserEntity{{Score: P(90), Memo: P("Great")}, {Score: P(55), Memo: P("Bad")}}
		cnt, err := userDataAccess.CreateMulti(tc, entities)
		if err != nil {
			t.Error("Error", err)
			return
		}
		if !(cnt == 2) {
			t.Errorf("\nExpected: %d\nBut got : %d", 2, cnt)
		}
		_ = tc.Rollback()
	})

	t.Run("Create 0 Entity", func(t *testing.T) {
		tc, err := tm.StartTransaction(ctx)
		var entities []UserEntity
		cnt, err := userDataAccess.CreateMulti(tc, entities)
		if err != nil {
			t.Error("Error", err)
			return
		}
		if cnt != 0 {
			t.Errorf("\nExpected: %d\nBut got : %d", 0, cnt)
		}
		_ = tc.Rollback()
	})

	t.Run("Update Entity", func(t *testing.T) {
		tc, err := tm.StartTransaction(ctx)
		entity := UserEntity{Score: P(90), Memo: P("Great")}
		entity.Id = 2
		cnt, err := userDataAccess.Update(tc, entity)
		if err != nil {
			t.Error("Error", err)
			return
		}
		userEntity, err := userDataAccess.Get(tc, 2)

		if !(cnt == 1 && *userEntity.Score == 90) {
			t.Errorf("\nExpected: %d\n     Got: %d", 1, cnt)
			t.Errorf("\nExpected: %d\n     Got: %d", 90, *userEntity.Score)
		}
		_ = tc.Rollback()
	})

	t.Run("Patch Entity", func(t *testing.T) {
		tc, err := tm.StartTransaction(ctx)
		entity := UserEntity{Score: P(90)}
		entity.Id = 2
		cnt, err := userDataAccess.Patch(ctx, entity)
		if err != nil {
			t.Error("Error", err)
			return
		}
		userEntity, err := userDataAccess.Get(ctx, 2)

		if !(cnt == 1 && *userEntity.Score == 90 && *userEntity.Memo == "Bad") {
			t.Errorf("\nExpected: %d %d %s\nBut got : %d %d %s",
				2, 90, "Bad", userEntity.Id, *userEntity.Score, *userEntity.Memo)
		}
		_ = tc.Rollback()
	})

	t.Run("Patch Entity By Query", func(t *testing.T) {
		tc, err := tm.StartTransaction(ctx)
		entity := UserEntity{Memo: P("Add Memo")}
		query := UserQuery{MemoNull: P(true)}
		cnt, err := userDataAccess.PatchByQuery(tc, entity, &query)

		if cnt != 1 {
			t.Errorf("\nExpected: %d\nBut got : %d", 1, err)
		}
		count, err := userDataAccess.Count(tc, &query)

		if count != 0 {
			t.Errorf("\nExpected: %d\nBut got : %d", 0, count)
		}
		_ = tc.Rollback()
	})
}
