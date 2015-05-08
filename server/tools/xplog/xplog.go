package xplog

//package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"
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

// 系统日志 - 查询 - 单条 (返回结果为数组)
type LogSysQueryRes struct {
	Uname  string // 用户名
	Op     string // 操作
	Info   string // 内容
	Result string // 结果
	Time   string // 时间
}

// 拦截日志 - 查询 - 单条 (返回结果为数组)
type LogEventQueryRes struct {
	Module string // 模块 (安全防护 | 增强防护)
	Mode   string // 模式 (防护模式 | 监控模式)
	User   string // 用户名
	Sub    string // 主体进程
	Obj    string // 对象
	Op     string // 操作
	Ret    string // 操作结果
	Time   string // 时间
}

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
func LogInsertEvent(Module, Mode, User, Sub, Obj, Op, Ret string) {
	sql := "insert into log_event (id, Module, Mode, User, Sub, Obj, Op, Ret, Time) values " +
		fmt.Sprintf("(null, '%s', '%s', '%s', '%s', '%s', '%s', '%s', datetime())", Module, Mode, User, Sub, Obj, Op, Ret)

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
		//fmt.Println(sql)
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
			Module varchar(128) not null,
			Mode   varchar(32) not null,
			User   varchar(128) not null,
			Sub    varchar(260) not null,
			Obj    varchar(520) not null,
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

// 查询系统日志总数
func LogQuerySysTotle() (totCount int, err error) {
	db := hDBLog
	tx, err := db.Begin()
	if err != nil {
		log.Printf("LogQuerySysTotle:DB.Begin(): %s\n", err)
		return totCount, err
	}

	sql := "select count(*) from log_sys"
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("LogQuerySysTotle(): %s", err)
		return totCount, errors.New("错误:查询系统日志总数失败")
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&totCount)
		break
	}
	rows.Close()

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("LogQuerySysTotle(commit transaction): %s\n", err)
		tx.Rollback()
		return totCount, err
	}

	return totCount, nil
}

// 查询安全日志总数
func LogQueryEventTotle() (totCount int, err error) {
	db := hDBLog
	tx, err := db.Begin()
	if err != nil {
		log.Printf("LogQueryEventTotle:DB.Begin(): %s\n", err)
		return totCount, err
	}

	sql := "select count(*) from log_event"
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("LogQueryEventTotle(): %s", err)
		return totCount, errors.New("错误:查询安全日志总数失败")
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&totCount)
		break
	}
	rows.Close()

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("LogQueryEventTotle(commit transaction): %s\n", err)
		tx.Rollback()
		return totCount, err
	}

	return totCount, nil
}

// 检测时间格式
func IsTimeRangeRight(timeStart string, timeEnd string) bool {
	var unix_stime int64
	var unix_etime int64

	if strings.EqualFold(timeStart, "") != true {
		s_time, err := time.Parse("2006-01-02 15:04:05", timeStart)
		if err == nil {
			unix_stime = s_time.Unix()
		} else {
			return false
		}
	} else {
		unix_stime = 0
	}
	if strings.EqualFold(timeEnd, "") != true {
		e_time, err1 := time.Parse("2006-01-02 15:04:05", timeEnd)
		if err1 == nil {
			unix_etime = e_time.Unix()
		} else {
			return false
		}
	} else {
		unix_etime = time.Now().Unix()
	}
	if unix_stime > unix_etime {
		return false
	} else {
		return true
	}
}

// 查询系统日志
func LogQuerySys(KeyWord, TimeStart, TimeStop string, Start, Length int) (resArray []LogSysQueryRes, err error) {
	if IsTimeRangeRight(TimeStart, TimeStop) == false {
		return resArray, errors.New("错误:查询时间格式不正确:[" + TimeStart + "~" + TimeStop + "]")
	}

	db := hDBLog
	tx, err := db.Begin()
	if err != nil {
		log.Printf("LogQuerySys:DB.Begin(): %s\n", err)
		return resArray, err
	}

	sql := "select Uname, Op, Info, Result, strftime('%Y-%m-%d %H:%M:%S', Time) from log_sys where " +
		"Time >= '" + TimeStart + "' and " +
		"Time <= '" + TimeStop + "' and " +
		"( Uname like '%" + KeyWord + "%' or " +
		"Op like '%" + KeyWord + "%' or " +
		"Info like '%" + KeyWord + "%' or " +
		"Result like '%" + KeyWord + "%' ) " +
		fmt.Sprintf("order by Id desc limit %d, %d ", Start, Length)

	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("LogQuerySys(): %s", err)
		return resArray, errors.New("错误:查询系统日志失败")
	}
	defer rows.Close()

	var res LogSysQueryRes
	for rows.Next() {
		rows.Scan(&res.Uname, &res.Op, &res.Info, &res.Result, &res.Time)
		resArray = append(resArray, res)
	}
	rows.Close()

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("LogQuerySys(commit transaction): %s\n", err)
		tx.Rollback()
		return resArray, err
	}

	return resArray, nil
}

// 查询安全日志
func LogQueryEvent(KeyWord, TimeStart, TimeStop string, Start, Length int) (resArray []LogEventQueryRes, err error) {
	if IsTimeRangeRight(TimeStart, TimeStop) == false {
		return resArray, errors.New("错误:查询时间格式不正确:[" + TimeStart + "~" + TimeStop + "]")
	}

	db := hDBLog
	tx, err := db.Begin()
	if err != nil {
		log.Printf("LogQueryEvent:DB.Begin(): %s\n", err)
		return resArray, err
	}

	sql := "select Module, Mode, User, Sub, Obj, Op, Ret, strftime('%Y-%m-%d %H:%M:%S', Time) from log_event where " +
		"Time >= '" + TimeStart + "' and " +
		"Time <= '" + TimeStop + "' and ( " +
		"Mode like '%" + KeyWord + "%' or " +
		"Module like '%" + KeyWord + "%' or " +
		"User like '%" + KeyWord + "%' or " +
		"Sub like '%" + KeyWord + "%' or " +
		"Obj like '%" + KeyWord + "%' or " +
		"Op like '%" + KeyWord + "%' or " +
		"Ret like '%" + KeyWord + "%' ) " +
		fmt.Sprintf("order by Id desc limit %d, %d ", Start, Length)

	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("LogQueryEvent(): %s", err)
		return resArray, errors.New("错误:查询系统日志失败")
	}
	defer rows.Close()

	var res LogEventQueryRes
	for rows.Next() {
		rows.Scan(&res.Module, &res.Mode, &res.User, &res.Sub, &res.Obj, &res.Op, &res.Ret, &res.Time)
		resArray = append(resArray, res)
	}
	rows.Close()

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("LogQueryEvent(commit transaction): %s\n", err)
		tx.Rollback()
		return resArray, err
	}

	return resArray, nil
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
