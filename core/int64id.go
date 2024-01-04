package core

import (
	"fmt"
	"reflect"
	"strconv"
)

type Int64Id struct {
	Id int64 `json:"id,omitempty"`
}

func (e Int64Id) GetId() any {
	return e.Id
}

func NewIntId(id int64) Int64Id {
	return Int64Id{Id: id}
}

func (e Int64Id) SetId(self any, id any) (err error) {
	Id, ok := id.(int64)
	if !ok {
		s := fmt.Sprintf("%s", id)
		Id, err = strconv.ParseInt(s, 10, 64)
	}
	if NoError(err) {
		elem := reflect.ValueOf(self).Elem()
		elem.FieldByName("Id").SetInt(Id)
	}
	return
}
