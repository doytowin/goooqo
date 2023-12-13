package goquery

import (
	"database/sql"
	"regexp"
)

type RestAPI[E any, Q GoQuery] interface {
	Page(q Q) (PageList[E], error)
	Get(id any) (*E, error)
}

type Service[E any, Q GoQuery] struct {
	db           *sql.DB
	prefix       string
	dataAccess   DataAccess[E]
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

func BuildController[E any, Q GoQuery](
	prefix string,
	db *sql.DB,
	createEntity func() E,
	createQuery func() Q,
) *Service[E, Q] {
	dataAccess := BuildDataAccess[E](createEntity)
	rc := &Service[E, Q]{
		db:           db,
		prefix:       prefix,
		dataAccess:   dataAccess,
		createQuery:  createQuery,
		createEntity: createEntity,
		idRgx:        regexp.MustCompile(prefix + `(\d+)$`),
	}
	return rc
}
