/*
 * The Clear BSD License
 *
 * Copyright (c) 2026, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package core

import (
	"reflect"
	"testing"

	log "github.com/sirupsen/logrus"
)

type UserEntity struct {
	Int64Id
	Score *int    `json:"score"`
	Memo  *string `json:"memo"`

	Roles   []RoleEntity `entitypath:"user,role" json:"roles,omitempty"`
	Friends []UserEntity `entitypath:""user,friend,friend,friend"" json:"roles,omitempty"`
}

type RoleEntity struct {
	IntId
	RoleName     *string
	RoleCode     *string
	CreateUserId *int

	Users []UserEntity `entitypath:"role,user"`
}

func TestBuildFieldMetas(t *testing.T) {

	log.SetLevel(log.DebugLevel)
	RegisterJoinTable("role", "user", "a_user_and_role")
	RegisterVirtualEntity("friend", "user")

	t.Run("Build FieldMetadata", func(t *testing.T) {
		entityType := reflect.TypeOf(UserEntity{})
		fieldMetas := BuildFieldMetas(entityType)
		actual := len(fieldMetas)
		expect := 5
		if actual != expect {
			t.Errorf("\nExpected: %d\n     Got: %d", expect, actual)
		}
	})
}
