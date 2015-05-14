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
		(1, 'CenterAdmin', 1, 1, 'bb149ef481514784aa75833d76be7b39'),
		(2, 'Admin', 2, 2, 'b40fdc1791396dc11b4ad54b5744bcd6'),
		(3, 'Audit', 3, 3, 'b40fdc1791396dc11b4ad54b5744bcd6');
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
			Mode       integer default 0,
			WinDir     integer default 0,
			WinStart   integer default 0,
			WinFormat  integer default 0,
			WinProc    integer default 0,
			WinService integer default 0
		);`
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("CreateTable(base_safe): %s, %s\n", err, sql)
		tx.Rollback()
		return err
	}

	sql = `insert into safe_base (id, Mode, WinDir, WinStart, WinFormat, WinProc, WinService) values 
		(1, 0, 0, 0, 0, 0, 0);	`
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
			Mode         integer default 0,
			AddService   integer default 0,
			AutoRun      integer default 0,
			AddStart     integer default 0,
			ReadWrite    integer default 0,
			CreateExe    integer default 0,
			LoadSys      integer default 0,
			ProcInject   integer default 0
		);`
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("CreateTableSafeHigh: %s, %s\n", err, sql)
		tx.Rollback()
		return err
	}

	sql = `insert into safe_high (id, Mode, AddService, AutoRun, AddStart, ReadWrite, CreateExe, LoadSys, ProcInject) values 
		(1, 0, 0, 0, 0, 0, 0, 0, 0);	`
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
			Mode          integer default 0,
			SafeLev       integer default 0
		);`
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("CreateTable(base_safe): %s, %s\n", err, sql)
		tx.Rollback()
		return err
	}

	sql = `insert into safe_account (id, Mode, SafeLev) values 
		(1, 0, 0);`
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
			path      varchar(260) not null unique,
			perm      varchar(8)
		);`
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("CreateTable(base_safe): %s, %s\n", err, sql)
		tx.Rollback()
		return err
	}

	sql = `insert into win_dir (id, path, perm) values 
		(null, 'C:\\Windows\\', 'r');`
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
			path      varchar(260) not null unique,
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
			path      varchar(260) not null unique
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

// 增强防护 - 开机启动项
func CreateTableHighWinStart(db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("CreateTable:DB.Begin(): %s\n", err)
		return err
	}

	var sql string

	sql = `create table if not exists high_winstart (
			id integer not null primary key autoincrement, 
			path      varchar(260) not null unique
		);`
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("CreateTable(high_winstart): %s, %s\n", err, sql)
		tx.Rollback()
		return err
	}

	sql = `insert into high_winstart (id, path) values 
		(null, 'HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\Run'),
		(null, 'HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\RunOnce')
		;`
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("InsertTable(high_winstart): %s, %s\n", err, sql)
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
			path      varchar(260) not null unique
		);`
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("CreateTable(): %s, %s\n", err, sql)
		tx.Rollback()
		return err
	}

	sql = `insert into whitelist (id, path) values 
		(null, 'c:\\windows\\pchealth\\helpctr\\binaries\\msconfig.exe');`
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("InsertTable(whitelist): %s, %s\n", err, sql)
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
			path      varchar(260) not null unique
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

	err = CreateTableHighWinStart(db)
	fmt.Println(err)

	err = CreateTableWhiteList(db)
	fmt.Println(err)

	err = CreateTableBlackList(db)
	fmt.Println(err)

	CloseSqlite(db)
}
