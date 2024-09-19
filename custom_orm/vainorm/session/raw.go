package session

import (
	"database/sql"
	"go_stu/custom_orm/vainorm/clause"
	"go_stu/custom_orm/vainorm/dialect"
	"go_stu/custom_orm/vainorm/log"
	"go_stu/custom_orm/vainorm/schema"
	"strings"
)

type Session struct {
	db      *sql.DB
	sql     strings.Builder
	sqlVars []interface{}

	dialect dialect.Dialect

	// 对应的表结构
	refTable *schema.Schema

	// 操作的语句
	clause clause.Clause

	tx *sql.Tx
}

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{db: db, dialect: dialect}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
}

func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

func (s *Session) Exec() (sql.Result, error) {
	defer s.Clear()
	log.Infof("SQL: %s, VARS: %v", s.sql.String(), s.sqlVars)

	result, err := s.db.Exec(s.sql.String(), s.sqlVars...)
	if err != nil {
		log.Error(err)

		return nil, err
	}
	return result, nil
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Infof("SQL: %s, VARS: %v", s.sql.String(), s.sqlVars)

	return s.db.QueryRow(s.sql.String(), s.sqlVars...)
}

func (s *Session) QueryRows() (*sql.Rows, error) {
	defer s.Clear()
	log.Infof("SQL: %s, VARS: %v", s.sql.String(), s.sqlVars)

	rows, err := s.db.Query(s.sql.String(), s.sqlVars...)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return rows, nil
}

type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}
