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

import "fmt"

var Dialect DbDialect = &BaseDialect{}

type DbDialect interface {
	BuildPageClause(sql string, offset int, size int) string
}

type BaseDialect struct {
}

func (d *BaseDialect) BuildPageClause(sql string, offset int, size int) string {
	return fmt.Sprintf("%s LIMIT %d OFFSET %d", sql, size, offset)
}

type MySQLDialect struct {
	BaseDialect
}
