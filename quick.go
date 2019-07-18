package lql

import (
	"database/sql"
	"errors"
)

func (db *DBPool) QuickInsert(p interface{}) (sql.Result, error) {
	mt := structConvMysqlTag(p)
	sql, param := mt.sqlInsert()
	if len(sql) > 0 && len(param) > 0 {
		return db.Exec(sql, param...)
	}
	return nil, NO_FIELDS
}

func (db *DBPool) QuickFind(selector interface{}, columns []string) ([]map[string]string, error) {
	return db.QuickPageFind(selector, columns, 0, 0)
}

func (db *DBPool) QuickPageFind(selector interface{}, columns []string, pageSize, pageNo int) ([]map[string]string, error) {
	if selector == nil {
		return nil, NO_TABLE
	}
	mt := structConvMysqlTag(selector)
	sql, param := mt.sqlSelect(columns, pageSize, pageNo)
	return db.QueryRows(sql, param...)
}

func (db *DBPool) QuickUpdate(p interface{}) (sql.Result, error) {
	if p == nil {
		return nil, NO_TABLE
	}
	mt := structConvMysqlTag(p)

	if len(mt.pk) == 0 {
		return nil, errors.New("can not get pk column")
	}
	sql, param := mt.sqlUpdate()
	return db.Exec(sql, param...)

}


func  (db *DBPool) QuickCheckTableStruct(p interface{}) {
	mt := structConvMysqlTag(p)
	row,err := db.QueryRow(mt.sqlCheckTbExists())
	if err != nil {
		db.debug("check table error",err.Error())
		return
	}
	if len(row) > 0 {
		currColumn := db.Query(mt.sqlCheckColumn())
		if currColumn == nil {
			db.debug("create table error cannot find column")
			return
		}
		currColumnList := make([]string,len(currColumn))
		for i,v := range currColumn {
			currColumnList[i] = v["Field"]
		}
		for _,v := range mt.getField() {
			if !inStringArrays(v.name,currColumnList){
				_,err = db.Exec(mt.sqlAddColumn(v))
				if err != nil {
					db.debug("alter table add column error",err.Error())
					return
				}
			}
		}
	} else {
		_,err = db.Exec(mt.sqlCreateTable())
		if err != nil {
			db.debug("create table error",err.Error())
			return
		}
	}
}
