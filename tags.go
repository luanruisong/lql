package lql

import (
	"fmt"
	"sort"
	"strings"
)

const (
	TAG_NAME_COLUMN = "sql"
	TAG_NAME_PK     = "pk"
	TAG_NAME_ORDER  = "order"
	TAG_NAME_SORT   = "sort"
)

type (
	mysqlTag struct {
		tname  string
		pk     string
		fields []mysqlFileTag
	}
	mysqlFileTag struct {
		name  string
		value interface{}
		tags  map[string]string
	}
)

func (mt *mysqlTag) getNotEmptyField() []mysqlFileTag {
	res := make([]mysqlFileTag, 0)
	for _, v := range mt.fields {
		if len(convStructField(v.value)) > 0 {
			res = append(res, v)
		}
	}
	return res
}

func (mt *mysqlTag) getOrderField() []mysqlFileTag {
	res := make(map[string]mysqlFileTag, 0)
	for _, v := range mt.fields {
		if oi := v.tags[TAG_NAME_ORDER]; len(oi) > 0 {
			res[oi] = v
		}
	}
	l := len(res)
	keyList := make([]string, l)
	i := 0
	for k, _ := range res {
		keyList[i] = k
		i++
	}
	sort.Strings(keyList)
	finalRes := make([]mysqlFileTag, l)
	for i, v := range keyList {
		finalRes[i] = res[v]
	}
	return finalRes
}

func (mt *mysqlTag) sqlInsert() (string, []interface{}) {
	fs := mt.getNotEmptyField()
	if len(fs) > 0 {
		fields := make([]string, 0)
		zw := make([]string, 0)
		param := make([]interface{}, 0)
		for _, v := range fs {
			fields = append(fields, v.name)
			param = append(param, v.value)
			zw = append(zw, "?")
		}
		sql := fmt.Sprintf(sql_insert, mt.tname, strings.Join(fields, ","), strings.Join(zw, ","))
		return sql, param
	}
	return "", nil
}

func (mt *mysqlTag) sqlSelect(column []string, pageSize, pageNo int) (string, []interface{}) {
	cs := "*"
	if len(column) > 0 {
		cs = strings.Join(column, ",")
	}
	//根据结构体名转换蛇形命名的表名
	//拼接处理好的查询列
	sql := fmt.Sprintf(sql_select, cs, mt.tname)
	param := make([]interface{}, 0)
	//根据结构体变量的非空属性加载查询条件
	filter := mt.getNotEmptyField()
	if len(filter) > 0 {
		sql += " where "
		fields := make([]string, 0)
		for _, v := range filter {
			fields = append(fields, v.name+" = ?")
			param = append(param, v.value)
		}
		sql += strings.Join(fields, " and ")
	}

	//查询是否制定排序
	order := mt.getOrderField()
	if len(order) > 0 {
		orderList := make([]string, 0)
		for _, v := range order {
			orderList = append(orderList, v.name+" "+v.tags[TAG_NAME_SORT])
		}
		sql += " order by " + strings.Join(orderList, ",")

	}

	//分页sql组建
	if pageSize > 0 && pageNo > 0 {
		start := (pageNo - 1) * pageSize
		sql += " limit ?,?"
		param = append(param, start, pageSize)
	}

	return sql, param
}

func (mt *mysqlTag) sqlUpdate() (string, []interface{}) {
	fs := mt.getNotEmptyField()
	var pk mysqlFileTag
	for i, v := range fs {
		if v.name == mt.pk {
			pk = v
			fs = append(fs[:i], fs[i+1:]...)
			break
		}
	}
	if len(fs) > 0 {
		fields := make([]string, 0)
		param := make([]interface{}, 0)
		for _, v := range fs {
			fields = append(fields, v.name+" = ?")
			param = append(param, v.value)
		}
		pksql := pk.name + " = ?"
		param = append(param, pk.value)
		sql := fmt.Sprintf(sql_update, mt.tname, strings.Join(fields, ","), pksql)
		return sql, param
	}
	return "", nil
}

func structConvMysqlTag(p interface{}) *mysqlTag {
	v, t := getStructValueAndType(p)
	var data = make(map[string]interface{})

	mt := mysqlTag{
		tname:  snakeString(t.Name()),
		pk:     "",
		fields: make([]mysqlFileTag, 0),
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name := f.Tag.Get(TAG_NAME_COLUMN)
		if len(name) == 0 {
			name = snakeString(f.Name)
		}

		if len(mt.pk) == 0 {
			if len(f.Tag.Get(TAG_NAME_PK)) > 0 {
				mt.pk = name
			}
		}
		value := v.Field(i).Interface()
		mft := mysqlFileTag{
			name:  name,
			value: value,
			tags: map[string]string{
				TAG_NAME_ORDER: f.Tag.Get(TAG_NAME_ORDER),
				TAG_NAME_SORT:  f.Tag.Get(TAG_NAME_SORT),
			},
		}
		mt.fields = append(mt.fields, mft)
		data[name] = value
	}
	return &mt
}
