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
	"fmt"
	. "github.com/doytowin/goooqo/core"
)

type UniqueKey struct {
	K1 any
	K2 any
}

type AssociationSqlBuilder struct {
	tableName        string
	k1Column         string
	k2Column         string
	createUserColumn string

	// SQL templates
	SelectK1ColumnByK2Id string
	SelectK2ColumnByK1Id string
	DeleteByK1           string
	DeleteByK2           string
	InsertSQL            string
	DeleteSQL            string
	CountSQL             string
}

func NewAssociationSqlBuilder(e1, e2 string) *AssociationSqlBuilder {
	tableName := fmt.Sprintf("a_%s_and_%s", e1, e2)
	k1Column := fmt.Sprintf("%s_id", e1)
	k2Column := fmt.Sprintf("%s_id", e2)

	return &AssociationSqlBuilder{
		tableName: tableName,
		k1Column:  k1Column,
		k2Column:  k2Column,

		SelectK1ColumnByK2Id: "SELECT " + k1Column + " FROM " + tableName + " WHERE " + k2Column + " = ?",
		SelectK2ColumnByK1Id: "SELECT " + k2Column + " FROM " + tableName + " WHERE " + k1Column + " = ?",
		DeleteByK1:           "DELETE FROM " + tableName + " WHERE " + k1Column + " = ?",
		DeleteByK2:           "DELETE FROM " + tableName + " WHERE " + k2Column + " = ?",
		InsertSQL:            "INSERT OR IGNORE INTO " + tableName + " (" + k1Column + ", " + k2Column + ") VALUES ",
		DeleteSQL:            "DELETE FROM " + tableName + " WHERE (" + k1Column + ", " + k2Column + ") IN (",
		CountSQL:             "SELECT count(*) FROM " + tableName + " WHERE (" + k1Column + ", " + k2Column + ") IN (",
	}
}

func (b *AssociationSqlBuilder) BuildInsert(keys []UniqueKey) (string, []any) {
	sql := b.InsertSQL
	args := make([]any, 0)
	for i, key := range keys {
		if i > 0 {
			sql += ", "
		}
		sql += "(?, ?)"
		args = append(args, key.K1, key.K2)
	}
	return sql, args
}

func (b *AssociationSqlBuilder) BuildDelete(keys []UniqueKey) (string, []any) {
	sql := b.DeleteSQL
	args := make([]any, 0)
	for i, key := range keys {
		if i > 0 {
			sql += ", "
		}
		sql += "(?, ?)"
		args = append(args, key.K1, key.K2)
	}
	sql += ")"
	return sql, args
}

func (b *AssociationSqlBuilder) BuildCount(keys []UniqueKey) (string, []any) {
	sql := b.CountSQL
	args := make([]any, 0)
	for i, key := range keys {
		if i > 0 {
			sql += ", "
		}
		sql += "(?, ?)"
		args = append(args, key.K1, key.K2)
	}
	sql += ")"
	return sql, args
}

func (b *AssociationSqlBuilder) WithCreateUserColumn(column string) {
	b.createUserColumn = column
}

func (b *AssociationSqlBuilder) BuildInsertWithUser(keys []UniqueKey, userId int) (string, []any) {
	sql := "INSERT OR IGNORE INTO " + b.tableName + " (" + b.k1Column + ", " + b.k2Column + ", " + b.createUserColumn + ") VALUES "
	args := make([]any, 0)
	for i, key := range keys {
		if i > 0 {
			sql += ", "
		}
		sql += "(?, ?, ?)"
		args = append(args, key.K1, key.K2, userId)
	}
	return sql, args
}

type JdbcAssociationService struct {
	builder AssociationSqlBuilder
	tm      TransactionManager
	conn    Connection
}

func NewJdbcAssociationService(tm TransactionManager, e1, e2, createUserColumn string) *JdbcAssociationService {
	builder := *NewAssociationSqlBuilder(e1, e2)
	builder.WithCreateUserColumn(createUserColumn)
	return &JdbcAssociationService{
		builder: builder,
		tm:      tm,
		conn:    tm.GetClient().(Connection),
	}
}

func (s *JdbcAssociationService) getConn(ctx context.Context) Connection {
	if tc, ok := ctx.(*rdbTransactionContext); ok {
		return tc.tx
	}
	return s.conn
}

