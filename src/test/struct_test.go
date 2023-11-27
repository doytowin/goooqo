package test

type UserEntity struct {
	Id    int
	Score int
	Memo  string
}

type UserQuery struct {
	IdGt     *int
	ScoreLt  *int
	MemoNull bool
}
