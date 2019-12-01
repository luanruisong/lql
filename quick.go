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

func (db *DBPool) QuickFind(selector interface{}, columns ...string) ([]map[string]string, error) {
	return db.QuickPageFind(selector, 0, 0, columns)
}

func (db *DBPool) QuickPageFind(selector interface{}, pageSize, pageNo int, columns []string) ([]map[string]string, error) {
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

func (db *DBPool) QuickCheckTableStruct(p interface{}) {
	mt := structConvMysqlTag(p)
	db.debug("start check table", mt.tname)
	row, err := db.QueryRow(mt.sqlCheckTbExists())
	if err != nil {
		db.debug("check table error", err.Error())
		return
	}
	if len(row) > 0 {
		db.debug("table", mt.tname, "exists start check column")
		currColumn := db.Query(mt.sqlCheckColumn())
		if currColumn == nil {
			db.debug("create table error cannot find column")
			return
		}
		currColumnList := make([]string, len(currColumn))
		for i, v := range currColumn {
			currColumnList[i] = v["Field"]
		}
		for _, v := range mt.getField() {
			if !inStringArrays(v.name, currColumnList) {
				db.debug("table", mt.tname, "column", v.name, "not exists add column")
				_, err = db.Exec(mt.sqlAddColumn(v))
				if err != nil {
					db.debug("alter table add column error", err.Error())
					return
				}
			}
		}
	} else {
		db.debug("table", mt.tname, "not exists create table")
		_, err = db.Exec(mt.sqlCreateTable())
		if err != nil {
			db.debug("create table error", err.Error())
			return
		}
		db.QuickInsert(p)
	}
}
