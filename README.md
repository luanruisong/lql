# lql




data base util for mysql


```sh
go get github.com/luanruisong/lql
```

初始化

```go
    import (
        "github.com/luanruisong/lql"
    )


    //init db

    dbconfig := factory.Config{
        Host:"127.0.0.1:3306",
        Database:"databaseName",
        UserName:"root",
        Password:"123456",
        MaxConn:10
        MaxIdleConn:3
    }

    db := factory.NewDataSource(dbconfig)


    if !db.isConn {
        //TODO err

    }

    //close db connection
    db1.Close()

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
        Id   int  `sql:"id" pk:"1"`
        Name string `sql:"name"`
        Age int `sql:"age" order:"1" sort:"desc"`
    }

    //快速查询
    //param 1 ：根据name = lihua 查询，
    //param 2 ：查询列为 name,age 参数为nil时 返回所有字段
    row := db.QuickFind(User{Name:"lihua"},[]string{"name","age"})

    //分页查询
    //param 1 ：根据name = lihua 查询，
    //param 2 ：查询列为 name,age 参数为nil时 返回所有字段
    //param 3 ：每页20条数据
    //param 4 ：查询第一页
    row := db.QuickPageFind(User{Name:"lihua"},[]string{"name","age",20,1})

    //快速插入
    //param ：指定插入name = lihua 的数据
    res,err := db.QuickUpdate(User{Name:"lihua"})

    //快速修改
    //param ：修改 id = 1 的数据 name = lihua
    res,err := db.QuickUpdate(User{Name:"lihua",Id:1})

```

tag 介绍


lql 的tag 使用分以下几种

sql   ： 指定在转换sql的时候的字段名
order ： 标明在快速查询的时候 是需要排序的字段，多个order可以并存 优先级根据order内容来确定
sort  ： 排序方式 和 order的时候可以指明 desc 不然排序使用默认的asc排序
pk    ： 快速修改时，声明pk是表示当前字段为表的主键，多个pk只采用第一个，慎重填写


