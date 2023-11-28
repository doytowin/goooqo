package test

import goquery "github.com/doytowin/doyto-query-go-sql"

type UserEntity struct {
	Id    int
	Score int
	Memo  string
}

type UserQuery struct {
	goquery.PageQuery
	IdGt     *int
	ScoreLt  *int
	MemoNull bool
}

func (q UserQuery) GetPageQuery() goquery.PageQuery {
	return q.PageQuery
}
