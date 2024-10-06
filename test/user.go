/*
 * The Clear BSD License
 *
 * Copyright (c) 2024, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package test

import . "github.com/doytowin/goooqo/core"

type UserEntity struct {
	Int64Id
	Score *int    `json:"score"`
	Memo  *string `json:"memo"`
}

type UserQuery struct {
	PageQuery
	IdGt     *int
	IdIn     *[]int
	IdNotIn  *[]int
	Cond     *string `condition:"(score = ? OR memo = ?)"`
	ScoreLt  *int
	MemoNull *bool
	MemoLike *string
	Deleted  *bool

	ScoreLtAvg *UserQuery `subquery:"select avg(score) from User"`
	ScoreLtAny *UserQuery `subquery:"SELECT score FROM User"`
	ScoreLtAll *UserQuery `subquery:"select score from UserEntity"`
	ScoreGtAvg *UserQuery `select:"avg(score)" from:"UserEntity"`

	ScoreInScoreOfUser    *UserQuery
	ScoreGtAvgScoreOfUser *UserQuery

	Role *RoleQuery `erpath:"user,role"`
	Perm *PermQuery `erpath:"user,role,perm"`
}
