package session

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"go_stu/custom_orm/vainorm/dialect"

	"os"
	"testing"
)

var (
	TestDB      *sql.DB
	TestDial, _ = dialect.GetDialect("mysql")
)

func TestMain(m *testing.M) {
	dsn := "root:root008@@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True"
	var err error
	TestDB, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	code := m.Run()
	_ = TestDB.Close()
	os.Exit(code)
}

func NewSession() *Session {
	return New(TestDB, TestDial)
}

func TestSession_Exec(t *testing.T) {
	s := NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text, Age int);").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text, Age int);").Exec()
	result, _ := s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Jack").Exec()
	if count, err := result.RowsAffected(); err != nil || count != 2 {
		t.Fatal("expect 2, but got", count)
	}
}

func TestSession_QueryRows(t *testing.T) {
	s := NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name text);").Exec()
	_, _ = s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Jack").Exec()
	rows, _ := s.Raw("SELECT * FROM User LIMIT 1").QueryRows()

	var names []string
	for rows.Next() {
		var name string
		_ = rows.Scan(&name)
		names = append(names, name)
	}

	if len(names) != 1 {
		t.Fatal("failed to query db")
	}
}
