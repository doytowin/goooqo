[![Sonar Stats](https://sonarcloud.io/api/project_badges/measure?project=win.doyto.goooqo&metric=alert_status)](https://sonarcloud.io/dashboard?id=win.doyto.goooqo)
[![Code Lines](https://sonarcloud.io/api/project_badges/measure?project=win.doyto.goooqo&metric=ncloc)](https://sonarcloud.io/component_measures?id=win.doyto.goooqo&metric=ncloc)
[![Coverage Status](https://sonarcloud.io/api/project_badges/measure?project=win.doyto.goooqo&metric=coverage)](https://sonarcloud.io/component_measures?id=win.doyto.goooqo&metric=coverage)

GoooQo - An OQM Implementation That Can Automatically Build SQL Statements from Objects
---

## Introduction to OQM

OQM is a technology that constructs database query statements only through objects, 
focusing on the mapping relationship between object-oriented programming languages and database query languages.

The biggest difference between OQM (object-query mapping) technology and traditional 
ORM (object-relational mapping) technology is that OQM proposes to build CRUD statements directly through objects.

The core function of OQM is to build query clauses for a table through a query object, 
which is also the origin of Q in OQM.

Another significant discovery in OQM technology is that the field names in query objects and the conditions in query clauses can be converted interchangeably.

In this way, we only need to create an entity object and a query object to build CRUD statements. 
The entity object is used to determine the table name and the column names, 
and the instance of the query object is used to control the construction of the query clause.

## Introduction to GoooQo

`GoooQo` is an OQM implementation that can automatically build SQL statements from objects.

The first three Os in the name `GoooQo` stands for the three major object concepts in the OQM technique:

- `Entity Object` is used to map the static part in the SQL statements, such as table name and column names;
- `Query Object` is used to map the dynamic part in the SQL statements, such as filter conditions, pagination, and sorting;
- `View Object` is used to map the static part in the complex query statements, such as table names, column names, nested views, and group-by columns.

Where `Qo` represents `Query Object`, which is the core concept in the OQM technique.

Check this [article](https://blog.doyto.win/post/introduction-to-goooqo-en/) for more details. 

Check this [demo](https://github.com/doytowin/goooqo-demo) to take a quick tour.

Product documentation: https://goooqo.docs.doyto.win/

## Quick Start

### Init Project

Use `go mod init` to init the project and add GoooQo dependency:

```
go get -u github.com/doytowin/goooqo
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
type UserEntity struct {
	Int64Id
	Name    *string `json:"name"`
	Score   *int    `json:"score"`
	Memo    *string `json:"memo"`
	Deleted *bool   `json:"deleted"`
}

func (u UserEntity) GetTableName() string {
    return "t_user"
}

type UserQuery struct {
    PageQuery
    IdGt *int64
    IdIn *[]int64
    ScoreLt *int
    MemoNull *bool
    MemoLike *string
    Deleted *bool
    UserOr *[]UserQuery
    
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

> This is currently an experimental project and is not suitable for production usage.
