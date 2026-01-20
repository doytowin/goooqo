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
	"reflect"
	"testing"
	"time"
)

func TestReadValue(t *testing.T) {

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

}

func TestToSnakeCase(t *testing.T) {
	if got := ToSnakeCase("t_user"); got != ("t_user") {
		t.Errorf("ToSnakeCase() = %v, want %v", got, "t_user")
	}
	if got := ToSnakeCase("UserEntity"); got != ("user_entity") {
		t.Errorf("ToSnakeCase() = %v, want %v", got, "user_entity")
	}
}
