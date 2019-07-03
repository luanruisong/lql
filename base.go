package lql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Host        string
	UserName    string
	Password    string
	Database    string
	MaxConn     int
	MaxIdleConn int
}

type DBPool struct {
	currDB *sql.DB
	isConn bool
	err    error
	d  func(...interface{})
}

func defDebuger(i ...interface{}) {
	fmt.Println(i...)
}

func NewDataSource(config Config) *DBPool {
	return initDb(config)
}

func convConfig2Str(config Config) string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", config.UserName, config.Password, config.Host, config.Database)
	//log.Info("init db conn pool --->",dsn)
	return dsn
}

func initDb(config Config) *DBPool {
	connStr := convConfig2Str(config)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return &DBPool{
			isConn: false,
			err:    err,
		}
	} else {
		db.SetMaxOpenConns(config.MaxConn)
		db.SetMaxIdleConns(config.MaxIdleConn)
		return &DBPool{
			currDB: db,
			isConn: true,
			d: defDebuger,
		}
	}
}

func parseRow(r *sql.Rows) map[string]string {
	cols, _ := r.ColumnTypes() // Remember to check err afterwards
	values := make([]sql.RawBytes, len(cols))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	err := r.Scan(scanArgs...)
	if err != nil {
		return nil
	}
	item := make(map[string]string, 0)
	for i, col := range values {
		if col != nil {
			k := cols[i].Name()
			item[k] = string(col)
		}
	}
	return item
}
