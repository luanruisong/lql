package lql

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
	"sort"
)

const (
	COLUMN_TAG_NAME = "sql"
	PK_TAG_NAME = "pk"
	ORDER_TAG_NAME = "order"
	SORT_TAG_NAME = "sort"
)

var (
	NO_TABLE = errors.New("can not get table name")
)

func structConvMap(p interface{}) (string,string,string,map[string]interface{}) {
	v, t := getStructValueAndType(p)
	var data = make(map[string]interface{})
	var odrkeymap = make(map[string]string)
	pkName,order := "",""
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name := f.Tag.Get(COLUMN_TAG_NAME)
		if len(name) == 0 {
			name = snakeString(f.Name)
		}
		if len(pkName) == 0 {
			if len(f.Tag.Get(PK_TAG_NAME)) > 0 {
				pkName = name
			}
		}
		if odrkey := f.Tag.Get(ORDER_TAG_NAME);len(odrkey) > 0 {
			odrkeymap[odrkey] = name + " " + f.Tag.Get(SORT_TAG_NAME)
		}

		value := v.Field(i).Interface()
		valueStr := convStructField(value)
		if len(valueStr) == 0 {
			continue
		}
		data[name] = value
	}

	if len(odrkeymap) > 0 {
		keys := make([]string,0)
		values := make([]string,0)
		for i,_ := range odrkeymap {
			keys = append(keys,i)
		}
		sort.Strings(keys)
		for _,v := range keys {
			values = append(values,odrkeymap[v])
		}
		order = strings.Join(values,",")
	}
	return snakeString(t.Name()),pkName,order,data
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
			return fmt.Sprintf("%d",p)
		}
	case int64:
		s := p.(int64)
		if s != 0 {
			return fmt.Sprintf("%d",p)
		}
	case float64:
		s := p.(float64)
		if s != 0 {
			return fmt.Sprintf("%f",p)
		}
	case string:
		return p.(string)
	}
	return ""
}


func snakeString(s string) string {
	data := make([]rune,0)
	rs := []rune(s)
	for i,v := range rs {
		if unicode.IsUpper(v) {
			v = unicode.ToLower(v)
			if i > 0 {
				data = append(data,'_')
			}
		}
		data = append(data,v)
	}
	return strings.ToLower(string(data))
}

func (db *DBPool)QuickInsert(p interface{})(sql.Result,error){
	tn,_,_,m := structConvMap(p)
	l := len(m)
	if l > 0 {
		fields := make([]string,0)
		zw := make([]string,0)
		param := make([]interface{},0)
		for i,v := range m {
			fields = append(fields,i)
			param = append(param,v)
			zw = append(zw,"?")
		}
		sql := `insert into %s (%s) values (%s)`
		sql = fmt.Sprintf(sql,tn,strings.Join(fields,","),strings.Join(zw,","))
		return db.Exec(sql,param...)
	}
	return nil,errors.New("can not get fields")
}


func (db *DBPool)buildQuickQuery(selector interface{},columns []string,pageSize, pageNo int)(string,[]interface{}){
	tn,_,order,m := structConvMap(selector)
	sql := `select %s from %s `
	//判断是否自己指定返回列
	cs := "*"
	if len(columns) > 0 {
		cs = strings.Join(columns,",")
	}
	//根据结构体名转换蛇形命名的表名
	//拼接处理好的查询列
	sql = fmt.Sprintf(sql,cs,snakeString(tn))
	param := make([]interface{},0)
	l := len(m)
	//根据结构体变量的非空属性加载查询条件
	if l > 0 {
		sql += " where "
		fields := make([]string,0)
		for i,v := range m {
			fields = append(fields,i + " = ?")
			param = append(param,v)
		}

		sql += strings.Join(fields," and ")
	}
	if len(order) > 0 {
		sql += " order by " + order
	}

	//分页sql组建
	if pageSize > 0 && pageNo > 0 {
		start := (pageNo - 1) * pageSize
		sql += " limit ?,?"
		param = append(param,start,pageSize)
	}

	return sql,param
}

func (db *DBPool)QuickFind(selector interface{},columns []string)([]map[string]string,error){
	return db.QuickPageFind(selector,columns,0,0)
}

func (db *DBPool)QuickPageFind(selector interface{},columns []string,pageSize, pageNo int)([]map[string]string,error){
	if selector == nil {
		return nil,NO_TABLE
	}
	sql,param := db.buildQuickQuery(selector,columns,pageSize, pageNo)
	return db.QueryRows(sql,param...)
}

func (db *DBPool)QuickUpdate(p interface{})(sql.Result,error){
	if p == nil {
		return nil,NO_TABLE
	}
	tn,pkName,_,m := structConvMap(p)

	if len(pkName) == 0 {
		return nil,errors.New("can not get pk column")
	}

	pkValue,ok := m[pkName]
	if !ok {
		return nil,errors.New("can not get pk column value")
	}

	delete(m,pkName)
	l := len(m)
	if l > 0 {
		fields := make([]string,0)
		param := make([]interface{},0)
		for i,v := range m {
			fields = append(fields,i +" = ?")
			param = append(param,v)
		}
		sql := `update %s set %s where %s`
		pksql := pkName + " = ?"
		param = append(param,pkValue)
		sql = fmt.Sprintf(sql,tn,strings.Join(fields,","),pksql)
		return db.Exec(sql,param...)
	}
	return nil,errors.New("can not get fields")

}