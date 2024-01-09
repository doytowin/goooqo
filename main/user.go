package main

import . "github.com/doytowin/goooqo/core"

type UserEntity struct {
	Int64Id
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
	ScoreLt1 *UserQuery `subquery:"select:avg(score),from:User"`
	MemoNull bool
	MemoLike *string
	Deleted  *bool
}
