package xplog

//package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
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
	rwLockLog      sync.Mutex                       // 记录操作锁
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

// 客户端首页统计信息 - 总量
type LogHomeCount struct {
	Totle          int // 总数
	White          int // 白名单事件数量
	Black          int // 黑名单事件数量
	BaseWinDir     int // 基本防护-系统文件及目录保护
	BaseWinStart   int // 基本防护-系统启动文件保护
	BaseWinFormat  int // 基本防护-防止格式化系统磁盘
	BaseWinProc    int // 基本防护-防止系统关键进程被杀死
	BaseWinService int // 基本防护-防止篡改系统服务
	HighAddService int // 增强防护-防止服务被添加
	HighAutoRun    int // 增强防护-防止自动运行
	HighAddStart   int // 增强防护-防止开机自启动
	HighReadWrite  int // 增强防护-防止磁盘被直接读写
	HighCreateExe  int // 增强防护-禁止创建.exe文件
	HighLoadSys    int // 增强防护-防止驱动程序被加载
	HighProcInject int // 增强防护-防止进程被注入
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

	logStatus = 0
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

func LogGetDbName() (dbName string, err error) {
	rootDir, err := RootDir.GetRootDir()
	if err != nil {
		return dbName, err
	}

	configpath := filepath.Join(rootDir, configFile)
	cfgIni, err := config.ReadDefault(configpath)
	if err != nil {
		return dbName, errors.New("错误:读取配置文件失败:" + configpath)
	}

	dbName, err = cfgIni.String("Log", "LogDbFile")
	if err != nil {
		return dbName, errors.New("错误:[Log]=>LogDbFile")
	}
	return dbName, nil
}

// 连接数据库
func LogConnectSqlite() (err error) {
	dbName, err := LogGetDbName()
	if err != nil {
		return err
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
		log.Printf("CreateLogTable(log_sys): %s, %s\n", err, sql)
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
		log.Printf("CreateLogTable(log_event): %s, %s\n", err, sql)
		tx.Rollback()
		return err
	}

	// 统计表
	sql = `create table if not exists log_count (
			Id integer not null primary key, 
			Time date unique,
			Totle integer default 0,
			White integer default 0,
			Black integer default 0,
			BaseWinDir integer default 0,
			BaseWinStart integer default 0,
			BaseWinFormat integer default 0,
			BaseWinProc integer default 0,
			BaseWinService integer default 0,
			HighAddService integer default 0,
			HighAutoRun integer default 0,
			HighAddStart integer default 0,
			HighReadWrite integer default 0,
			HighCreateExe integer default 0,
			HighLoadSys integer default 0,
			HighProcInject integer default 0
		);`
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("CreateLogTable(log_count): %s, %s\n", err, sql)
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

// 客户端首页 - 查询统计信息 - 总量
func LogQueryHomeCount() (homeCnt LogHomeCount, err error) {
	db := hDBLog
	tx, err := db.Begin()
	if err != nil {
		log.Printf("LogQueryCount:DB.Begin(): %s\n", err)
		return homeCnt, err
	}

	// 查询总数totle
	sql := "select count(*) from log_event"
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("LogQueryCount(): %s", err)
		return homeCnt, errors.New("错误:查询首页统计信息失败")
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&homeCnt.Totle)
		break
	}
	rows.Close()

	// 查询单项
	sql = "select Module, count(*) as cnt from log_event group by Module"
	rows, err = db.Query(sql)
	if err != nil {
		log.Printf("LogQueryCount(): %s", err)
		return homeCnt, errors.New("错误:查询首页统计信息失败")
	}
	defer rows.Close()

	var name string
	var cnt int
	for rows.Next() {
		rows.Scan(&name, &cnt)
		switch name {
		case "白名单":
			homeCnt.White = cnt
		case "黑名单":
			homeCnt.Black = cnt
		case "基本防护-系统文件及目录保护":
			homeCnt.BaseWinDir = cnt
		case "基本防护-系统启动文件保护":
			homeCnt.BaseWinStart = cnt
		case "基本防护-防止格式化系统磁盘":
			homeCnt.BaseWinFormat = cnt
		case "基本防护-防止系统关键进程被杀死":
			homeCnt.BaseWinProc = cnt
		case "基本防护-防止篡改系统服务":
			homeCnt.BaseWinService = cnt
		case "增强防护-防止服务被添加":
			homeCnt.HighAddService = cnt
		case "增强防护-防止自动运行":
			homeCnt.HighAutoRun = cnt
		case "增强防护-防止开机自启动":
			homeCnt.HighAddStart = cnt
		case "增强防护-防止磁盘被直接读写":
			homeCnt.HighReadWrite = cnt
		case "增强防护-禁止创建.exe文件":
			homeCnt.HighCreateExe = cnt
		case "增强防护-防止驱动程序被加载":
			homeCnt.HighLoadSys = cnt
		case "增强防护-防止进程被注入":
			homeCnt.HighProcInject = cnt
		}
	}
	rows.Close()

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("LogQueryCount(commit transaction): %s\n", err)
		tx.Rollback()
		return homeCnt, err
	}

	return homeCnt, nil
}

