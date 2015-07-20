package tools

import (
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/larspensjo/config"
	_ "github.com/mattn/go-sqlite3"
)

type KeyValue struct {
	Key   string
	Value int
}

// 统计信息 - 总量
type LogTypeCount struct {
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

// 全局变量定义
var (
	pwdUsedInfo string  = "作者:李振逢 QQ:24324962"
	configFile  string  = "config.ini" // 配置文件相对路径
	hDbRules    *sql.DB                // 规则数据库句柄
	hDbLog      *sql.DB                // 日志数据库句柄

)

func GetRootDir() (rootdir string, err error) {
	rootdir, err = filepath.Abs("./")
	if err != nil {
		return "", errors.New("错误:获取根路径失败")
	}
	return rootdir, nil
}

// 规则模块初始化
func RulesInit() (err error) {
	//连接数据库
	err = ConnectRuleDb()
	if err != nil {
		return err
	}

	//连接数据库
	err = ConnectLogDb()
	if err != nil {
		return err
	}

	// 初始化表
	/*{
		err = CreateTableUser(hDbRules)
		if err != nil {
			return err
		}

		err = CreateTableIPGroup(hDbRules)
		if err != nil {
			return err
		}

		err = CreateTableIPList(hDbRules)
		if err != nil {
			return err
		}

	}
	*/
	err = CreateTableLog(hDbLog)
	if err != nil {
		return err
	}
	return err
}

func ReleaseRuleDb() (err error) {
	CloseRuleDb()
	CloseLogDb()
	return nil
}

// 链接规则数据库
func ConnectRuleDb() (err error) {
	rootDir, err := GetRootDir()
	if err != nil {
		return err
	}

	configpath := filepath.Join(rootDir, configFile)
	cfgIni, err := config.ReadDefault(configpath)
	if err != nil {
		return errors.New("错误:读取配置文件失败:" + configpath)
	}

	dbName, err := cfgIni.String("Rules", "RuleDbFile")
	if err != nil {
		return errors.New("错误:[Rules]=>RuleDbFile失败")
	}

	dbPath := filepath.Join(rootDir, dbName)
	hDbRules, err = sql.Open("sqlite3", dbPath)

	return err
}

// 关闭数据库
func CloseRuleDb() {
	hDbRules.Close()
}

// 链接日志数据库
func ConnectLogDb() (err error) {
	rootDir, err := GetRootDir()
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
		return errors.New("错误:[Log]=>LogDbFile失败")
	}

	dbPath := filepath.Join(rootDir, dbName)
	hDbLog, err = sql.Open("sqlite3", dbPath)

	return err
}

// 关闭数据库
func CloseLogDb() {
	hDbLog.Close()
}

// 创建用户名密码表
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
		log.Printf("CreateTable(): %s, %s\n", err, sql)
		tx.Rollback()
		return err
	}

	sql = `insert into user (uid, uname, user_type, user_group, password) values 
		(2, 'Admin', 2, 2, 'b40fdc1791396dc11b4ad54b5744bcd6'),
		(3, 'Audit', 3, 3, 'b40fdc1791396dc11b4ad54b5744bcd6');
	`
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("InsertTable(): %s, %s\n", err, sql)
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

// 创建IP分组表
func CreateTableIPGroup(db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("CreateTableIPGroup:DB.Begin(): %s\n", err)
		return err
	}

	var sql string
	// 创建用户分组表
	sql = `create table if not exists ip_group (
			Id integer not null primary key autoincrement,
			Name char(128) unique not null
		);`
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("CreateTableIPGroup(): %s, %s\n", err, sql)
		tx.Rollback()
		return err
	}

	sql = `insert into ip_group (Id, Name) values (1, '默认组');`
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("CreateTableIPGroup(): %s, %s\n", err, sql)
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

// 创建IP列表
func CreateTableIPList(db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("CreateTableIPList:DB.Begin(): %s\n", err)
		return err
	}

	var sql string
	// 创建IP列表
	sql = `create table if not exists ip_list (
			Id integer not null primary key autoincrement,
			Gid integer not null default 1,
			IP    char(32) unique not null,
			Port  char(8) not null
		);`
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("CreateTableIPList(): %s, %s\n", err, sql)
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

