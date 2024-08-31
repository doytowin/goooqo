/*
 * The Clear BSD License
 *
 * Copyright (c) 2024, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package core

import (
	"strconv"
)

type IntId struct {
	Id int `json:"id,omitempty"`
}

func (e IntId) GetId() any {
	return e.Id
}

func NewIntId(id int) IntId {
	return IntId{Id: id}
}

func (e IntId) SetId(self any, id any) (err error) {
	var Id int
	switch x := id.(type) {
	case int:
		Id = x
	case int64:
		Id = int(x)
	case string:
		Id, err = strconv.Atoi(x)
	}
	self.(intIdSetter).setId(Id)
	return
}

type intIdSetter interface {
	setId(id int)
}

func (e *IntId) setId(id int) {
	e.Id = id
}
