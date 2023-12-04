package goquery

type UserEntity struct {
	Id    int
	Score *int
	Memo  *string
}

type UserQuery struct {
	PageQuery
	IdGt     *int
	IdIn     *[]int
	ScoreLt  *int
	MemoNull bool
	MemoLike *string
}

func (q *UserQuery) GetPageQuery() *PageQuery {
	return &q.PageQuery
}

type TestEntity struct {
	Id int
}

func (e TestEntity) GetTableName() string {
	return "t_user"
}
