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

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func InitDB(db *sql.DB) {
	_, _ = db.Exec("drop table t_user")
	_, _ = db.Exec("create table t_user(id integer constraint user_pk primary key autoincrement, score int, memo varchar(255))")
	_, _ = db.Exec("insert into t_user(score, memo) values (85, 'Good'), (40, 'Bad'), (55, null), (62, 'Well')")
}
