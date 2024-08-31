/*
 * The Clear BSD License
 *
 * Copyright (c) 2024, DoytoWin, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package goooqo

import (
	"database/sql"
	"github.com/doytowin/goooqo/rdb"
	. "github.com/doytowin/goooqo/test"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestBuild(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	t.Run("Export Interface", func(t *testing.T) {
		var db *sql.DB
		tm := rdb.NewTransactionManager(db)
		dataAccess := rdb.NewTxDataAccess[UserEntity](tm)

		BuildRestService[UserEntity, UserQuery]("/user/", dataAccess)
	})
}
