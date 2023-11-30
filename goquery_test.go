package goquery

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestBuild(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	t.Run("Build with Custom Table Name", func(t *testing.T) {
		em := buildEntityMetadata[TestEntity](TestEntity{})
		actual := em.TableName
		expect := "t_user"
		if actual != expect {
			t.Errorf("\nExpected: %s\n     Got: %s", expect, actual)
		}
	})
}
