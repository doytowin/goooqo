[![Sonar Stats](https://sonarcloud.io/api/project_badges/measure?project=win.doyto.goooqo&metric=alert_status)](https://sonarcloud.io/dashboard?id=win.doyto.goooqo)
[![Code Lines](https://sonarcloud.io/api/project_badges/measure?project=win.doyto.goooqo&metric=ncloc)](https://sonarcloud.io/component_measures?id=win.doyto.goooqo&metric=ncloc)
[![Coverage Status](https://sonarcloud.io/api/project_badges/measure?project=win.doyto.goooqo&metric=coverage)](https://sonarcloud.io/component_measures?id=win.doyto.goooqo&metric=coverage)
<a href="https://www.producthunt.com/posts/goooqo?embed=true&utm_source=badge-featured&utm_medium=badge&utm_souce=badge-goooqo" target="_blank"><img src="https://api.producthunt.com/widgets/embed-image/v1/featured.svg?post_id=516822&theme=light" alt="ProductHunt" style="width: 250px; height: 54px;" width="250" height="54"/></a>

GoooQo
---

## Introduction

GoooQo is a CRUD framework in Golang based on the OQM technique.

The OQM (Object-Query Mapping) technique is a database access technique that constructs database query statements through objects.

OQM proposes a new method to solve the problem of dynamic combination of n query conditions
by mapping 2^n assignment combinations of an object instance with n fields to 2^n combinations of n query conditions.

This approach enables developers to define and construct objects only to build dynamic query statements, 
setting OQM apart from ORM. Such objects are called query objects, which is the Qo in GoooQo.

The first three Os in the name GoooQo stands for the three major object concepts in the OQM technique:

- `Entity Object` is used to map the static part in the CRUD statements, such as table name and column names;
- `Query Object` is used to map the dynamic part in the CRUD statements, such as filter conditions, pagination and sorting clause;
- `View Object` is used to map the static part in the complex query statements, such as table names, column names, group-by clause and joins.

Check this [article](https://blog.doyto.win/post/introduction-to-goooqo-en/) for more details. 

Check this [demo](https://github.com/doytowin/goooqo-demo) to take a quick tour.

Product documentation: https://goooqo.docs.doyto.win/

## Quick Start

### Init Project

Use `go mod init` to init the project and add GoooQo dependency:

```
go get -u github.com/doytowin/goooqo/rdb
```

Init the database connection and transaction manager:

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

We define an entity object and a query object for this table:

```go
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
	ScoreLtAll *UserQuery `subquery:"select score from UserEntity"`
	ScoreGtAvg *UserQuery `select:"avg(score)" from:"UserEntity"`
}
```

Then we create a `userDataAccess` interface to perform CRUD operations:

```go
userDataAccess := rdb.NewTxDataAccess[UserEntity](tm)
```

### Query example: 

```go
userQuery := UserQuery{ScoreLt: P(80)}
users, err := userDataAccess.Query(ctx, userQuery)
// SQL="SELECT id, name, score, memo, deleted FROM t_user WHERE score < ?" args="[80]"

userQuery := UserQuery{PageQuery: PageQuery{PageSize: P(20), Sort: P("id,desc;score")}, MemoLike: P("Great")}
users, err := userDataAccess.Query(ctx, userQuery)
// SQL="SELECT id, name, score, memo, deleted FROM t_user WHERE memo LIKE ? ORDER BY id DESC, score LIMIT 20 OFFSET 0" args="[Great]"

userQuery := UserQuery{IdIn: &[]int64{1, 4, 12}, Deleted: P(true)}
users, err := userDataAccess.Query(ctx, userQuery)
// SQL="SELECT id, name, score, memo, deleted FROM t_user WHERE id IN (?, ?, ?) AND deleted = ?" args="[1 4 12 true]"

userQuery := UserQuery{UserOr: &[]UserQuery{{IdGt: P(int64(10)), MemoNull: P(true)}, {ScoreLt: P(80), MemoLike: P("Good")}}}
users, err := userDataAccess.Query(ctx, userQuery)
// SQL="SELECT id, name, score, memo, deleted FROM t_user WHERE (id > ? AND memo IS NULL OR score < ? AND memo LIKE ?)" args="[10 80 Good]"

userQuery := UserQuery{ScoreGtAvg: &UserQuery{Deleted: P(true)}, ScoreLtAny: &UserQuery{}}
users, err := userDataAccess.Query(ctx, userQuery)
// SQL="SELECT id, name, score, memo, deleted FROM t_user WHERE score > (SELECT avg(score) FROM t_user WHERE deleted = ?) AND score < ANY(SELECT score FROM t_user)" args="[true]"
```

For more CRUD examples, please refer to: https://goooqo.docs.doyto.win/v/zh/api/crud

### Transaction Examples

Use `TransactionManager#StartTransaction` to start a transaction, then manually commit or rollback the transaction:
```go
tc, err := userDataAccess.StartTransaction(ctx)
userQuery := UserQuery{ScoreLt: PInt(80)}
cnt, err := userDataAccess.DeleteByQuery(tc, userQuery)
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