func LogQueryDayInMonth() (dayInMon map[string]int, err error) {
	dayInMon = make(map[string]int)
	db := hDBLog
	tx, err := db.Begin()
	if err != nil {
		log.Printf("LogQueryDayInMonth:DB.Begin(): %s\n", err)
		return dayInMon, err
	}

	tm := time.Now()
	Year := int(tm.Year())
	Mon := int(tm.Month())
	Day := 0

	// 查询当月中每天的总数
	for Day = 1; Day <= 31; Day++ {
		sql := fmt.Sprintf("select count(*) from log_event where Time like '%04d-%02d-%02d%%';", Year, Mon, Day)
		rows, err := db.Query(sql)
		defer rows.Close()
		if err != nil {
			log.Printf("LogQueryDayInMonth: %s\n", err)
			return dayInMon, errors.New("错误:查询当月每天总数失败")
		}
		var cnt int
		for rows.Next() {
			rows.Scan(&cnt)
			dayInMon[fmt.Sprintf("%02d", Day)] = cnt
			break
		}
		rows.Close()
	}

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("LogQueryDayInMonth(commit transaction): %s\n", err)
		tx.Rollback()
		return dayInMon, err
	}
	return dayInMon, err
}

func LogQueryMonthEventTot() (MonTop map[string]int, err error) {
	// 初始化
	MonTop = make(map[string]int)
	MonTop["白名单"] = 0
	MonTop["黑名单"] = 0
	MonTop["系统文件及目录保护"] = 0
	MonTop["系统启动文件保护"] = 0
	MonTop["防止格式化系统磁盘"] = 0
	MonTop["防止系统关键进程被杀死"] = 0
	MonTop["防止篡改系统服务"] = 0
	MonTop["防止服务被添加"] = 0
	MonTop["防止自动运行"] = 0
	MonTop["防止开机自启动"] = 0
	MonTop["防止磁盘被直接读写"] = 0
	MonTop["禁止创建exe文件"] = 0
	MonTop["防止驱动程序被加载"] = 0
	MonTop["防止进程被注入"] = 0

	db := hDBLog
	tx, err := db.Begin()
	if err != nil {
		log.Printf("LogQueryMonthTop:DB.Begin(): %s\n", err)
		return MonTop, err
	}

	tm := time.Now()
	Year := int(tm.Year())
	Mon := int(tm.Month())

	// 查询当月分类统计TOP
	sql := fmt.Sprintf("select count(*) as cnt, Module from log_event where Time like '%04d-%02d%%' group by Module;", Year, Mon)
	rows, err := db.Query(sql)
	defer rows.Close()
	if err != nil {
		log.Printf("LogQueryMonthTop: %s\n", err)
		return MonTop, errors.New("错误:查询当月分类统计失败")
	}
	var cnt int
	var name string
	for rows.Next() {
		rows.Scan(&cnt, &name)
		switch name {
		case "白名单":
			MonTop["白名单"] = cnt
		case "黑名单":
			MonTop["黑名单"] = cnt
		case "基本防护-系统文件及目录保护":
			MonTop["系统文件及目录保护"] = cnt
		case "基本防护-系统启动文件保护":
			MonTop["系统启动文件保护"] = cnt
		case "基本防护-防止格式化系统磁盘":
			MonTop["防止格式化系统磁盘"] = cnt
		case "基本防护-防止系统关键进程被杀死":
			MonTop["防止系统关键进程被杀死"] = cnt
		case "基本防护-防止篡改系统服务":
			MonTop["防止篡改系统服务"] = cnt
		case "增强防护-防止服务被添加":
			MonTop["防止服务被添加"] = cnt
		case "增强防护-防止自动运行":
			MonTop["防止自动运行"] = cnt
		case "增强防护-防止开机自启动":
			MonTop["防止开机自启动"] = cnt
		case "增强防护-防止磁盘被直接读写":
			MonTop["防止磁盘被直接读写"] = cnt
		case "增强防护-禁止创建.exe文件":
			MonTop["禁止创建exe文件"] = cnt
		case "增强防护-防止驱动程序被加载":
			MonTop["防止驱动程序被加载"] = cnt
		case "增强防护-防止进程被注入":
			MonTop["防止进程被注入"] = cnt
		}
	}

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("LogQueryMonthTop(commit transaction): %s\n", err)
		tx.Rollback()
		return MonTop, err
	}
	return MonTop, err
}

