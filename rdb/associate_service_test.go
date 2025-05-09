/*
 * The Clear BSD License
 *
 * Copyright (c) 2025, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package rdb

import (
	"context"
	. "github.com/doytowin/goooqo/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAssociationSqlBuilder(t *testing.T) {
	builder := NewAssociationSqlBuilder("user", "role")

	t.Run("testSelectK1ColumnByK2Id", func(t *testing.T) {
		assert.Equal(t, "SELECT user_id FROM a_user_and_role WHERE role_id = ?", builder.SelectK1ColumnByK2Id)
	})

	t.Run("testSelectK2ColumnByK1Id", func(t *testing.T) {
		assert.Equal(t, "SELECT role_id FROM a_user_and_role WHERE user_id = ?", builder.SelectK2ColumnByK1Id)
	})

	t.Run("testDeleteByK1", func(t *testing.T) {
		assert.Equal(t, "DELETE FROM a_user_and_role WHERE user_id = ?", builder.DeleteByK1)
	})

	t.Run("testDeleteByK2", func(t *testing.T) {
		assert.Equal(t, "DELETE FROM a_user_and_role WHERE role_id = ?", builder.DeleteByK2)
	})

	keys := []UniqueKey{{K1: 1, K2: 1}, {K1: 2, K2: 3}}

	t.Run("testInsert", func(t *testing.T) {
		sql, args := builder.BuildInsert(keys)
		assert.Equal(t, "INSERT OR IGNORE INTO a_user_and_role (user_id, role_id) VALUES (?, ?), (?, ?)", sql)
		assert.Equal(t, []any{1, 1, 2, 3}, args)
	})

	t.Run("testDelete", func(t *testing.T) {
		sql, args := builder.BuildDelete(keys)
		assert.Equal(t, "DELETE FROM a_user_and_role WHERE (user_id, role_id) IN ((?, ?), (?, ?))", sql)
		assert.Equal(t, []any{1, 1, 2, 3}, args)
	})

	t.Run("testCount", func(t *testing.T) {
		sql, args := builder.BuildCount(keys)
		assert.Equal(t, "SELECT count(*) FROM a_user_and_role WHERE (user_id, role_id) IN ((?, ?), (?, ?))", sql)
		assert.Equal(t, []any{1, 1, 2, 3}, args)
	})

	t.Run("testInsertWithUser", func(t *testing.T) {
		builder.WithCreateUserColumn("create_user_id")
		sql, args := builder.BuildInsertWithUser(keys, 1)
		assert.Equal(t, "INSERT OR IGNORE INTO a_user_and_role (user_id, role_id, create_user_id) VALUES (?, ?, ?), (?, ?, ?)", sql)
		assert.Equal(t, []any{1, 1, 1, 2, 3, 1}, args)
	})
}

func TestJdbcAssociationService(t *testing.T) {
	db := Connect()
	InitDB(db)
	defer Disconnect(db)
	ctx := context.Background()
	tm := NewTransactionManager(db)

	service := NewJdbcAssociationService(tm, "user", "role", "create_user_id")

	t.Run("testAssociate", func(t *testing.T) {
		tc, _ := tm.StartTransaction(ctx)
		defer tc.Rollback()

		ret, err := service.Associate(tc, 1, 20)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), ret)
	})

	t.Run("testQueryK1ByK2", func(t *testing.T) {
		userIds, err := service.QueryK1ByK2(ctx, 2)
		assert.NoError(t, err)
		assert.Equal(t, []int64{1, 4}, userIds)
	})

	t.Run("testQueryK2ByK1", func(t *testing.T) {
		roleIds, err := service.QueryK2ByK1(ctx, 1)
		assert.NoError(t, err)
		assert.Equal(t, []int{1, 2}, roleIds)
	})

	t.Run("testDeleteByK1", func(t *testing.T) {
		tc, _ := tm.StartTransaction(ctx)
		defer tc.Rollback()

		ret, err := service.DeleteByK1(tc, 1)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), ret)
	})

	t.Run("testDeleteByK2", func(t *testing.T) {
		tc, _ := tm.StartTransaction(ctx)
		defer tc.Rollback()

		ret, err := service.DeleteByK2(tc, 1)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), ret)
	})

	t.Run("testReassociateForK1", func(t *testing.T) {
		tc, _ := tm.StartTransaction(ctx)
		defer tc.Rollback()

		ret, err := service.ReassociateForK1(tc, 1, []int{2, 3})
		assert.NoError(t, err)
		assert.Equal(t, int64(2), ret)

		roleIds, err := service.QueryK2ByK1(tc, 1)
		assert.NoError(t, err)
		assert.Equal(t, []int{2, 3}, roleIds)
	})

	t.Run("testReassociateForK1WithEmptyK2", func(t *testing.T) {
		tc, _ := tm.StartTransaction(ctx)
		defer tc.Rollback()

		ret, err := service.ReassociateForK1(tc, 1, []int{})
		assert.NoError(t, err)
		assert.Equal(t, int64(0), ret)

		roleIds, err := service.QueryK2ByK1(tc, 1)
		assert.NoError(t, err)
		assert.Empty(t, roleIds)
	})

	t.Run("testReassociateForK2", func(t *testing.T) {
		tc, _ := tm.StartTransaction(ctx)
		defer tc.Rollback()

		k1List := []int64{1, 2, 3, 4}
		ret, err := service.ReassociateForK2(tc, 1, k1List)
		assert.NoError(t, err)
		assert.Equal(t, int64(4), ret)

		userIds, err := service.QueryK1ByK2(tc, 1)
		assert.NoError(t, err)
		assert.Equal(t, k1List, userIds)
	})

	t.Run("testReassociateForK2WithEmptyK1", func(t *testing.T) {
		tc, _ := tm.StartTransaction(ctx)
		defer tc.Rollback()

		ret, err := service.ReassociateForK2(tc, 1, []int64{})
		assert.NoError(t, err)
		assert.Equal(t, int64(0), ret)

		userIds, err := service.QueryK1ByK2(tc, 1)
		assert.NoError(t, err)
		assert.Empty(t, userIds)
	})

	t.Run("testCount", func(t *testing.T) {
		tc, _ := tm.StartTransaction(ctx)
		defer tc.Rollback()

		keys := []UniqueKey{{K1: 1, K2: 2}, {K1: 1, K2: 3}, {K1: 1, K2: 4}}
		ret, err := service.Count(tc, keys)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), ret)
	})

	t.Run("testDissociate", func(t *testing.T) {
		tc, _ := tm.StartTransaction(ctx)
		defer tc.Rollback()

		ret, err := service.Dissociate(tc, 1, 2)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), ret)

		ret, err = service.Dissociate(tc, 1, 2)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), ret)
	})

	t.Run("testBuildUniqueKeys", func(t *testing.T) {
		keys := service.BuildUniqueKeys(1, []any{2, 2, 3, 4})
		assert.Equal(t, []UniqueKey{{K1: 1, K2: 2}, {K1: 1, K2: 3}, {K1: 1, K2: 4}}, keys)
	})

	t.Run("testExists", func(t *testing.T) {
		exists, err := service.Exists(ctx, 1, 2)
		assert.NoError(t, err)
		assert.True(t, exists)

		exists, err = service.Exists(ctx, 1, 5)
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}
