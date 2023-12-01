package goquery

import (
	"database/sql"
	"reflect"
)

type RestAPI[E comparable, Q GoQuery] interface {
	Page(q Q) (PageList[E], error)
}

type Service[E comparable, Q GoQuery] struct {
	db         *sql.DB
	prefix     string
	dataAccess DataAccess[E]
	queryType  reflect.Type
	entityType reflect.Type
}

func createModel[T any](t reflect.Type) T {
	if t.Kind() == reflect.Pointer {
		return reflect.New(t.Elem()).Interface().(T)
	}
	return reflect.New(t).Elem().Interface().(T)
}

func (s *Service[E, Q]) Page(q Q) (PageList[E], error) {
	return s.dataAccess.Page(s.db, q)
}

func BuildController[E comparable, Q GoQuery](prefix string, db *sql.DB, e E, q Q) *Service[E, Q] {
	dataAccess := BuildDataAccess[E](e)
	rc := &Service[E, Q]{
		db:         db,
		prefix:     prefix,
		dataAccess: dataAccess,
		queryType:  reflect.TypeOf(q),
		entityType: reflect.TypeOf(e),
	}
	return rc
}
