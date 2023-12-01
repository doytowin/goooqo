package goquery

type UserEntity struct {
	Id    int
	Score *int
	Memo  *string
}

type UserQuery struct {
	PageQuery
	IdGt     *int
	ScoreLt  *int
	MemoNull bool
}

func (q UserQuery) GetPageQuery() *PageQuery {
	return &q.PageQuery
}

type TestEntity struct {
	Id int
}

func (e TestEntity) GetTableName() string {
	return "t_user"
}
