package main

import (
	"fmt"
)

type User struct {
	Id      int64  `sql:"id" pk:"1"`
	Name    string `sql:"name" order:"1" sort:"desc"`
	Account string `sql:"account"`
	Pwd     string `sql:"pwd"`
}

func main() {


	dbConfig := Config{

	}
	dbpool := NewDataSource(dbConfig)
	if dbpool.IsConn() {
		fmt.Println("connection ", dbConfig.Host, dbConfig.Database, "database successful")
		defer dbpool.Close()
	} else {
		fmt.Println(dbpool.Err())
	}

	dbpool.QuickFind(User{},nil)
	dbpool.QuickUpdate(User{Id:1,Name:"test update"})
	dbpool.QuickInsert(User{
		Name:    "test name",
		Account: "test account",
		Pwd:     "test pwd",
	})


}
