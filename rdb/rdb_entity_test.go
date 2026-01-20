/*
 * The Clear BSD License
 *
 * Copyright (c) 2024-2026, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package rdb

import (
	"fmt"
	"strconv"
	"time"

	. "github.com/doytowin/goooqo/core"
)

type TestEntity struct {
	Id         *int
	Username   *string
	Email      *string
	Mobile     *string
	CreateTime *time.Time
}

func (e TestEntity) GetTableName() string {
	return "t_user"
}

func (e TestEntity) GetId() any {
	return e.Id
}

func (e TestEntity) SetId(self any, id any) (err error) {
	v, ok := id.(int64)
	if !ok {
		s := fmt.Sprintf("%s", id)
		v, err = strconv.ParseInt(s, 10, 64)
	}
	if NoError(err) {
		self.(*TestEntity).Id = P(int(v))
	}
	return
}

type TestQuery struct {
	PageQuery
	Username   *string
	Email      *string
	EmailStart *string
	EmailNull  *bool
	Mobile     *string
	Or         *TestQuery
	And        *TestQuery
	EmailEndOr *[]string
	TestsOr    *[]TestQuery
	Account    *string `condition:"(username = ? OR email = ?)"`
	Deleted    *bool
}
