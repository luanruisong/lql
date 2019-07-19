package lql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"strings"
	"unicode"
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
	d      func(...interface{})
	sqlDebug bool
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
			d:      defDebuger,
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

func getStructValueAndType(p interface{}) (reflect.Value, reflect.Type) {
	v := reflect.ValueOf(p)
	if v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}
	return v, v.Type()
}

func convStructField(p interface{}) string {
	switch p.(type) {
	case int:
		s := p.(int)
		if s != 0 {
			return fmt.Sprintf("%d", p)
		}
	case int64:
		s := p.(int64)
		if s != 0 {
			return fmt.Sprintf("%d", p)
		}
	case float64:
		s := p.(float64)
		if s != 0 {
			return fmt.Sprintf("%f", p)
		}
	case string:
		return p.(string)
	}
	return ""
}

func snakeString(s string) string {
	data := make([]rune, 0)
	rs := []rune(s)
	for i, v := range rs {
		if unicode.IsUpper(v) {
			v = unicode.ToLower(v)
			if i > 0 {
				data = append(data, '_')
			}
		}
		data = append(data, v)
	}
	return strings.ToLower(string(data))
}

func getColumnDateTypeAndLength(p interface{}) string {
	switch p.(type) {
	case int,int32:
		return "int(11)"
	case int64:
		return "bigint(11)"
	case float32:
		return "float(11)"
	case float64:
		return "double(11)"
	}
	return "varchar(20)"
}

func inStringArrays(s string,ss []string)bool {
	for _,v := range ss {
		if s == v {
			return true
		}
	}
	return false
}