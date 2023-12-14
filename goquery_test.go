package goquery

import (
	"database/sql"
	"github.com/doytowin/goquery/rdb"
	. "github.com/doytowin/goquery/test"
	log "github.com/sirupsen/logrus"
	"net/http"
	"testing"
)

func TestBuild(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	t.Run("Export Interface", func(t *testing.T) {
		var db *sql.DB
		dataAccess := rdb.BuildRelationalDataAccess[UserEntity](func() UserEntity { return UserEntity{} })
		rc := BuildController[rdb.Connection, UserEntity, *UserQuery](
			"/user/", db, dataAccess,
			func() UserEntity { return UserEntity{} },
			func() *UserQuery { return &UserQuery{} },
		)
		http.Handle(rc.Prefix, rc)
	})
}
