[![Sonar Stats](https://sonarcloud.io/api/project_badges/measure?project=win.doyto.goooqo&metric=alert_status)](https://sonarcloud.io/dashboard?id=win.doyto.goooqo)
[![Code Lines](https://sonarcloud.io/api/project_badges/measure?project=win.doyto.goooqo&metric=ncloc)](https://sonarcloud.io/component_measures?id=win.doyto.goooqo&metric=ncloc)
[![Coverage Status](https://sonarcloud.io/api/project_badges/measure?project=win.doyto.goooqo&metric=coverage)](https://sonarcloud.io/component_measures?id=win.doyto.goooqo&metric=coverage)

GoooQo - An OQM Implementation That Can Automatically Build SQL Statements from Objects
---

## Introduction to OQM

OQM technology focuses on studying the mapping relationship between object-oriented programming languages and
database query languages. It is a technology that constructs database query statements only through objects.

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

Check our [wiki](https://github.com/doytowin/goooqo/wiki) for the incoming documentations.

## Quick Tutorial

First, use `go mod init` to initialize the project and add GoooQo:
```bash
go get -u github.com/doytowin/goooqo
```

Then, define entity objects and query objects for the user table:

```go
type UserEntity struct {
	Int64Id
	Name  *string `json:"name"`
	Score *int `json:"score"`
	Memo  *string `json:"memo"`
}

func (u UserEntity) GetTableName() string {
	return "t_user"
}

type UserQuery struct {
	PageQuery
	IdIn *[]int
	ScoreLt *int
	ScoreLt1 *UserQuery `subquery:"select:avg(score),from:UserEntity"`
	MemoNull *bool
	UserOr   *UserQuery
}
```

Finally, initialize the database connection, transaction manager, and data access interface as follows, 
and then we can use `userDataAccess` to access the table `t_user`:
```go
package main

import (
	"github.com/doytowin/goooqo"
	"github.com/doytowin/goooqo/rdb"
)

func main() {
	db := rdb.Connect("local.properties")
	defer rdb.Disconnect(db)
	
	tm := rdb.NewTransactionManager(db)
	
	userDataAccess := rdb.NewTxDataAccess[UserEntity](tm)

	//...
}
```

#### Paging Query Example
```go
userQuery := UserQuery{PageQuery: PageQuery{PageSize: PInt(2)}, ScoreLt: PInt(80), MemoStart: PStr("Well")}
page, err := userDataAccess.Page(ctx, &userQuery)
``` 

This code snippet will execute the following SQL statements:
```sql
SELECT id, score, memo FROM t_user WHERE score < ? AND memo LIKE ? LIMIT 2 OFFSET 0; -- args="[80 Well%]"
SELECT count(0) FROM t_user WHERE score < ? AND memo LIKE ? -- args="[80 Well%]"
```

#### Transaction Example

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

---
> This is currently an experimental project and is not suitable for production usage.
