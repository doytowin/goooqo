package goquery

type GoQuery interface {
	GetPageQuery() PageQuery
}

type PageList[E comparable] struct {
	Data  []E
	Total int
}

type DataAccess[E comparable] interface {
	Get(conn connection, id any) (E, error)
	DeleteById(conn connection, id any) (int64, error)
	Query(conn connection, query GoQuery) ([]E, error)
	Count(conn connection, query GoQuery) (int, error)
	Delete(conn connection, query any) (int64, error)
	IsZero(entity E) bool
	Page(conn connection, query GoQuery) (PageList[E], error)
}

func BuildDataAccess[E comparable](entity any) DataAccess[E] {
	e := buildEntityMetadata[E](entity)
	return &e
}
