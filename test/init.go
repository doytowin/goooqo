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
	"github.com/doytowin/goooqo/core"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

func InitDB(db *sql.DB) {
	sqlText := `
drop table if exists a_user_and_role;
drop table if exists t_user;
drop table if exists t_role;

create table t_user(id integer constraint user_pk primary key autoincrement, score integer, memo varchar(255));
create table t_role(id integer constraint role_pk primary key autoincrement, role_name varchar(30), role_code varchar(30), create_user_id integer, valid boolean DEFAULT true);
create table a_user_and_role (user_id int, role_id int, PRIMARY KEY (user_id, role_id));

INSERT INTO t_user(score, memo) VALUES (85, 'Good'), (40, 'Bad'), (55, null), (62, 'Well');
INSERT INTO t_role (role_name, role_code, create_user_id) VALUES ('admin', 'ADMIN', 1);
INSERT INTO t_role (role_name, role_code, create_user_id) VALUES ('vip', 'VIP', 2);
INSERT INTO t_role (role_name, role_code, create_user_id) VALUES ('vip2', 'VIP2', 2);
INSERT INTO t_role (role_name, role_code, create_user_id) VALUES ('vip3', 'VIP3', 0);
INSERT INTO t_role (role_name, role_code, create_user_id) VALUES ('vip4', 'VIP4', null);

INSERT INTO a_user_and_role (user_id, role_id) VALUES (1, 1);
INSERT INTO a_user_and_role (user_id, role_id) VALUES (1, 2);
INSERT INTO a_user_and_role (user_id, role_id) VALUES (3, 1);
INSERT INTO a_user_and_role (user_id, role_id) VALUES (4, 1);
INSERT INTO a_user_and_role (user_id, role_id) VALUES (4, 2);
`
	for _, statement := range strings.Split(sqlText, ";") {
		_, err := db.Exec(statement)
		core.NoError(err)
	}
}
