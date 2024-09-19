package vainorm

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"go_stu/custom_orm/vainorm/session"

	"testing"
)

func OpenDB(t *testing.T) *Engine {
	t.Helper()
	dsn := "root:root008@@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True"
	engine, err := NewEngine("mysql", dsn)
	if err != nil {
		t.Fatal("failed to connect", err)
	}
	return engine
}

func TestNewEngine(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
}

type User struct {
	Name string `vainorm:"PRIMARY KEY"`
	Age  int
}

// 在MySQL中，执行DDL（数据定义语言）语句时，事务会自动提交。
// 这意味着在事务中执行的任何DDL操作都会立即提交，而不会被回滚。
func transactionRollback(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_ = s.Model(&User{}).CreateTable()

	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		_, err = s.Insert(&User{"sbydx", 10})
		return nil, errors.New("Error")
	})

	u := &User{}
	_ = s.Where("Name = ? AND Age = ?", "sbydx", 10).First(u)

	if err == nil || u.Name == "Tom" {
		t.Logf("u: %v", u)
		t.Fatal("failed to rollback")
	}
}

func transactionCommit(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		_ = s.Model(&User{}).CreateTable()
		_, err = s.Insert(&User{"Tom", 18})
		return
	})
	u := &User{}
	_ = s.First(u)
	if err != nil || u.Name != "Tom" {
		t.Fatal("failed to commit")
	}
}

func TestEngine_Transaction(t *testing.T) {
	t.Run("rollback", func(t *testing.T) {
		transactionRollback(t)
	})
	t.Run("commit", func(t *testing.T) {
		transactionCommit(t)
	})
}
