package rdb

import (
	. "github.com/doytowin/goooqo/core"
	"time"
)

type AccountOr struct {
	Username *string
	Email    *string
	Mobile   *string
}

type TestEntity struct {
	Id         *int
	Username   *string
	Email      *string
	Mobile     *string
	CreateTime *time.Time
}

func (e TestEntity) GetTableName() string {
	return "t_user"
}

type TestQuery struct {
	PageQuery
	AccountOr *AccountOr
	Account   *string `condition:"(username = ? OR email = ?)"`
	Deleted   *bool
}
