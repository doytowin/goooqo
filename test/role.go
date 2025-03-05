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

type RoleEntity struct {
	IntId
	RoleName     *string
	RoleCode     *string
	CreateUserId *int
}

type RoleQuery struct {
	PageQuery
	Id    *int
	Valid *bool

	User *UserQuery `entitypath:"user,role"`
}
