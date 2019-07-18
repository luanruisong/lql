package lql

import (
	"errors"
)

var (
	NO_TABLE  = errors.New("can not get table name")
	NO_FIELDS = errors.New("can not get fields")

	sql_insert = "insert into %s (%s) values (%s)"
	sql_update = "update %s set %s where %s"
	sql_select = "select %s from %s"

	sql_check_tb = "show tables like '%s'"
	sql_create_table = "CREATE TABLE %s(%s);"
	sql_check_column = "show columns from %s"
	sql_add_column = "ALTER TABLE %s ADD %s"

)