func (s *JdbcAssociationService) Associate(ctx context.Context, k1 any, k2 any) (int64, error) {
	sql := s.builder.InsertSQL + "(?, ?)"
	return parse(s.getConn(ctx).ExecContext(ctx, sql, k1, k2))
}

func (s *JdbcAssociationService) QueryK1ByK2(ctx context.Context, k2 int) ([]int64, error) {
	rows, err := s.getConn(ctx).QueryContext(ctx, s.builder.SelectK1ColumnByK2Id, k2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var k1List []int64
	for rows.Next() {
		var k1 int64
		if err = rows.Scan(&k1); err != nil {
			return nil, err
		}
		k1List = append(k1List, k1)
	}
	return k1List, nil
}

func (s *JdbcAssociationService) QueryK2ByK1(ctx context.Context, k1 int64) ([]int, error) {
	rows, err := s.getConn(ctx).QueryContext(ctx, s.builder.SelectK2ColumnByK1Id, k1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var k2List []int
	for rows.Next() {
		var k2 int
		if err := rows.Scan(&k2); err != nil {
			return nil, err
		}
		k2List = append(k2List, k2)
	}
	return k2List, nil
}

func (s *JdbcAssociationService) DeleteByK1(ctx context.Context, k1 any) (int64, error) {
	return parse(s.getConn(ctx).ExecContext(ctx, s.builder.DeleteByK1, k1))
}

func (s *JdbcAssociationService) DeleteByK2(ctx context.Context, k2 any) (int64, error) {
	return parse(s.getConn(ctx).ExecContext(ctx, s.builder.DeleteByK2, k2))
}

func (s *JdbcAssociationService) ReassociateForK1(ctx context.Context, k1 int64, k2List []int) (int64, error) {
	if cnt, err := s.DeleteByK1(ctx, k1); err != nil {
		return cnt, err
	}

	if len(k2List) == 0 {
		return 0, nil
	}
	sql := s.builder.InsertSQL
	args := make([]any, 0)
	for i, k2 := range k2List {
		if i > 0 {
			sql += ", "
		}
		sql += "(?, ?)"
		args = append(args, k1, k2)
	}
	return parse(s.getConn(ctx).ExecContext(ctx, sql, args...))
}

func (s *JdbcAssociationService) ReassociateForK2(ctx context.Context, k2 any, k1List []int64) (int64, error) {
	if cnt, err := s.DeleteByK2(ctx, k2); err != nil {
		return cnt, err
	}

	if len(k1List) == 0 {
		return 0, nil
	}
	sql := s.builder.InsertSQL
	args := make([]any, 0)
	for i, k1 := range k1List {
		if i > 0 {
			sql += ", "
		}
		sql += "(?, ?)"
		args = append(args, k1, k2)
	}
	return parse(s.getConn(ctx).ExecContext(ctx, sql, args...))
}

func (s *JdbcAssociationService) Count(ctx context.Context, keys []UniqueKey) (int64, error) {
	sql, args := s.builder.BuildCount(keys)
	var count int64
	err := s.getConn(ctx).QueryRowContext(ctx, sql, args...).Scan(&count)
	return count, err
}

func (s *JdbcAssociationService) Dissociate(ctx context.Context, k1 any, k2 any) (int64, error) {
	sql := "DELETE FROM " + s.builder.tableName + " WHERE " + s.builder.k1Column + " = ? AND " + s.builder.k2Column + " = ?"
	return parse(s.getConn(ctx).ExecContext(ctx, sql, k1, k2))
}

func (s *JdbcAssociationService) BuildUniqueKeys(k1 any, k2List []any) []UniqueKey {
	keys := make([]UniqueKey, 0)
	seen := make(map[UniqueKey]struct{})

	for _, k2 := range k2List {
		key := UniqueKey{K1: k1, K2: k2}
		if _, exists := seen[key]; !exists {
			keys = append(keys, key)
			seen[key] = struct{}{}
		}
	}
	return keys
}

func (s *JdbcAssociationService) Exists(ctx context.Context, k1 any, k2 any) (bool, error) {
	sql := "SELECT COUNT(*) FROM " + s.builder.tableName + " WHERE " + s.builder.k1Column + " = ? AND " + s.builder.k2Column + " = ?"
	var count int64
	err := s.getConn(ctx).QueryRowContext(ctx, sql, k1, k2).Scan(&count)
	return count > 0, err
}
