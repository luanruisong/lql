package lql

import (
	"database/sql"
)

func (m *DBPool) SetDebuger(debuger func(...interface{})) {
	m.d = debuger
	m.debug("dbpool set debuger success")
}

func (m *DBPool) debug(p ...interface{}) {
	m.d(p...)
}

func (m *DBPool) debugSql(p ...interface{}) {
	if m.sqlDebug {
		m.debug(p...)
	}
}

func (m *DBPool) Query(query string, args ...interface{}) []map[string]string {
	m.debugSql("Query sql :(", query, ")", args)
	rows, err := m.currDB.Query(query, args...)
	if err != nil {
		return nil
	}
	defer rows.Close()
	res := make([]map[string]string, 0)
	for rows.Next() {
		res = append(res, parseRow(rows))
	}
	return res
}

func (m *DBPool) QueryRows(query string, args ...interface{}) ([]map[string]string, error) {
	m.debugSql("Query sql :(", query, ")", args)
	rows, err := m.currDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := make([]map[string]string, 0)
	for rows.Next() {
		res = append(res, parseRow(rows))
	}
	return res, nil
}

func (m *DBPool) QueryRow(query string, args ...interface{}) (map[string]string, error) {
	m.debugSql("Query Row sql :(", query, ")", args)
	rows, err := m.currDB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		return parseRow(rows), nil
	}
	return nil, err
}

func (m *DBPool) Exec(sql string, args ...interface{}) (sql.Result, error) {
	m.debugSql("Exec sql :(", sql, ")", args)
	res, err := m.currDB.Exec(sql, args...)
	if err != nil {
	}
	return res, err
}

func (m *DBPool) Prepare(query string) (*sql.Stmt, error) {
	return m.currDB.Prepare(query)
}

func (m *DBPool) GetDb() *sql.DB {
	return m.currDB
}

func (m *DBPool) Close() {
	m.currDB.Close()
}

func (m *DBPool) Ping() error {
	return m.currDB.Ping()
}

func (m *DBPool) IsConn() bool {
	return m.isConn
}

func (m *DBPool) Err() error {
	return m.err
}
func (m *DBPool) Begin() (*sql.Tx, error) {
	return m.currDB.Begin()
}
func (m *DBPool) OpenSqlDebugger() {
	m.sqlDebug = true
}