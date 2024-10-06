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

type MenuEntity struct {
	IntId
	ParentId *int    `json:"parentId,omitempty"`
	Name     *string `json:"name,omitempty"`
}

type MenuQuery struct {
	PageQuery
	Id *int

	Parent *MenuQuery `erpath:"menu" localField:"ParentId"`
}
