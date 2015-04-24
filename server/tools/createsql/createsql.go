package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// 连接数据库
func ConnectSqlite(dbName string) (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(dbName, err)
	}

	return db, err
}

func CloseSqlite(db *sql.DB) {
	db.Close()
}

func CreateTableUser(db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("CreateTable:DB.Begin(): %s\n", err)
		return err
	}

	var sql string
	// 创建用户分组表
	sql = `create table if not exists user (
			uid integer not null primary key, 
			uname char(128) unique,
			user_type  integer,
			user_group integer,
			password   char(128) not null
		);`
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("CreateTable(user): %s, %s\n", err, sql)
		tx.Rollback()
		return err
	}

	sql = `insert into user (uid, uname, user_type, user_group, password) values 
		(1, 'centerAdmin', 1, 1, '123456'),
		(2, 'Admin', 2, 2, '123456'),
		(3, 'Audit', 3, 3, '123456');
	`
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("InsertTable(user): %s, %s\n", err, sql)
		tx.Rollback()
		return err
	}

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("CreateTable(commit transaction): %s\n", err)
		tx.Rollback()
		return err
	}
	return err
}

func main() {
	dbpath := "./config.db"
	dbpath, _ = filepath.Abs(dbpath)
	fmt.Println(dbpath)

	os.Remove(dbpath)

	db, err := ConnectSqlite(dbpath)
	err = CreateTableUser(db)
	fmt.Println(err)
	CloseSqlite(db)
}
