package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go_stu/custom_orm/vaingorm"
	"log"
)

func main() {
	//origin()
	dsn := "root:root008@@tcp(127.0.0.1:3306)/?charset=utf8mb4&parseTime=True"

	engine, err := vaingorm.NewEngine("mysql", dsn)
	if err != nil {
		log.Println(err)
		return
	}

	defer engine.Close()

	s := engine.NewSession()
	_, _ = s.Raw("USE test_db").Exec()
	_, _ = s.Raw("DROP TABLE IF EXISTS user;").Exec()
	_, _ = s.Raw("CREATE TABLE user (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255));").Exec()
	_, _ = s.Raw("CREATE TABLE user (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255));").Exec()
	result, _ := s.Raw("INSERT INTO user (name) VALUES (?)", "test").Exec()
	count, _ := result.RowsAffected()
	fmt.Printf("Exec success, %d affected\n", count)

}

func origin() {
	dsn := "root:root008@@tcp(127.0.0.1:3306)/?charset=utf8mb4&parseTime=True"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.Exec("CREATE DATABASE IF NOT EXISTS test_db")
	db.Exec("USE test_db")
	result, err := db.Exec("CREATE TABLE IF NOT EXISTS user (id INT AUTO_INCREMENT PRIMARY KEY, name VARCHAR(255))")
	if err != nil {
		panic(err)
	}

	result, err = db.Exec("INSERT INTO user (name) VALUES ('test')")
	if err != nil {
		panic(err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}

	log.Println(affected)

	row := db.QueryRow("SELECT * FROM user")
	var id int
	var name string
	err = row.Scan(&id, &name)
	if err != nil {
		panic(err)
	}

	log.Println(id, name)
}
