/*
 * The Clear BSD License
 *
 * Copyright (c) 2024, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package main

import . "github.com/doytowin/goooqo/core"

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
	Cond     *string `condition:"(Score = ? OR Memo = ?)"`
	ScoreLt  *int
	ScoreLt1 *UserQuery `subquery:"select:avg(score),from:User"`
	MemoNull *bool
	MemoLike *string
	Deleted  *bool
}