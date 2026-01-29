/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2026, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package core

type PageQuery struct {
	Page int    `json:"page,omitempty"`
	Size int    `json:"size,omitempty"`
	Sort string `json:"sort,omitempty"`
}

func (pq PageQuery) GetPageNumber() int {
	return Ternary(pq.Page > 1, pq.Page-1, 0)
}

func (pq PageQuery) GetPageSize() int {
	return Ternary(pq.Size > 0, pq.Size, 10)
}

func (pq PageQuery) CalcOffset() int {
	return pq.GetPageNumber() * pq.GetPageSize()
}

func (pq PageQuery) GetSort() string {
	return pq.Sort
}

func (pq PageQuery) NeedPaging() bool {
	return pq.Size > 0 || pq.Page > 0
}