// 创建日志表
func CreateTableLog(db *sql.DB) (err error) {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("CreateTableLog:DB.Begin(): %s\n", err)
		return err
	}

	var sql string
	// 统计表
	sql = `create table if not exists log_count (
			Id integer not null primary key, 
			IP  varchar(32) not null,
			Time char(20),
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
		log.Printf("CreateTableLog(): %s, %s\n", err, sql)
		tx.Rollback()
		return err
	}

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("CreateTableLog(commit transaction): %s\n", err)
		tx.Rollback()
		return err
	}
	return err
}

// 取字符串md5值，返回32字节字符串
func RulesGetMd5String(s string) (md5string string) {
	md5string = ""
	h := md5.New()
	h.Write([]byte(s))
	digest := h.Sum(nil)
	md5string = fmt.Sprintf("%x", digest)

	return md5string
}

// 用户密码验证
func RulesCheckUserPassword(user string, pwd string) (uid int, user_type int, err error) {
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RulesCheckUserPassword: %s\n", err)
		return uid, user_type, err
	}

	sql := fmt.Sprintf("select uid, user_type, password from user where uname = '%s'", user)
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("RulesCheckUserPassword(): %s", err)
		return uid, user_type, errors.New("错误:查询用户名密码失败")
	}
	defer rows.Close()

	var password string = ""
	for rows.Next() {
		rows.Scan(&uid, &user_type, &password)
		break
	}
	rows.Close()

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RulesCheckUserPassword(commit transaction): %s\n", err)
		tx.Rollback()
		return uid, user_type, err
	}

	if password == "" {
		return uid, user_type, errors.New("错误:用户不存在")
	}
	// 校验密码
	md5Pwd := RulesGetMd5String(pwdUsedInfo + pwd)
	if md5Pwd != password {
		return uid, user_type, errors.New("错误:密码不正确")
	}

	return uid, user_type, nil
}

// 修改用户密码
func RulesChangeUserPassword(user string, pwdold, pwdnew string) (err error) {
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RulesChangeUserPassword: %s\n", err)
		return err
	}

	sql := fmt.Sprintf("select password from user where uname = '%s'", user)
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("RulesChangeUserPassword(): %s", err)
		return errors.New("错误:查询用户密码失败")
	}
	defer rows.Close()

	var password string
	for rows.Next() {
		rows.Scan(&password)
		break
	}
	rows.Close()

	if password == "" {
		return errors.New("错误:用户不存在")
	}

	// 校验旧密码
	md5oldPwd := RulesGetMd5String(pwdUsedInfo + pwdold)
	if md5oldPwd != password {
		log.Printf("RulesChangeUserPassword(): %s", err)
		return errors.New("错误:旧密码错误")
	}

	// 更新新密码
	md5newPwd := RulesGetMd5String(pwdUsedInfo + pwdnew)
	sql = fmt.Sprintf("update user set password = '%s' where uname = '%s'", md5newPwd, user)
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("RulesChangeUserPassword(): %s", err)
		tx.Rollback()
		return errors.New("错误:更新密码失败")
	}

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RulesChangeUserPassword(commit transaction): %s\n", err)
		tx.Rollback()
		return err
	}

	return nil
}

// IP分组 - 添加
func RuleIPAddGroup(gname string) (err error) {
	if gname == "默认组" {
		return errors.New("错误:不允许添加默认组")
	}
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("IPAddGroup:DB.Begin(): %s\n", err)
		return err
	}

	var sql string
	// 添加IP组
	sql = fmt.Sprintf("insert into ip_group (Id, Name) values (null, '%s');", gname)
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("IPAddGroup(): %s, %s\n", err, sql)
		tx.Rollback()
		return err
	}

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("IPAddGroup(commit transaction): %s\n", err)
		tx.Rollback()
		return err
	}
	return err
}

// IP分组 - 删除
func RuleIPDelGroup(gname string) (err error) {
	if gname == "默认组" {
		return errors.New("错误:不允许删除默认组")
	}
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("IPDelGroup:DB.Begin(): %s\n", err)
		return err
	}

	var sql string
	// 将本组的IP移动到默认组
	sql = fmt.Sprintf("select A.IP from ip_list A left join ip_group B on A.Gid = B.Id where B.Name = '%s';", gname)
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("IPDelGroup(): %s", err)
		return errors.New("错误:删除分组失败")
	}
	defer rows.Close()

	var tip string
	for rows.Next() {
		rows.Scan(&tip)
		// 移动该IP到默认组
		sql = fmt.Sprintf("update ip_list set Gid = 1 where IP = '%s';", tip)
		_, err = tx.Exec(sql)
		if err != nil {
			log.Printf("IPDelGroup(): %s, %s\n", err, sql)
			tx.Rollback()
			return err
		}
	}
	rows.Close()

	// 删除IP组
	sql = fmt.Sprintf("delete from ip_group where Name = '%s';", gname)
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("IPDelGroup(): %s, %s\n", err, sql)
		tx.Rollback()
		return err
	}

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("IPDelGroup(commit transaction): %s\n", err)
		tx.Rollback()
		return err
	}
	return err
}

// IP分组 - 查询
func RuleIPQueryGroup() (groups []string, err error) {
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("IPQueryGroup:DB.Begin(): %s\n", err)
		return groups, err
	}

	var sql string
	// 查找IP组
	sql = fmt.Sprintf("select Name from ip_group;")
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("IPQueryGroup(): %s", err)
		return groups, errors.New("错误:查询IP分组失败")
	}
	defer rows.Close()

	var Name string
	for rows.Next() {
		rows.Scan(&Name)
		groups = append(groups, Name)
	}
	rows.Close()

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("IPQueryGroup(commit transaction): %s\n", err)
		tx.Rollback()
		return groups, err
	}
	return groups, err
}

// IP添加
func RuleIPAdd(ip, port, gname string) (err error) {
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RuleIPAdd:DB.Begin(): %s\n", err)
		return err
	}

	// 获取组ID
	var sql string
	sql = fmt.Sprintf("select Id from ip_group where Name = '%s'", gname)
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("IPQueryGroup(): %s", err)
		return errors.New("错误:添加IP失败")
	}
	defer rows.Close()

	var Gid int = 0
	for rows.Next() {
		rows.Scan(&Gid)
		break
	}
	rows.Close()

	if Gid == 0 {
		return errors.New("错误:添加IP失败,分组不存在")
	}

	// 查询IP是否存在
	sql = fmt.Sprintf("select count(*) from ip_list where IP = '%s'", ip)
	rows, err = db.Query(sql)
	if err != nil {
		log.Printf("IPQueryGroup(): %s", err)
		return errors.New("错误:添加IP失败")
	}
	defer rows.Close()

	var cnt int = 0
	for rows.Next() {
		rows.Scan(&cnt)
	}
	rows.Close()

	sql = ""
	if cnt == 0 {
		// IP不存在，添加
		sql = fmt.Sprintf("insert into ip_list (Id, Gid, IP, Port) values (null, %d, '%s', '%s');", Gid, ip, port)
	} else {
		// IP存在，更新
		sql = fmt.Sprintf("update ip_list set Gid = %d, Port = '%s' where IP = '%s';", Gid, port, ip)
	}
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("RuleIPAdd(): %s, %s\n", err, sql)
		tx.Rollback()
		return err
	}

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RuleIPAdd(commit transaction): %s\n", err)
		tx.Rollback()
		return err
	}
	return err
}

// IP删除
func RuleIPDel(ip string) (err error) {
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RuleIPDel:DB.Begin(): %s\n", err)
		return err
	}

	// 删除IP
	sql := fmt.Sprintf("delete from ip_list where IP = '%s';", ip)
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("RuleIPDel(): %s, %s\n", err, sql)
		tx.Rollback()
		return err
	}

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RuleIPDel(commit transaction): %s\n", err)
		tx.Rollback()
		return err
	}
	return err
}

type IpPort struct { // 组列表
	IP    string
	Port  string
	Gname string
}

// IP - 查询
func RuleIPQuery() (ipport []IpPort, err error) {
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RuleIPQuery:DB.Begin(): %s\n", err)
		return ipport, err
	}

	var sql string
	// 查找IP组
	sql = fmt.Sprintf("select A.IP, A.Port, B.Name from ip_list as A left join ip_group as B on A.Gid = B.Id order by B.Name;")
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("RuleIPQuery(): %s", err)
		return ipport, errors.New("错误:查询IP列表失败")
	}
	defer rows.Close()

	var ipp IpPort
	for rows.Next() {
		rows.Scan(&ipp.IP, &ipp.Port, &ipp.Gname)
		ipport = append(ipport, ipp)
	}
	rows.Close()

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RuleIPQuery(commit transaction): %s\n", err)
		tx.Rollback()
		return ipport, err
	}
	return ipport, err
}

// IP - 按组查询
func RuleIPQueryByGroup(group string) (ipport []IpPort, err error) {
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RuleIPQueryByGroup:DB.Begin(): %s\n", err)
		return ipport, err
	}

	var sql string
	// 查找IP组
	sql = fmt.Sprintf("select A.IP, A.Port from ip_list as A left join ip_group as B on A.Gid = B.Id where B.Name = '%s' order by A.Ip;", group)
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("RuleIPQueryByGroup(): %s", err)
		return ipport, errors.New("错误:查询IP列表失败")
	}
	defer rows.Close()

	var ipp IpPort
	for rows.Next() {
		rows.Scan(&ipp.IP, &ipp.Port)
		ipp.Gname = group
		ipport = append(ipport, ipp)
	}
	rows.Close()

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RuleIPQueryByGroup(commit transaction): %s\n", err)
		tx.Rollback()
		return ipport, err
	}
	return ipport, err
}

// 客户端登录时候将客户端IP+端口添加到管理中心
func RuleClientLogin(ip, port string) (err error) {
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RuleClientLogin:DB.Begin(): %s\n", err)
		return err
	}

	var sql string
	// 查找IP是否存在
	sql = fmt.Sprintf("select count(*) from ip_list where IP = '%s'", ip)
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("RuleClientLogin(): %s", err)
		return errors.New("错误:添加IP到管理中心失败")
	}
	defer rows.Close()

	var cnt int = 0
	for rows.Next() {
		rows.Scan(&cnt)
	}
	rows.Close()

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RuleClientLogin(commit transaction): %s\n", err)
		tx.Rollback()
		return err
	}

	if cnt == 0 {
		// ip不存在，插入
		err = RuleIPAdd(ip, port, "默认组")
	}
	return err
}

// 客户端将日志统计上传到管理中心
func RuleClientLogToday(IP, Time string, Totle, White, Black, BaseWinDir, BaseWinStart,
	BaseWinFormat, BaseWinProc, BaseWinService, HighAddService, HighAutoRun,
	HighAddStart, HighReadWrite, HighCreateExe, HighLoadSys, HighProcInject int) (err error) {
	db := hDbLog
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RuleClientLogToday:DB.Begin(): %s\n", err)
		return err
	}

	var sql string
	// 查找记录是否存在
	sql = fmt.Sprintf("select count(*) from log_count where IP = '%s' and Time = '%s'", IP, Time)
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("RuleClientLogToday(): %s", err)
		return errors.New("错误:接收客户端日志统计失败")
	}
	defer rows.Close()

	var cnt int = 0
	for rows.Next() {
		rows.Scan(&cnt)
	}
	rows.Close()

	sql = ""
	if cnt == 0 {
		// insert
		sql = fmt.Sprintf("insert into log_count (Id, IP, Time, Totle, White, Black, BaseWinDir, BaseWinStart, BaseWinFormat, BaseWinProc, BaseWinService, HighAddService, HighAutoRun, HighAddStart, HighReadWrite, HighCreateExe, HighLoadSys, HighProcInject) values (null, '%s', '%s', %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d)",
			IP, Time, Totle, White, Black, BaseWinDir, BaseWinStart, BaseWinFormat, BaseWinProc,
			BaseWinService, HighAddService, HighAutoRun, HighAddStart, HighReadWrite, HighCreateExe,
			HighLoadSys, HighProcInject)
	} else {
		// update
		sql = fmt.Sprintf("update log_count set Totle = %d, White = %d, Black = %d, BaseWinDir = %d, BaseWinStart = %d, BaseWinFormat = %d, BaseWinProc = %d, BaseWinService = %d, HighAddService = %d, HighAutoRun = %d, HighAddStart = %d, HighReadWrite = %d, HighCreateExe = %d, HighLoadSys = %d, HighProcInject = %d where IP = '%s' and Time = '%s'",
			Totle, White, Black, BaseWinDir, BaseWinStart, BaseWinFormat, BaseWinProc,
			BaseWinService, HighAddService, HighAutoRun, HighAddStart, HighReadWrite, HighCreateExe,
			HighLoadSys, HighProcInject, IP, Time)
	}

	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("LogInsertToday(): %s, %s\n", err, sql)
		tx.Rollback()
		return err
	}

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RuleClientLogToday(commit transaction): %s\n", err)
		tx.Rollback()
		return err
	}

	return err
}

// 查询当年 每个月的事件总数
func RuleQueryMonEventInYear() (moneveyear []KeyValue, err error) {
	db := hDbLog
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RuleQueryMonEventInYear:DB.Begin(): %s\n", err)
		return moneveyear, err
	}

	tm := time.Now()
	Year := tm.Year()
	var sql string

	// 查找IP组
	sql = fmt.Sprintf("select sum(Totle) as Tot, substr(Time,1,7) as Time from log_count where Time like '%d-%%' group by substr(Time,1,7) order by Time asc;", Year)
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("RuleQueryMonEventInYear(): %s", err)
		return moneveyear, errors.New("错误:查询当前每个月的事件总数失败")
	}
	defer rows.Close()

	var data KeyValue
	for rows.Next() {
		rows.Scan(&data.Value, &data.Key)
		moneveyear = append(moneveyear, data)
	}
	rows.Close()

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RuleQueryMonEventInYear(commit transaction): %s\n", err)
		tx.Rollback()
		return moneveyear, err
	}
	return moneveyear, err
}

// 查询 本月主机排名 Top10
func RuleQueryTopIPInMon() (topIPInMon []KeyValue, err error) {
	db := hDbLog
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RuleQueryTopIPInMon:DB.Begin(): %s\n", err)
		return topIPInMon, err
	}

	tm := time.Now()
	Year := tm.Year()
	Mon := int(tm.Month())
	var sql string

	// 查找IP组
	sql = fmt.Sprintf("select sum(Totle) as Totle, IP from log_count where Time like '%04d-%02d-%%' group by IP order by Totle desc limit 0, 10;", Year, Mon)
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("RuleQueryTopIPInMon(): %s", err)
		return topIPInMon, errors.New("错误:查询本月主机排名失败")
	}
	defer rows.Close()

	var data KeyValue
	for rows.Next() {
		rows.Scan(&data.Value, &data.Key)
		topIPInMon = append(topIPInMon, data)
	}
	rows.Close()

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RuleQueryTopIPInMon(commit transaction): %s\n", err)
		tx.Rollback()
		return topIPInMon, err
	}
	return topIPInMon, err
}

// 查询 本年主机排名 Top10
func RuleQueryTopIPInYear() (topIPInYear []KeyValue, err error) {
	db := hDbLog
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RuleQueryTopIPInYear:DB.Begin(): %s\n", err)
		return topIPInYear, err
	}

	tm := time.Now()
	Year := tm.Year()
	var sql string

	// 查找IP组
	sql = fmt.Sprintf("select sum(Totle) as Totle, IP from log_count where Time like '%04d-%%' group by IP order by Totle desc limit 0, 10;", Year)
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("RuleQueryTopIPInYear(): %s", err)
		return topIPInYear, errors.New("错误:查询本年主机排名失败")
	}
	defer rows.Close()

	var data KeyValue
	for rows.Next() {
		rows.Scan(&data.Value, &data.Key)
		topIPInYear = append(topIPInYear, data)
	}
	rows.Close()

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RuleQueryTopIPInYear(commit transaction): %s\n", err)
		tx.Rollback()
		return topIPInYear, err
	}
	return topIPInYear, err
}

// 本月安全事件统计
func RuleQueryTotEventInMon() (logtypecnt LogTypeCount, err error) {
	db := hDbLog
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RuleQueryTotEventInMon:DB.Begin(): %s\n", err)
		return logtypecnt, err
	}

	tm := time.Now()
	Year := tm.Year()
	Mon := int(tm.Month())
	var sql string

	// 查找IP组
	sql = "select sum(White), sum(Black), sum(BaseWinDir), sum(BaseWinStart), sum(BaseWinFormat)," +
		" sum(BaseWinProc), sum(BaseWinService), sum(HighAddService), sum(HighAutoRun)," +
		" sum(HighAddStart),  sum(HighReadWrite), sum(HighCreateExe), sum(HighLoadSys), sum(HighProcInject)" +
		fmt.Sprintf(" from log_count where Time like '%04d-%02d-%%';", Year, Mon)
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("RuleQueryTotEventInMon(): %s", err)
		return logtypecnt, errors.New("错误:查询本月安全事件统计失败")
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&logtypecnt.White, &logtypecnt.Black, &logtypecnt.BaseWinDir,
			&logtypecnt.BaseWinStart, &logtypecnt.BaseWinFormat, &logtypecnt.BaseWinProc,
			&logtypecnt.BaseWinService, &logtypecnt.HighAddService, &logtypecnt.HighAutoRun,
			&logtypecnt.HighAddStart, &logtypecnt.HighReadWrite, &logtypecnt.HighCreateExe,
			&logtypecnt.HighLoadSys, &logtypecnt.HighProcInject)
	}
	rows.Close()

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RuleQueryTotEventInMon(commit transaction): %s\n", err)
		tx.Rollback()
		return logtypecnt, err
	}
	return logtypecnt, err
}
