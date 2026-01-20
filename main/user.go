/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2026, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

import . "github.com/doytowin/goooqo/core"

//go:generate gooogen
type UserEntity struct {
	Int64Id
	Score *int    `json:"score"`
	Memo  *string `json:"memo"`
}

func (u UserEntity) GetTableName() string {
	return "User"
}

type UserQuery struct {
	PageQuery
	IdGt     *int
	IdIn     *[]int
	IdNotIn  *[]int
	Cond     *string `condition:"(score = ? OR memo = ?)"`
	ScoreLt  *int
	MemoNull *bool
	Deleted  *bool

	MemoLike       *string
	MemoNotLike    *string
	MemoContain    *string
	MemoNotContain *string
	MemoStart      *string
	MemoNotStart   *string
	MemoEnd        *string
	MemoNotEnd     *string
	MemoRx         *string

	ScoreLtAvg *UserQuery `subquery:"select avg(score) from User"`
	ScoreLtAny *UserQuery `subquery:"SELECT score FROM User"`
	ScoreLtAll *UserQuery `subquery:"select score from User"`
	ScoreGtAvg *UserQuery `select:"avg(score)" from:"User"`
}

var UserDataAccess TxDataAccess[UserEntity]
