/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2025, DoytoWin, Inc.
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

	Roles []RoleEntity `entitypath:"user,role" json:"roles,omitempty"`
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

	/**
	id IN (
		SELECT user_id FROM a_user_and_role WHERE role_id IN (
			SELECT id FROM t_role WHERE ...
		)
	)*/
	Role      *RoleQuery `entitypath:"user,role"`
	WithRoles *RoleQuery

	/**
	id IN (
		SELECT user_id FROM a_user_and_role WHERE role_id IN (
			SELECT id FROM t_role WHERE ... INTERSECT
			SELECT role_id FROM a_role_and_perm WHERE perm_id IN (
				SELECT id FROM t_perm WHERE ...
			)
		)
	)*/
	Perm *PermQuery `entitypath:"user,role,perm"`
}

var UserDataAccess TxDataAccess[UserEntity]
