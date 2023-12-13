package web

import (
	"database/sql"
	. "github.com/doytowin/goquery/core"
	"github.com/doytowin/goquery/rdb"
	"regexp"
)

type RestAPI[E any, Q GoQuery] interface {
	Page(q Q) (PageList[E], error)
	Get(id any) (*E, error)
}

type Service[E any, Q GoQuery] struct {
	db           *sql.DB
	prefix       string
	dataAccess   DataAccess[Connection, E]
	createQuery  func() Q
	createEntity func() E
	idRgx        *regexp.Regexp
}

func (s *Service[E, Q]) Page(q Q) (PageList[E], error) {
	return s.dataAccess.Page(s.db, q)
}

func (s *Service[E, Q]) Get(id any) (*E, error) {
	return s.dataAccess.Get(s.db, id)
}

func BuildService[E any, Q GoQuery](
	prefix string,
	db *sql.DB,
	createEntity func() E,
	createQuery func() Q,
) *Service[E, Q] {
	rc := &Service[E, Q]{
		db:           db,
		prefix:       prefix,
		createQuery:  createQuery,
		createEntity: createEntity,
		idRgx:        regexp.MustCompile(prefix + `(\d+)$`),
	}
	rc.dataAccess = rdb.BuildDataAccess[E](createEntity)
	return rc
}
