[![Sonar Stats](https://sonarcloud.io/api/project_badges/measure?project=win.doyto.goooqo&metric=alert_status)](https://sonarcloud.io/dashboard?id=win.doyto.goooqo)
[![Code Lines](https://sonarcloud.io/api/project_badges/measure?project=win.doyto.goooqo&metric=ncloc)](https://sonarcloud.io/component_measures?id=win.doyto.goooqo&metric=ncloc)
[![Coverage Status](https://sonarcloud.io/api/project_badges/measure?project=win.doyto.goooqo&metric=coverage)](https://sonarcloud.io/component_measures?id=win.doyto.goooqo&metric=coverage)

GoooQo
---

## 项目介绍

GoooQo是一个采用Go语言实现的基于OQM技术的的数据库访问框架，完全依赖对象构建各种数据库查询语句，
消除各种使用传统ORM框架带来的样板代码，帮助开发人员实现自动化的数据库访问操作。

GoooQo中用于动态构建查询子句的对象被称为查询对象 (Query Object)，即GoooQo中`Qo`。
而GoooQo名称中的前三个`O`代表了OQM技术中的三大对象概念：

- `Entity Object`用于映射增删查改语句中的静态部分，例如表名和列名；
- `Query Object`用于映射增删查改语句中的动态部分，例如过滤条件、分页和排序；
- `View Object`用于映射复杂查询语句中的静态部分，例如表名、列名、分组子句、嵌套视图和各类连接。

查看这篇[文章](https://blog.doyto.win/post/introduction-to-goooqo-en/)，了解更多详情。

参考[demo](https://github.com/doytowin/goooqo-demo)，快速上手。

产品文档: https://goooqo.docs.doyto.win/v/zh

## 快速开始

### 初始化项目

使用`go mod init`初始化项目并添加GoooQo依赖：

```
go get -u github.com/doytowin/goooqo/rdb
```

初始化数据库连接和事务管理器：

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

### 创建数据访问接口

假设我们在`test.db`中有以下用户表：

| id | name  | score | memo  | deleted |
|----|-------|-------|-------|---------|
| 1  | Alley | 80    | Good  | false   |
| 2  | Dave  | 75    | Well  | false   |
| 3  | Bob   | 60    |       | false   |
| 4  | Tim   | 92    | Great | true    |
| 5  | Emy   | 100   | Great | false   |

我们为该表定义一个实体对象和一个查询对象，以及一个用于执行数据库访问操作的`UserDataAccess`接口：

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
	ScoreLtAll *UserQuery `subquery:"SELECT score FROM t_user"`
	ScoreGtAvg *UserQuery `select:"avg(score)" from:"t_user"`

	Role *RoleQuery `entitypath:"user,role"`
	Perm *PermQuery `entitypath:"user,role,perm"`
}

var UserDataAccess TxDataAccess[UserEntity]
```

在建立数据库连接并创建事务管理器后，初始化`UserDataAccess`接口：

```go
UserDataAccess := rdb.NewTxDataAccess[UserEntity](tm)
```

### 代码生成

GoooQo提供的代码生成工具`gooogen`支持为查询对象自动生成动态查询语句的构建方法。

#### 安装 `gooogen`

通过以下命令安装代码生成工具`gooogen`：

```bash
go install github.com/doytowin/goooqo/gooogen@latest
```

#### 添加生成指令

在查询对象的定义上添加注释`go:generate gooogen`。例如：

```go
//go:generate gooogen -type sql -o user_query_builder.go
type UserQuery struct {
//...
}
```

- **`-type`**：（可选）指定生成的查询语句类型，如 `sql`。
- **`-f`**：（可选）定义包含查询对象的文件的名称，如 `user.go`。
- **`-o`**：（可选）定义生成文件的名称，如 `user_query_builder.go`。

#### 生成代码

执行`go generate`命令即可在指定的文件中生成相应的查询语句构建方法。

### 查询示例

```go
userQuery := UserQuery{ScoreLt: P(80)}
users, err := UserDataAccess.Query(ctx, userQuery)
// SQL="SELECT id, name, score, memo, deleted FROM t_user WHERE score < ?" args="[80]"

userQuery := UserQuery{PageQuery: PageQuery{PageSize: P(20), Sort: P("id,desc;score")}, MemoLike: P("Great")}
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

更多接口调用示例请参考：https://goooqo.docs.doyto.win/v/zh/api/crud

### 事务示例

使用`TransactionManager#StartTransaction`开启事务，手动提交或者回滚事务：
```go
tc, err := UserDataAccess.StartTransaction(ctx)
userQuery := UserQuery{ScoreLt: PInt(80)}
cnt, err := UserDataAccess.DeleteByQuery(tc, userQuery)
if err != nil {
	err = tc.Rollback()
	return 0
}
err = tc.Commit()
return cnt
```

或者使用`TransactionManager#SubmitTransaction`通过回调的方式提交事务：
```go
err := tm.SubmitTransaction(ctx, func(tc TransactionContext) (err error) {
	// transaction body
	return
})
```

许可证
---
本项目遵循[BSD 3-Clause Clear License](https://spdx.org/licenses/BSD-3-Clause-Clear)。
