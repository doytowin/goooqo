[![Sonar Stats](https://sonarcloud.io/api/project_badges/measure?project=win.doyto.goooqo&metric=alert_status)](https://sonarcloud.io/dashboard?id=win.doyto.goooqo)
[![Code Lines](https://sonarcloud.io/api/project_badges/measure?project=win.doyto.goooqo&metric=ncloc)](https://sonarcloud.io/component_measures?id=win.doyto.goooqo&metric=ncloc)
[![Coverage Status](https://coveralls.io/repos/github/doytowin/goooqo/badge.svg?branch=main)](https://coveralls.io/github/doytowin/goooqo?branch=main)
[![CI](https://github.com/doytowin/goooqo/actions/workflows/go.yml/badge.svg)](https://github.com/doytowin/goooqo/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/doytowin/goooqo/rdb)](https://goreportcard.com/report/github.com/doytowin/goooqo/rdb)
[![Go Reference](https://pkg.go.dev/badge/github.com/doytowin/goooqo/rdb.svg)](https://pkg.go.dev/github.com/doytowin/goooqo/rdb)

GoooQo
---

## Introduction

GoooQo is a database access framework implemented in Go, based on OQM techniques. It relies entirely on objects to construct various database query statements, eliminating boilerplate code associated with traditional ORM frameworks, and assists developers in achieving automated database access operations.

The first three Os in the name GoooQo stands for the three major object concepts in the OQM technique:

- `Entity Object` is used to map the static part in the CRUD statements, e.g. table name and column names;
- `Query Object` is used to map the dynamic part in the CRUD statements, e.g. filter conditions, pagination and sorting clause;
- `View Object` is used to map the static part in the complex query statements, e.g. table names, column names, group-by clause and joins.

Check this [article](https://blog.doyto.win/post/introduction-to-goooqo-en/) for more details. 

Check this [demo](https://github.com/doytowin/goooqo-demo) to take a quick tour.

Product documentation: https://goooqo.docs.doyto.win/

## Quick start

### Init project

Use `go mod init` to init the project and add GoooQo dependency:

```
go get -u github.com/doytowin/goooqo/rdb
```

Init a database connection `db` and a transaction manager `tm`:

```go
package main

import (
	"database/sql"
	"github.com/doytowin/goooqo/rdb"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, _ := sql.Open("sqlite3", "./test.db")
	tm := rdb.NewTransactionManager(db)
	//...
}
```

### Create a data access interface

Suppose we have the following user table in `test.db`:

| id | name  | score | memo  | deleted |
|----|-------|-------|-------|---------|
| 1  | Alley | 80    | Good  | false   |
| 2  | Dave  | 75    | Well  | false   |
| 3  | Bob   | 60    |       | false   |
| 4  | Tim   | 92    | Great | true    |
| 5  | Emy   | 100   | Great | false   |

We define an entity object, a query object and a database access interface for this table:

```go
package main

import . "github.com/doytowin/goooqo/core"

type UserEntity struct {
	Int64Id
	Name    *string `json:"name"`
	Score   *int    `json:"score"`
	Memo    *string `json:"memo"`
	Deleted *bool   `json:"deleted"`
}

type UserQuery struct {
	PageQuery
	IdGt     *int64
	IdIn     *[]int64
	ScoreLt  *int
	MemoNull *bool
	MemoLike *string
	Deleted  *bool
	UserOr   *[]UserQuery

	ScoreLtAvg *UserQuery `subquery:"select avg(score) from t_user"`
	ScoreLtAny *UserQuery `subquery:"SELECT score FROM t_user"`
	ScoreLtAll *UserQuery `subquery:"SELECT score FROM t_user"`
	ScoreGtAvg *UserQuery `select:"avg(score)" from:"t_user"`

	Role *RoleQuery `entitypath:"user,role"`
	Perm *PermQuery `entitypath:"user,role,perm"`
}

var UserDataAccess TxDataAccess[UserEntity]
```

After establishing a database connection and creating the transaction manager `tm`,
initialize the `UserDataAccess` interface to perform CRUD operations:

```go
UserDataAccess = rdb.NewTxDataAccess[UserEntity](tm)
```

### Query examples: 

```go
userQuery := UserQuery{ScoreLt: P(80)}
users, err := UserDataAccess.Query(ctx, userQuery)
// SQL="SELECT id, name, score, memo, deleted FROM t_user WHERE score < ?" args="[80]"

userQuery := UserQuery{PageQuery: PageQuery{Size: 20, Sort: "id,desc;score"}, MemoLike: P("Great")}
users, err := UserDataAccess.Query(ctx, userQuery)
// SQL="SELECT id, name, score, memo, deleted FROM t_user WHERE memo LIKE ? ORDER BY id DESC, score LIMIT 20 OFFSET 0" args="[Great]"

userQuery := UserQuery{IdIn: &[]int64{1, 4, 12}, Deleted: P(true)}
users, err := UserDataAccess.Query(ctx, userQuery)
// SQL="SELECT id, name, score, memo, deleted FROM t_user WHERE id IN (?, ?, ?) AND deleted = ?" args="[1 4 12 true]"

userQuery := UserQuery{UserOr: &[]UserQuery{{IdGt: P(int64(10)), MemoNull: P(true)}, {ScoreLt: P(80), MemoLike: P("Good")}}}
users, err := UserDataAccess.Query(ctx, userQuery)
// SQL="SELECT id, name, score, memo, deleted FROM t_user WHERE (id > ? AND memo IS NULL OR score < ? AND memo LIKE ?)" args="[10 80 Good]"

userQuery := UserQuery{ScoreGtAvg: &UserQuery{Deleted: P(true)}, ScoreLtAny: &UserQuery{}}
users, err := UserDataAccess.Query(ctx, userQuery)
// SQL="SELECT id, name, score, memo, deleted FROM t_user WHERE score > (SELECT avg(score) FROM t_user WHERE deleted = ?) AND score < ANY(SELECT score FROM t_user)" args="[true]"
```

For more CRUD examples, please refer to: https://goooqo.docs.doyto.win/api/crud

### Code Generation

GoooQo provides a code generation tool, `gooogen`, which supports automatically generating dynamic query construction methods for query objects.

#### Install `gooogen`

Install the `gooogen` tool using the following command:

```bash
go install github.com/doytowin/goooqo/gooogen@latest
```

#### Add Generate Directive

Add the `go:generate gooogen` directive to the query object definition. For example:

```go
//go:generate gooogen -type sql -o user_query_builder.go
type UserQuery struct {
    // ...
}
```

- **`-type`**: (Optional) Specifies the type of query language to generate, e.g. `sql`.
- **`-f`**: (Optional) Defines the name of the input file which contains a query object, e.g. `user.go`.
- **`-o`**: (Optional) Defines the name of the generated file, e.g. `user_query_builder.go`.

#### Generate Code

Run the `go generate` command to generate the corresponding query construction methods in the specified file.

### Transaction Examples

Use `TransactionManager#StartTransaction` to start a transaction, then manually commit or rollback the transaction:
```go
tc, err := UserDataAccess.StartTransaction(ctx)
userQuery := UserQuery{ScoreLt: P(80)}
cnt, err := UserDataAccess.DeleteByQuery(tc, userQuery)
if err != nil {
	err = tc.Rollback()
	return 0, err
}
err = tc.Commit()
return cnt, err
```

Or use `TransactionManager#SubmitTransaction` to commit the transaction via callback function:
```go
err := tm.SubmitTransaction(ctx, func(tc TransactionContext) (err error) {
	// transaction body
	return
})
```

License
---
This project is under the [BSD 3-Clause Clear License](https://spdx.org/licenses/BSD-3-Clause-Clear).
