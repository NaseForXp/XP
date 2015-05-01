package xplog

//package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"sync"
	"time"

	"../RootDir"

	"github.com/larspensjo/config"
	_ "github.com/mattn/go-sqlite3"
)

// 常量
var (
	configFile     string         = "config.ini"    // 配置文件相对路径
	logFile        string         = "log\\log.db"   // 规则数据库文件
	hDBLog         *sql.DB                          // 数据库句柄
	rwLockLog      sync.RWMutex                     // 记录操作锁
	logWaitGroup   sync.WaitGroup                   // 线程等待组
	logCacheInsert []string                         // 拦截的日志放入这里
	logCacheWrite  []string                         // 写入数据库的实际缓存
	logWriteSpace  time.Duration  = time.Second * 1 // 日志写入间隔
	logStatus      int            = 0               // 线程退出标志
)

func LogInit() (err error) {
	err = LogConnectSqlite()
	if err != nil {
		return err
	}

	err = LogCreateTable()
	if err != nil {
		return err
	}

	logWaitGroup.Add(1)
	go LogWriteCacheToDb()
	return nil
}

func LogFini() {
	rwLockLog.Lock()
	logStatus = 1 // 通知日志写入线程退出
	rwLockLog.Unlock()
	logWaitGroup.Wait()
	LogCloseSqlite(hDBLog)
}

// 插入系统日志
func LogInsertSys(uname string, op string, info string, ret string) {
	sql := "insert into log_sys (Id, Uname, Op, Info, Result, Time) values " +
		fmt.Sprintf("(null, '%s', '%s', '%s', '%s', datetime())", uname, op, info, ret)

	rwLockLog.Lock()
	logCacheInsert = append(logCacheInsert, sql)
	rwLockLog.Unlock()
}

// 插入违规日志
func LogInsertEvent(Module string, Mode int, User, Sub, Obj, Op, Ret string) {
	sql := "insert into log_event (id, Module, Mode, User, Sub, Obj, Op, Ret, Time) values " +
		fmt.Sprintf("(null, '%s', %d, '%s', '%s', '%s', '%s', '%s', datetime())", Module, Mode, User, Sub, Obj, Op, Ret)

	rwLockLog.Lock()
	logCacheInsert = append(logCacheInsert, sql)
	rwLockLog.Unlock()
}

func LogWriteCacheToDbExe() (err error) {
	db := hDBLog
	tx, err := db.Begin()
	if err != nil {
		log.Printf("LogWriteCacheToDbExe:DB.Begin(): %s\n", err)
		return err
	}

	for _, sql := range logCacheWrite {
		_, err = tx.Exec(sql)
		if err != nil {
			log.Printf("LogWriteCacheToDbExe(user): %s, %s\n", err, sql)
			tx.Rollback()
			return err
		}
	}

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("LogWriteCacheToDbExe(commit transaction): %s\n", err)
		tx.Rollback()
		return err
	}

	return nil
}

// 线程 - 将缓存中的日志写入数据库文件
func LogWriteCacheToDb() {
	for {
		if logStatus == 1 {
			rwLockLog.Lock()
			logCacheWrite = logCacheInsert
			logCacheInsert = logCacheInsert[0:0]
			rwLockLog.Unlock()
			LogWriteCacheToDbExe()
			break
		}
		rwLockLog.Lock()
		logCacheWrite = logCacheInsert
		logCacheInsert = logCacheInsert[0:0]
		rwLockLog.Unlock()
		LogWriteCacheToDbExe()

		time.Sleep(logWriteSpace)
	}
	logWaitGroup.Done()
	fmt.Println("Write Process exists")
}

// 连接数据库
func LogConnectSqlite() (err error) {
	rootDir, err := RootDir.GetRootDir()
	if err != nil {
		return err
	}

	configpath := filepath.Join(rootDir, configFile)
	cfgIni, err := config.ReadDefault(configpath)
	if err != nil {
		return errors.New("错误:读取配置文件失败:" + configpath)
	}

	dbName, err := cfgIni.String("Log", "LogDbFile")
	if err != nil {
		return errors.New("错误:[Log]=>LogDbFile")
	}

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		log.Fatal(dbName, err)
	}

	hDBLog = db
	return err
}

func LogCloseSqlite(db *sql.DB) {
	db.Close()
}

func LogCreateTable() (err error) {
	db := hDBLog

	tx, err := db.Begin()
	if err != nil {
		log.Printf("CreateLogTable:DB.Begin(): %s\n", err)
		return err
	}

	var sql string
	// 操作记录表
	sql = `create table if not exists log_sys (
			Id integer not null primary key, 
			Uname  varchar(128) not null,
			Op     varchar(128) not null,
			Info   varchar(256) not null,
			Result varchar(12) not null,
			Time   datetime
		);`
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("CreateLogTable(user): %s, %s\n", err, sql)
		tx.Rollback()
		return err
	}

	// 违规记录表
	sql = `create table if not exists log_event (
			Id integer not null primary key, 
			Module varchar(32) not null,
			Mode   integer not null,
			User   varchar(128) not null,
			Sub    varchar(260) not null,
			Obj    varchar(260) not null,
			Op     varchar(64) not null,
			Ret    varchar(16) not null,
			Time   datetime
		);`
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("CreateLogTable(user): %s, %s\n", err, sql)
		tx.Rollback()
		return err
	}

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("CreateLogTable(commit transaction): %s\n", err)
		tx.Rollback()
		return err
	}
	return err
}

/*
func main() {
	err := LogInit()
	fmt.Println(err)
	if err != nil {
		return
	}

	err = LogCreateTable()
	fmt.Println(err)
	if err != nil {
		return
	}

	for i := 0; i < 1000; i++ {
		// 插入系统日志
		LogInsertSys("uname", "op", "info", "ret")

		// 插入违规日志
		LogInsertEvent("Module", 1, "User", "Sub", "Obj", "Op", "Ret")
	}

	time.Sleep(time.Second * 2)
	LogFini()
}
*/
