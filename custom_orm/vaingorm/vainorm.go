package vaingorm

import (
	"database/sql"
	"go_stu/custom_orm/vaingorm/log"
	"go_stu/custom_orm/vaingorm/session"
)

type Engine struct {
	db *sql.DB
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

	e := &Engine{db: db}
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
	return session.New(e.db)
}
