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
	. "github.com/doytowin/goooqo/core"
	. "github.com/doytowin/goooqo/test"
	"testing"
)

func TestWeb(t *testing.T) {

	db := Connect("app.properties")
	InitDB(db)
	defer Disconnect(db)

	ctx := context.Background()
	tm := NewTransactionManager(db)

	userDataAccess := NewTxDataAccess[UserEntity](tm)

	t.Run("Should not start tx repeated", func(t *testing.T) {
		tc, _ := tm.StartTransaction(ctx)
		tc2, _ := tm.StartTransaction(tc)
		defer tc2.Rollback()

		if tc2 != tc {
			t.Error("Should not start tx repeated")
		}
	})

	t.Run("Support get parent context from TransactionContext", func(t *testing.T) {
		tc, _ := tm.StartTransaction(ctx)
		defer tc.Rollback()

		if tc.Parent() != ctx {
			t.Error("Should return parent context")
		}
	})

	t.Run("Support save point", func(t *testing.T) {
		tc, _ := tm.StartTransaction(ctx)
		defer tc.Rollback()

		tx := tc.(*rdbTransactionContext).tx
		tx.ExecContext(tc, "DELETE FROM User WHERE id IN (1, 2)")
		NoError(tc.SavePoint("delete0"))
		tx.ExecContext(tc, "DELETE FROM User WHERE id IN (3, 4)")
		NoError(tc.RollbackTo("delete0"))
		entities, _ := userDataAccess.Query(tc, &UserQuery{})
		if !(len(entities) == 2 && entities[0].Id == 3 && entities[1].Id == 4) {
			t.Error("Should support SavePoint: ", entities)
		}
	})
}
