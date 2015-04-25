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

// 用户表
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

// 基本保护状态表
func CreateTableSafeBase(db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("CreateTable:DB.Begin(): %s\n", err)
		return err
	}

	var sql string

	sql = `create table if not exists safe_base (
			id integer not null primary key, 
			win_dir    integer default 0,
			win_start  integer default 0,
			win_format interger default 0,
			win_proc   interger default 0,
			win_sevice interger default 0
		);`
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("CreateTable(base_safe): %s, %s\n", err, sql)
		tx.Rollback()
		return err
	}

	sql = `insert into safe_base (id, win_dir, win_start, win_format, win_proc, win_sevice) values 
		(1, 0, 0, 0, 0, 0);	`
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

// 增强保护状态表
func CreateTableSafeHigh(db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("CreateTableSafeHigh:DB.Begin(): %s\n", err)
		return err
	}

	var sql string

	sql = `create table if not exists safe_high (
			id integer not null primary key, 
			add_sevice    integer default 0,
			auto_run      integer default 0,
			add_start     interger default 0,
			readwrite     interger default 0,
			create_exe    interger default 0,
			load_sys      interger default 0,
			proc_inject   interger default 0
		);`
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("CreateTableSafeHigh: %s, %s\n", err, sql)
		tx.Rollback()
		return err
	}

	sql = `insert into safe_high (id, add_sevice, auto_run, add_start, readwrite, create_exe, load_sys, proc_inject) values 
		(1, 0, 0, 0, 0, 0, 0, 0);	`
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("CreateTableSafeHigh: %s, %s\n", err, sql)
		tx.Rollback()
		return err
	}

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("CreateTableSafeHigh(commit transaction): %s\n", err)
		tx.Rollback()
		return err
	}
	return err
}

// 账户保护状态表
func CreateTableSafeAccount(db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("CreateTable:DB.Begin(): %s\n", err)
		return err
	}

	var sql string

	sql = `create table if not exists safe_account (
			id integer not null primary key, 
			safe_lev      integer default 0,
			pwd_complex   integer default 0,
			pwd_min_len   integer default 0,
			pwd_lock_time interger default 0,
			pwd_used_min  interger default 0,
			pwd_used_max  interger default 0,
			pwd_old_num   interger default 0
		);`
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("CreateTable(base_safe): %s, %s\n", err, sql)
		tx.Rollback()
		return err
	}

	sql = `insert into safe_account (id, safe_lev, pwd_complex, pwd_min_len, pwd_lock_time, pwd_used_min, pwd_used_max, pwd_old_num) values 
		(1, 0, 0, 0, 0, 0, 0, 0);	`
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

// 系统文件及目录
func CreateTableWinDir(db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("CreateTable:DB.Begin(): %s\n", err)
		return err
	}

	var sql string

	sql = `create table if not exists win_dir (
			id integer not null primary key autoincrement, 
			path      varchar(260) not null,
			perm      varchar(8)
		);`
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("CreateTable(base_safe): %s, %s\n", err, sql)
		tx.Rollback()
		return err
	}

	sql = `insert into win_dir (id, path, perm) values 
		(null, 'C:\\Windows\\*', 'r');`
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

// 系统启动文件
func CreateTableWinStart(db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("CreateTable:DB.Begin(): %s\n", err)
		return err
	}

	var sql string

	sql = `create table if not exists win_start (
			id integer not null primary key autoincrement, 
			path      varchar(260) not null,
			perm      varchar(8)
		);`
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("CreateTable(base_safe): %s, %s\n", err, sql)
		tx.Rollback()
		return err
	}

	sql = `insert into win_start (id, path, perm) values 
		(null, 'C:\\boot.ini', 'r'),
		(null, 'C:\\Ntldr', 'r');`
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

// 系统关键进程不被结束
func CreateTableWinProc(db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("CreateTable:DB.Begin(): %s\n", err)
		return err
	}

	var sql string

	sql = `create table if not exists win_proc (
			id integer not null primary key autoincrement, 
			path      varchar(260) not null
		);`
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("CreateTable(): %s, %s\n", err, sql)
		tx.Rollback()
		return err
	}

	sql = `insert into win_proc (id, path) values 
		(null, 'C:\\Windows\\System32\\csrss.exe'),
		(null, 'C:\\Windows\\System32\\lsass.exe'),
		(null, 'C:\\Windows\\System32\\services.exe'),
		(null, 'C:\\Windows\\System32\\smss.exe'),
		(null, 'C:\\Windows\\System32\\svchost.exe'),
		(null, 'C:\\Windows\\System32\\winlogon.exe');`
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

// 白名单
func CreateTableWhiteList(db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("CreateTable:DB.Begin(): %s\n", err)
		return err
	}

	var sql string

	sql = `create table if not exists whitelist (
			id integer not null primary key autoincrement, 
			path      varchar(260) not null
		);`
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("CreateTable(): %s, %s\n", err, sql)
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

// 黑名单
func CreateTableBlackList(db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("CreateTable:DB.Begin(): %s\n", err)
		return err
	}

	var sql string

	sql = `create table if not exists blacklist (
			id integer not null primary key autoincrement, 
			path      varchar(260) not null
		);`
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("CreateTable(): %s, %s\n", err, sql)
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

	err = CreateTableSafeBase(db)
	fmt.Println(err)

	err = CreateTableSafeHigh(db)
	fmt.Println(err)

	err = CreateTableSafeAccount(db)
	fmt.Println(err)

	err = CreateTableWinDir(db)
	fmt.Println(err)

	err = CreateTableWinStart(db)
	fmt.Println(err)

	err = CreateTableWinProc(db)
	fmt.Println(err)

	err = CreateTableWhiteList(db)
	fmt.Println(err)

	err = CreateTableBlackList(db)
	fmt.Println(err)

	CloseSqlite(db)
}
