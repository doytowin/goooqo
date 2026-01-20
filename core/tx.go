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
	"fmt"
)

type RollbackError struct {
	Err    error
	Origin error
}

func (e *RollbackError) Error() string {
	return e.Origin.Error()
}

func (e *RollbackError) Unwrap() error {
	return e.Origin
}

func TransactionCallback(tc TransactionContext, callback func(tc TransactionContext) error) (err error) {
	defer func() {
		if re := recover(); re != nil {
			err = RollbackFor(tc, errors.New(fmt.Sprint(re)))
		}
	}()
	err = callback(tc)
	if err == nil {
		err = tc.Commit()
	} else {
		err = RollbackFor(tc, err)
	}
	return
}

func RollbackFor(tc TransactionContext, origin error) error {
	err := tc.Rollback()
	if err == nil {
		return origin
	}
	return &RollbackError{err, origin}
}
