package goquery

import (
	. "github.com/doytowin/goquery/test"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestBuild(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	t.Run("Export Interface", func(t *testing.T) {
		_ = BuildDataAccess[UserEntity](func() UserEntity { return UserEntity{} })
	})
}
