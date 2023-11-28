package goquery

type GoQuery interface {
	GetPageQuery() PageQuery
}

type DataAccess[E comparable] interface {
	Get(conn connection, id interface{}) (E, error)
	DeleteById(conn connection, id interface{}) (int64, error)
	Query(conn connection, query GoQuery) ([]E, error)
	Count(conn connection, query GoQuery) (int, error)
	Delete(conn connection, query interface{}) (int64, error)
	IsZero(entity E) bool
}

func BuildDataAccess[E comparable](entity interface{}) DataAccess[E] {
	e := buildEntityMetadata[E](entity)
	return &e
}
