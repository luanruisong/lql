# lql

sql 在我们的开发过程中基本上是绕不过去的一个坎.  
lql就是给各位酷爱偷懒的小伙伴带来一个半自动sql拼装工具盒。  
欢迎大家一起交流~  


安装
```sh
go get github.com/luanruisong/lql
```

初始化

```go
    import (
        "github.com/luanruisong/lql"
    )


    //init db

    dbconfig := lql.Config{
        Host:"127.0.0.1:3306",
        Database:"databaseName",
        UserName:"root",
        Password:"123456",
        MaxConn:10,
        MaxIdleConn:3,
    }

    db := lql.NewDataSource(dbconfig)


    if !db.isConn {
        //TODO err

    }


    //设置日志打印主要是sql等，可以加入自己的logger
    db.SetDebuger(func (msg ...interface{}){
        fmt.Println(msg...)
    })
    //开启sql打印
    db.OpenSqlDebugger()

    //close db connection
    db.Close()

```

native 使用方式

```go
    rows := db.Query("select id,a from tb_bbb")

    for i,v := range rows{
        fmt.Println(i);
        for i1,v1 := range v{
            fmt.Println(i1,v1)
        }
    }
```


半ORM 使用方式

```go


    type User struct {
        Id   int    `sql:"id" pk:"1"`
        Name string `sql:"name"`
        Age  int    `sql:"age" order:"1" sort:"desc"`
    }

    //快速查询
    //param 1 ：根据name = lihua 查询，

    //返回所有字段
    row := db.QuickFind(User{Name:"lihua"})
    //查询列为 name,age
    row := db.QuickFind(User{Name:"lihua"},"name","age")

    //分页查询
    //param 1 ：根据name = lihua 查询，
    //param 2 ：每页20条数据
    //param 3 ：查询第一页
    //返回所有字段
    row := db.QuickPageFind(User{Name:"lihua"},20,1)
    //查询列为 name,age
    row := db.QuickPageFind(User{Name:"lihua"},20,1,"name","age")

    //快速插入
    //param ：指定插入name = lihua 的数据
    res,err := db.QuickInsert(User{Name:"lihua"})

    //快速修改
    //param ：修改 id = 1 的数据 name = lihua
    res,err := db.QuickUpdate(User{Name:"lihua",Id:1})


    //新加入  数据结构检测
    //如果表不存在，创建表
    //如果表存在，检测字段做增量更新（不包含字段约束，类型等修改）
    db.QuickCheckTableStruct(User{})



```

其他介绍

半orm时，表名使用结构体名称进行蛇形命名转换
如：UserInfo => user_info

lql 的tag 使用分以下几种

tag_name | description
:- | :-
sql   | 字段名，指定方式 => `sql:"user_id"` 不填写按蛇形命名转换|
order | 排序字段，指定方式 => `order:"1"`， `order:"2"` 优先级根据tag排序 |
sort  | 排序方式，指定方式 => `sort:"desc"` 单独指定无效 |
pk    | 主键声明，指定方式 => `pk:"1"` 多个pk只采用第一个，慎重填写 |
dtype | 数据库对应类型，指定方式 => `dtype:"varchar(20)"` 不指定时有默认值 |
cdesc | 数据库字段约束或默认值，指定方式 => `cdesc:"NOT NULL"` 多个约束时 空格排列即可 |
unique | 唯一键指定 => `unique:"1"` |