func LogQueryYearEventTot() (YearTop map[string]int, err error) {
	// 初始化
	YearTop = make(map[string]int)
	YearTop["白名单"] = 0
	YearTop["黑名单"] = 0
	YearTop["系统文件及目录保护"] = 0
	YearTop["系统启动文件保护"] = 0
	YearTop["防止格式化系统磁盘"] = 0
	YearTop["防止系统关键进程被杀死"] = 0
	YearTop["防止篡改系统服务"] = 0
	YearTop["防止服务被添加"] = 0
	YearTop["防止自动运行"] = 0
	YearTop["防止开机自启动"] = 0
	YearTop["防止磁盘被直接读写"] = 0
	YearTop["禁止创建exe文件"] = 0
	YearTop["防止驱动程序被加载"] = 0
	YearTop["防止进程被注入"] = 0

	db := hDBLog
	tx, err := db.Begin()
	if err != nil {
		log.Printf("LogQueryYearEventTot:DB.Begin(): %s\n", err)
		return YearTop, err
	}

	tm := time.Now()
	Year := int(tm.Year())

	// 查询当月分类统计TOP
	sql := fmt.Sprintf("select count(*) as cnt, Module from log_event where Time like '%04d-%%' group by Module;", Year)
	rows, err := db.Query(sql)
	defer rows.Close()
	if err != nil {
		log.Printf("LogQueryYearEventTot: %s\n", err)
		return YearTop, errors.New("错误:查询当年分类统计失败")
	}
	var cnt int
	var name string
	for rows.Next() {
		rows.Scan(&cnt, &name)
		switch name {
		case "白名单":
			YearTop["白名单"] = cnt
		case "黑名单":
			YearTop["黑名单"] = cnt
		case "基本防护-系统文件及目录保护":
			YearTop["系统文件及目录保护"] = cnt
		case "基本防护-系统启动文件保护":
			YearTop["系统启动文件保护"] = cnt
		case "基本防护-防止格式化系统磁盘":
			YearTop["防止格式化系统磁盘"] = cnt
		case "基本防护-防止系统关键进程被杀死":
			YearTop["防止系统关键进程被杀死"] = cnt
		case "基本防护-防止篡改系统服务":
			YearTop["防止篡改系统服务"] = cnt
		case "增强防护-防止服务被添加":
			YearTop["防止服务被添加"] = cnt
		case "增强防护-防止自动运行":
			YearTop["防止自动运行"] = cnt
		case "增强防护-防止开机自启动":
			YearTop["防止开机自启动"] = cnt
		case "增强防护-防止磁盘被直接读写":
			YearTop["防止磁盘被直接读写"] = cnt
		case "增强防护-禁止创建.exe文件":
			YearTop["禁止创建exe文件"] = cnt
		case "增强防护-防止驱动程序被加载":
			YearTop["防止驱动程序被加载"] = cnt
		case "增强防护-防止进程被注入":
			YearTop["防止进程被注入"] = cnt
		}
	}

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("LogQueryYearEventTot(commit transaction): %s\n", err)
		tx.Rollback()
		return YearTop, err
	}
	return YearTop, err
}

// 拷贝文件
func CopyFile(srcName, dstName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return written, err
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return written, err
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

// 清空日志表
func LogClearSysEvent() (err error) {
	db := hDBLog
	tx, err := db.Begin()
	if err != nil {
		log.Printf("LogClearSysEvent:DB.Begin(): %s\n", err)
		return err
	}

	sql := "delete from log_sys"
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("LogClearSysEvent(log_sys): %s, %s\n", err, sql)
		tx.Rollback()
		return err
	}

	sql = "delete from log_event"
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("LogClearSysEvent(log_event): %s, %s\n", err, sql)
		tx.Rollback()
		return err
	}

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("LogQuerySysTotle(commit transaction): %s\n", err)
		tx.Rollback()
		return err
	}

	sql = "vacuum;"
	_, err = db.Exec(sql)
	if err != nil {
		log.Printf("LogClearSysEvent(log_sys): %s, %s\n", err, sql)
		return err
	}

	return nil
}

