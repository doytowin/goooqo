package test

import . "github.com/doytowin/go-query/core"

type UserEntity struct {
	Id    int     `json:"id"`
	Score *int    `json:"score"`
	Memo  *string `json:"memo"`
}

func (u UserEntity) GetTableName() string {
	return "User"
}

type UserQuery struct {
	PageQuery
	IdGt     *int
	IdIn     *[]int
	ScoreLt  *int
	MemoNull bool
	MemoLike *string
	Deleted  *bool
}
