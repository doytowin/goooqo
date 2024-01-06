package rdb

import "github.com/doytowin/goooqo/core"

type RdbEntity interface {
	core.Entity
	GetTableName() string
}
