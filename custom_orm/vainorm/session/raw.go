package session

import (
	"database/sql"
	"fmt"
	"go_stu/custom_orm/vainorm/dialect"
	"go_stu/custom_orm/vainorm/log"
	"go_stu/custom_orm/vainorm/schema"
	"reflect"
	"strings"
)

type Session struct {
	db      *sql.DB
	sql     strings.Builder
	sqlVars []interface{}

	dialect  dialect.Dialect
	refTable *schema.Schema
}

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{db: db, dialect: dialect}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
}

func (s *Session) DB() *sql.DB {
	return s.db
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

func (s *Session) Model(value interface{}) *Session {
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is not set")
	}
	return s.refTable
}

func (s *Session) CreateTable() error {
	table := s.RefTable()
	var columns []string
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err := s.Raw("CREATE TABLE ?(?)", table.Name, desc).Exec()
	return err
}

func (s *Session) DropTable() error {
	_, err := s.Raw("DROP TABLE IF EXISTS ?", s.RefTable().Name).Exec()
	return err
}

func (s *Session) HasTable() bool {
	sql, vars := s.dialect.TableExistSQL(s.RefTable().Name)
	row := s.Raw(sql, vars...).QueryRow()
	var tableName string
	_ = row.Scan(&tableName)
	return tableName == s.RefTable().Name
}
