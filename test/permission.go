/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2026, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package test

import . "github.com/doytowin/goooqo/core"

type PermQuery struct {
	PageQuery
	Id   *int
	Code *string

	// used before INTERSECT in UserQuery.Perm
	RoleQuery *RoleQuery
}
