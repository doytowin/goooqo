package core

import (
	"strconv"
)

type Int64Id struct {
	Id int64 `json:"id,omitempty"`
}

func (e Int64Id) GetId() any {
	return e.Id
}

func NewInt64Id(id int64) Int64Id {
	return Int64Id{Id: id}
}

func (e Int64Id) SetId(self any, id any) (err error) {
	var Id int64
	switch x := id.(type) {
	case int64:
		Id = x
	case string:
		Id, err = strconv.ParseInt(x, 10, 64)
	}
	self.(int64IdSetter).setId(Id)
	return
}

type int64IdSetter interface {
	setId(id int64)
}

func (e *Int64Id) setId(id int64) {
	e.Id = id
}
