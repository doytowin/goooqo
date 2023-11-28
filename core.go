package goquery

type GoQuery interface {
	GetPageQuery() PageQuery
}

type DataAccess[E comparable] interface {
	Query(conn connection, query GoQuery) ([]E, error)
	Get(conn connection, id interface{}) (E, error)
	Delete(conn connection, query interface{}) (int64, error)
	DeleteById(conn connection, id interface{}) (int64, error)
	IsZero(entity E) bool
}

func BuildDataAccess[E comparable](entity interface{}) DataAccess[E] {
	e := buildEntityMetadata[E](entity)
	return &e
}
