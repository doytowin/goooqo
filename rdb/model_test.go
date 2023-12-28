package rdb

import (
	"fmt"
	. "github.com/doytowin/goooqo/core"
	"strconv"
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

func (e TestEntity) GetId() any {
	return e.Id
}

func (e TestEntity) SetId(self any, id any) {
	if v, ok := id.(int64); ok {
		self.(*TestEntity).Id = PInt(int(v))
	} else {
		s := fmt.Sprintf("%s", id)
		v, _ := strconv.Atoi(s)
		self.(*TestEntity).Id = &v
	}
}

type TestQuery struct {
	PageQuery
	AccountOr *AccountOr
	Account   *string `condition:"(username = ? OR email = ?)"`
	Deleted   *bool
}
