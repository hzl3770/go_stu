package vainorm

import (
	"database/sql"
	"go_stu/custom_orm/vainorm/dialect"
	"go_stu/custom_orm/vainorm/log"
	"go_stu/custom_orm/vainorm/session"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driverName, dataSourceName string) (*Engine, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Error(err)
		return nil, err
	}

	dial, ok := dialect.GetDialect(driverName)
	if !ok {
		log.Errorf("Dialect Not Found: %s", driverName)
		return nil, err
	}

	e := &Engine{db: db, dialect: dial}
	log.Infof("Connect database success")
	return e, nil
}

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Error(err)
	}
	log.Infof("Close database success")
}

func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}

type TxFunc func(*session.Session) (interface{}, error)

// Transaction https://stackoverflow.com/questions/16184238/database-sql-tx-detecting-commit-or-rollback
func (e *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := e.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			_ = s.Rollback() // err is non-nil; don't change it
		} else {
			defer func() {
				if err != nil {
					_ = s.Rollback()
				}
			}()
			err = s.Commit() // err is nil; if Commit returns error update err
		}
	}()

	return f(s)
}
