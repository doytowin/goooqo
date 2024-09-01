[![Sonar Stats](https://sonarcloud.io/api/project_badges/measure?project=win.doyto.goooqo&metric=alert_status)](https://sonarcloud.io/dashboard?id=win.doyto.goooqo)
[![Code Lines](https://sonarcloud.io/api/project_badges/measure?project=win.doyto.goooqo&metric=ncloc)](https://sonarcloud.io/component_measures?id=win.doyto.goooqo&metric=ncloc)
[![Coverage Status](https://sonarcloud.io/api/project_badges/measure?project=win.doyto.goooqo&metric=coverage)](https://sonarcloud.io/component_measures?id=win.doyto.goooqo&metric=coverage)

GoooQo - OQM技术的Golang实现
===

## 项目介绍

GoooQo是一个可以自动从对象构建SQL语句的OQM框架。

OQM技术专注于研究面向对象编程语言和数据库查询语言之间的映射关系，
是一项仅通过对象来构建数据库查询语句的技术。
OQM主要依靠以下三类对象来映射数据库查询语句：
- `Entity Object`用于映射SQL语句中的静态部分，例如表名和列名；
- `Query Object`用于映射SQL语句中的动态部分，例如过滤条件、分页和排序；
- `View Object`用于映射复杂查询语句中的静态部分，例如表名、列名、嵌套视图和分组列。

GoooQo中的前三个o即代表上述三类对象，`Qo`代表`Query Object`，是OQM技术中最核心的对象概念。

查看这篇[文章](https://blog.doyto.win/post/introduction-to-goooqo-en/)，了解更多详情。

访问[demo](https://github.com/doytowin/goooqo-demo)，快速上手。

访问[wiki](https://github.com/doytowin/goooqo/wiki)，查阅产品文档。

## 快速开始

首先，使用`go mod init`初始化项目后，添加GoooQo：
```bash
go get -u github.com/doytowin/goooqo
```

然后，为用户表定义实体对象和查询对象：
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
	IdIn	 *[]int
	ScoreLt  *int
	ScoreLt1 *UserQuery `subquery:"select:avg(score),from:UserEntity"`
	MemoNull *bool
	UserOr   *UserQuery
}
```

最后，初始化数据库连接、事务管理器、数据访问接口如下，即可使用`userDataAccess`访问表`t_user`：
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

#### 分页查询示例：
```go
userQuery := UserQuery{PageQuery: PageQuery{PageSize: PInt(2)}, ScoreLt: PInt(80), MemoStart: PStr("Well")}
page, err := userDataAccess.Page(ctx, &userQuery)
```

这段代码将会执行以下SQL语句：
```sql
SELECT id, score, memo FROM t_user WHERE score < ? AND memo LIKE ? LIMIT 2 OFFSET 0; -- args="[80 Well%]"
SELECT count(0) FROM t_user WHERE score < ? AND memo LIKE ? -- args="[80 Well%]"
```

#### 事务示例：

使用`TransactionManager#StartTransaction`开启事务，手动提交或者回滚事务：
```go
tc, err := userDataAccess.StartTransaction(ctx)
userQuery := UserQuery{ScoreLt: PInt(80)}
cnt, err := userDataAccess.DeleteByQuery(tc, userQuery)
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

---
> 本产品目前尚处于验证阶段，请谨慎用于生产环境。
