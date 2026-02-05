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

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestUtil(t *testing.T) {

	t.Run("P *string", func(t *testing.T) {
		if got := P("t_user"); *got != ("t_user") {
			t.Errorf("P(any) = %v, want %v", *got, "t_user")
		}
	})

	t.Run("Read time.Time", func(t *testing.T) {
		expect := time.Now()
		actual := ReadValue(reflect.ValueOf(expect))
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})

	t.Run("Read *time.Time", func(t *testing.T) {
		expect := time.Now()
		actual := ReadValue(reflect.ValueOf(&expect))
		if actual != expect {
			t.Errorf("\nExpected: %s\nBut got : %s", expect, actual)
		}
	})

	t.Run("ToSnakeCase: t_user", func(t *testing.T) {
		if got := ToSnakeCase("t_user"); got != ("t_user") {
			t.Errorf("ToSnakeCase() = %v, want %v", got, "t_user")
		}
	})
	t.Run("ToSnakeCase: UserEntity", func(t *testing.T) {
		if got := ToSnakeCase("UserEntity"); got != ("user_entity") {
			t.Errorf("ToSnakeCase() = %v, want %v", got, "user_entity")
		}
	})

	t.Run("NoError", func(t *testing.T) {
		if NoError(errors.New("test")) {
			t.Errorf("NoError() should return false")
		}
	})

	t.Run("Capitalize: userEntity", func(t *testing.T) {
		if got := Capitalize("userEntity"); got != ("UserEntity") {
			t.Errorf("Capitalize() = %v, want %v", got, "UserEntity")
		}
	})

	t.Run("Ternary true", func(t *testing.T) {
		expect := 10
		actual := Ternary(true, expect, -1)
		if actual != expect {
			t.Errorf("\nExpected: %d\nBut got : %d", expect, actual)
		}
	})
	t.Run("Ternary false", func(t *testing.T) {
		expect := 10
		actual := Ternary(false, -1, expect)
		if actual != expect {
			t.Errorf("\nExpected: %d\nBut got : %d", expect, actual)
		}
	})
}
