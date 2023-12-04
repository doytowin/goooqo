package goquery

import (
	"database/sql"
)

type RestAPI[E comparable, Q GoQuery] interface {
	Page(q Q) (PageList[E], error)
}

type Service[E comparable, Q GoQuery] struct {
	db           *sql.DB
	prefix       string
	dataAccess   DataAccess[E]
	createQuery  func() Q
	createEntity func() E
}

func (s *Service[E, Q]) Page(q Q) (PageList[E], error) {
	return s.dataAccess.Page(s.db, q)
}

func BuildController[E comparable, Q GoQuery](
	prefix string,
	db *sql.DB,
	createEntity func() E,
	createQuery func() Q,
) *Service[E, Q] {
	dataAccess := BuildDataAccess[E](createEntity())
	rc := &Service[E, Q]{
		db:           db,
		prefix:       prefix,
		dataAccess:   dataAccess,
		createQuery:  createQuery,
		createEntity: createEntity,
	}
	return rc
}