// 日志导出
func LogExport(SaveDir string) (SaveFile string, err error) {
	SaveDir, _ = filepath.Abs(SaveDir)

	tm := time.Now()
	name := fmt.Sprintf("%04d%02d%02d_%02d%02d%02d.db", tm.Year(), int(tm.Month()), tm.Day(), tm.Hour(), tm.Minute(), tm.Second())

	SaveFile = filepath.Join(SaveDir, name)

	// 1、关闭数据库
	rwLockLog.Lock()
	logStatus = 1 // 通知日志写入线程退出
	rwLockLog.Unlock()
	logWaitGroup.Wait()
	LogCloseSqlite(hDBLog)

	// 2、copy出数据库文件
	dbname, err := LogGetDbName()
	if err != nil {
		return SaveFile, err
	}

	_, err = CopyFile(dbname, SaveFile)
	if err != nil {
		return SaveFile, err
	}

	// 3、打开数据库
	err = LogInit()
	if err != nil {
		return SaveFile, err
	}

	// 4、清空数据库表项目
	err = LogClearSysEvent()
	if err != nil {
		return SaveFile, err
	}
	return SaveFile, nil
}

// 当天统计信息 - 总量
func LogQueryTodayCount() (homeCnt LogHomeCount, err error) {
	tm := time.Now()
	today := fmt.Sprintf("%04d-%02d-%02d", tm.Year(), int(tm.Month()), tm.Day())

	db := hDBLog
	tx, err := db.Begin()
	if err != nil {
		log.Printf("LogQueryTodayCount:DB.Begin(): %s\n", err)
		return homeCnt, err
	}

	// 查询总数totle
	sql := fmt.Sprintf("select count(*) from log_event where Time like '%s%%'", today)
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("LogQueryTodayCount(): %s", err)
		return homeCnt, errors.New("错误:查询当天统计信息失败")
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&homeCnt.Totle)
		break
	}
	rows.Close()

	// 查询单项
	sql = fmt.Sprintf("select Module, count(*) as cnt from log_event where Time like '%s%%' group by Module", today)
	rows, err = db.Query(sql)
	if err != nil {
		log.Printf("LogQueryTodayCount(): %s", err)
		return homeCnt, errors.New("错误:查询当天统计信息失败")
	}
	defer rows.Close()

	var name string
	var cnt int
	for rows.Next() {
		rows.Scan(&name, &cnt)
		switch name {
		case "白名单":
			homeCnt.White = cnt
		case "黑名单":
			homeCnt.Black = cnt
		case "基本防护-系统文件及目录保护":
			homeCnt.BaseWinDir = cnt
		case "基本防护-系统启动文件保护":
			homeCnt.BaseWinStart = cnt
		case "基本防护-防止格式化系统磁盘":
			homeCnt.BaseWinFormat = cnt
		case "基本防护-防止系统关键进程被杀死":
			homeCnt.BaseWinProc = cnt
		case "基本防护-防止篡改系统服务":
			homeCnt.BaseWinService = cnt
		case "增强防护-防止服务被添加":
			homeCnt.HighAddService = cnt
		case "增强防护-防止自动运行":
			homeCnt.HighAutoRun = cnt
		case "增强防护-防止开机自启动":
			homeCnt.HighAddStart = cnt
		case "增强防护-防止磁盘被直接读写":
			homeCnt.HighReadWrite = cnt
		case "增强防护-禁止创建.exe文件":
			homeCnt.HighCreateExe = cnt
		case "增强防护-防止驱动程序被加载":
			homeCnt.HighLoadSys = cnt
		case "增强防护-防止进程被注入":
			homeCnt.HighProcInject = cnt
		}
	}
	rows.Close()

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("LogQueryTodayCount(commit transaction): %s\n", err)
		tx.Rollback()
		return homeCnt, err
	}

	return homeCnt, nil
}

// 写入当天统计信息
func LogInsertToday(data LogHomeCount, today string) (err error) {
	return nil
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
