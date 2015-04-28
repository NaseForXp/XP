package rules

import (
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"../RootDir"

	"github.com/larspensjo/config"
	_ "github.com/mattn/go-sqlite3"
)

// 全局变量定义
var (
	rwLockRule sync.RWMutex                     // 全局读写锁 - 内存中的规则
	hMemRules  RuleMemHandle                    // 规则内存句柄
	configFile string        = "config.ini"     // 配置文件相对路径
	ruleDB     string        = "rule\\rules.db" // 规则数据库文件
	hDbRules   *sql.DB                          // 规则数据库句柄
)

// 安全防护 - 基本防护 - 配置
type SafeBaseConfig struct {
	Mode       int // 模式：0:监视模式 1:防护模式
	WinDir     int // 系统文件及目录保护状态 0:关闭 1:开启
	WinStart   int // 系统启动文件保护状态   0:关闭 1:开启
	WinFormat  int // 防止格式化磁盘状态    0:关闭 1:开启
	WinProc    int // 防止系统关键进程被杀死 0:关闭 1:开启
	WinService int // 防止篡改系统服务      0:关闭 1:开启
}

// 安全防护 - 增强防护 - 配置
type SafeHighConfig struct {
	Mode       int // 模式：0:监视模式 1:防护模式
	AddService int // 防止服务被添加       0:关闭 1:开启
	AutoRun    int // 防止自动运行恶意程序  0:关闭 1:开启
	AddStart   int // 防止添加开机启动项    0:关闭 1:开启
	ReadWrite  int // 防止磁盘直接读写      0:关闭 1:开启
	CreateExe  int // 防止创建EXE文件      0:关闭 1:开启
	LoadSys    int // 防止驱动被加载        0:关闭 1:开启
	ProcInject int // 防止进程被注入        0:关闭 1:开启
}

// 账户安全 - 配置
type AccountConfig struct {
	Mode          int // 模式：0:关闭 1:开启
	SafeLev       int // 账户策略设置 0:自定义 1:低级 2:中级 3:高级
	PwdComplex    int // 密码复杂度  0:关闭 1:开启
	PwdMinLen     int // 密码最小长度(字符个数)
	PwdUsedMin    int // 最短使用期限(天)
	PwdUsedMax    int // 最长使用期限(天)
	PwdOldNum     int // 强制密码历史次数(次)
	AccountTimes  int // 账户锁定次数(无效登录次数)
	AccountMinute int // 账户锁定时长(分钟)
}

// 模块初始化
func RulesInit() (err error) {
	//连接数据库
	err = RulesConnectDb()
	if err != nil {
		return err
	}

	// 规则加载到内存
	err = RulesMemInit()
	if err != nil {
		return err
	}
	return err
}

func RulesRelease() (err error) {
	CloseSqlite()
	return nil
}

// 链接数据库
func RulesConnectDb() (err error) {
	rootDir, err := RootDir.GetRootDir()
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
func CloseSqlite() {
	hDbRules.Close()
}

var pwdUsedInfo string = "作者:李振逢 QQ:24324962"

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

// 查找是否在白名单
func RuleCheckInWhite(fpath string) (bFind bool, err error) {
	bFind = false
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RuleCheckInWhite: %s\n", err)
		return bFind, err
	}

	absPath, _ := filepath.Abs(fpath)
	sql := fmt.Sprintf("select * from whitelist where path = '%s'", absPath)
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("RuleCheckInWhite(): %s", err)
		return bFind, errors.New("错误:查找白名单数据失败")
	}
	defer rows.Close()

	for rows.Next() {
		bFind = true
		break
	}
	rows.Close()

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RuleCheckInWhite(commit transaction): %s\n", err)
		tx.Rollback()
		return bFind, err
	}

	return bFind, nil
}

// 添加白名单(可能是目录或末尾带*)
func RulesAddWhite(fpath string) (err error) {
	bFind, err := RuleCheckInBlack(fpath)
	if err != nil {
		return err
	}

	if bFind == true {
		return errors.New("错误:该记录已经存在于黑名单中")
	}

	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RulesAddWhite: %s\n", err)
		return err
	}

	absPath, _ := filepath.Abs(fpath)
	sql := fmt.Sprintf("insert into whitelist (id, path) values (null, '%s')", absPath)
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("RulesAddWhite(): %s", err)
		tx.Rollback()
		return errors.New("错误:添加白名单失败")
	}

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RulesAddWhite(commit transaction): %s\n", err)
		tx.Rollback()
		return err
	}

	// 更新内存
	rwLockRule.Lock()
	hMemRules.White[absPath] = 0
	rwLockRule.Unlock()
	return nil
}

// 删除白名单(可能是目录或末尾带*)
func RulesDelWhite(fpath string) (err error) {
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RulesDelWhite: %s\n", err)
		return err
	}

	absPath, _ := filepath.Abs(fpath)
	sql := fmt.Sprintf("delete from whitelist where path = '%s'", absPath)
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("RulesDelWhite(): %s", err)
		tx.Rollback()
		return errors.New("错误:删除白名单失败")
	}

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RulesDelWhite(commit transaction): %s\n", err)
		tx.Rollback()
		return err
	}

	// 更新内存
	rwLockRule.Lock()
	delete(hMemRules.White, absPath)
	rwLockRule.Unlock()

	return nil
}

// 查询白名单总数
func RulesGetWhiteTotle() (totCnt int, err error) {
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RulesGetWhiteTotle: %s\n", err)
		return totCnt, err
	}

	sql := fmt.Sprintf("select count(*) from whitelist")
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("RulesGetWhiteTotle(): %s", err)
		return totCnt, errors.New("错误:获取白名单总数失败")
	}
	defer rows.Close()

	totCnt = 0
	for rows.Next() {
		rows.Scan(&totCnt)
		break
	}
	rows.Close()

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RulesGetWhiteTotle(commit transaction): %s\n", err)
		tx.Rollback()
		return totCnt, err
	}

	return totCnt, nil
}

// 查询白名单记录(start, length用来分页)
// limit(start, length)
// 从start开始，取length条记录
func RulesQueryWhite(start int, length int) (files []string, err error) {
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RulesQueryWhite: %s\n", err)
		return files, err
	}

	sql := fmt.Sprintf("select path from whitelist order by path asc limit %d, %d", start, length)
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("RulesQueryWhite(): %s", err)
		return files, errors.New("错误:获取白名单记录失败")
	}
	defer rows.Close()

	var file string
	for rows.Next() {
		rows.Scan(&file)
		files = append(files, file)
	}
	rows.Close()

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RulesQueryWhite(commit transaction): %s\n", err)
		tx.Rollback()
		return files, err
	}

	return files, nil
}

// 查找是否在黑名单
func RuleCheckInBlack(fpath string) (bFind bool, err error) {
	bFind = false
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RuleCheckInBlack: %s\n", err)
		return bFind, err
	}

	absPath, _ := filepath.Abs(fpath)
	sql := fmt.Sprintf("select * from blacklist where path = '%s'", absPath)
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("RuleCheckInBlack(): %s", err)
		return bFind, errors.New("错误:查找黑名单数据失败")
	}
	defer rows.Close()

	for rows.Next() {
		bFind = true
		break
	}
	rows.Close()

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RuleCheckInBlack(commit transaction): %s\n", err)
		tx.Rollback()
		return bFind, err
	}

	return bFind, nil
}

// 添加黑名单(可能是目录或末尾带*)
func RulesAddBlack(fpath string) (err error) {
	bFind, err := RuleCheckInWhite(fpath)
	if err != nil {
		return err
	}

	if bFind == true {
		return errors.New("错误:该记录已经存在于白名单中")
	}
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RulesAddBlack: %s\n", err)
		return err
	}

	absPath, _ := filepath.Abs(fpath)
	sql := fmt.Sprintf("insert into blacklist (id, path) values (null, '%s')", absPath)
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("RulesAddBlack(): %s", err)
		tx.Rollback()
		return errors.New("错误:添加白名单失败")
	}

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RulesAddBlack(commit transaction): %s\n", err)
		tx.Rollback()
		return err
	}

	// 更新内存
	rwLockRule.Lock()
	hMemRules.Black[absPath] = 0
	rwLockRule.Unlock()

	return nil
}

// 删除黑名单(可能是目录或末尾带*)
func RulesDelBlack(fpath string) (err error) {
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RulesDelBlack: %s\n", err)
		return err
	}

	absPath, _ := filepath.Abs(fpath)
	sql := fmt.Sprintf("delete from blacklist where path = '%s'", absPath)
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("RulesDelBlack(): %s", err)
		tx.Rollback()
		return errors.New("错误:删除白名单失败")
	}

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RulesDelBlack(commit transaction): %s\n", err)
		tx.Rollback()
		return err
	}

	// 更新内存
	rwLockRule.Lock()
	delete(hMemRules.Black, absPath)
	rwLockRule.Unlock()

	return nil
}

// 查询黑名单总数
func RulesGetBlackTotle() (totCnt int, err error) {
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RulesGetBlackTotle: %s\n", err)
		return totCnt, err
	}

	sql := fmt.Sprintf("select count(*) from blacklist")
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("RulesGetBlackTotle(): %s", err)
		return totCnt, errors.New("错误:获取黑名单总数失败")
	}
	defer rows.Close()

	totCnt = 0
	for rows.Next() {
		rows.Scan(&totCnt)
		break
	}
	rows.Close()

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RulesGetBlackTotle(commit transaction): %s\n", err)
		tx.Rollback()
		return totCnt, err
	}

	return totCnt, nil
}

// 查询黑名单记录(start, length用来分页)
// limit(start, length)
// 从start开始，取length条记录
func RulesQueryBlack(start int, length int) (files []string, err error) {
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RulesQueryBlack: %s\n", err)
		return files, err
	}

	sql := fmt.Sprintf("select path from blacklist order by path asc limit %d, %d", start, length)
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("RulesQueryBlack(): %s", err)
		return files, errors.New("错误:获取黑名单记录失败")
	}
	defer rows.Close()

	var file string
	for rows.Next() {
		rows.Scan(&file)
		files = append(files, file)
	}
	rows.Close()

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RulesQueryBlack(commit transaction): %s\n", err)
		tx.Rollback()
		return files, err
	}

	return files, nil
}

/////////////////////////////////
// Safe
/////////////////////////////////

// 获取系统基本防护配置
func RulesSafeBaseGet() (cfg SafeBaseConfig, err error) {
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RulesSafeBaseGet: %s\n", err)
		return cfg, err
	}

	sql := fmt.Sprintf("select Mode, WinDir, WinStart, WinFormat, WinProc, WinService from safe_base where id = 1")
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("RulesSafeBaseGet(): %s", err)
		return cfg, errors.New("错误:获取系统基本防护配置失败")
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&cfg.Mode, &cfg.WinDir, &cfg.WinStart, &cfg.WinFormat, &cfg.WinProc, &cfg.WinService)
		break
	}
	rows.Close()

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RulesSafeBaseGet(commit transaction): %s\n", err)
		tx.Rollback()
		return cfg, err
	}
	return cfg, nil
}

// 设置系统基本防护配置
func RulesSafeBaseSet(cfg SafeBaseConfig) (err error) {
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RulesSafeBaseSet: %s\n", err)
		return err
	}

	sql := fmt.Sprintf("update safe_base set Mode = %d, WinDir = %d, WinStart = %d, WinFormat = %d, WinProc = %d, WinService = %d where id = 1", cfg.Mode, cfg.WinDir, cfg.WinStart, cfg.WinFormat, cfg.WinProc, cfg.WinService)
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("RulesSafeBaseSet(): %s", err)
		tx.Rollback()
		return errors.New("错误:设置系统基本防护配置失败")
	}

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RulesSafeBaseSet(commit transaction): %s\n", err)
		tx.Rollback()
		return err
	}

	// 更新内存
	rwLockRule.Lock()
	hMemRules.SafeBaseCfg.Mode = cfg.Mode
	hMemRules.SafeBaseCfg.WinDir = cfg.WinDir
	hMemRules.SafeBaseCfg.WinStart = cfg.WinStart
	hMemRules.SafeBaseCfg.WinFormat = cfg.WinFormat
	hMemRules.SafeBaseCfg.WinProc = cfg.WinProc
	hMemRules.SafeBaseCfg.WinService = cfg.WinService
	rwLockRule.Unlock()
	return nil
}

// 导出系统基本防护配置 ini
func RulesSafeBaseSave() (saveString string, err error) {
	base, err := RulesSafeBaseGet()
	if err != nil {
		return saveString, err
	}

	saveString = ""
	saveString += "### 配置文件\n\n"

	saveString += "[INFO]\n"
	saveString += "Name = SafeBase\n"
	saveString += "\n"

	saveString += "[CONFIG]\n"
	saveString += "Mode = " + string(base.Mode) + "\n"
	saveString += "WinDir = " + string(base.WinDir) + "\n"
	saveString += "WinStart = " + string(base.WinStart) + "\n"
	saveString += "WinFormat = " + string(base.WinFormat) + "\n"
	saveString += "WinProc = " + string(base.WinProc) + "\n"
	saveString += "WinService = " + string(base.WinService) + "\n"
	saveString += "\n"

	return saveString, nil
}

// 获取系统增强防护配置
func RulesSafeHighGet() (cfg SafeHighConfig, err error) {
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RulesSafeHighGet: %s\n", err)
		return cfg, err
	}

	sql := fmt.Sprintf("select Mode, AddService, AutoRun, AddStart, ReadWrite, CreateExe, LoadSys, ProcInject from safe_high where id = 1")
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("RulesSafeHighGet(): %s", err)
		return cfg, errors.New("错误:获取系统增强防护配置失败")
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&cfg.Mode, &cfg.AddService, &cfg.AutoRun, &cfg.AddStart, &cfg.ReadWrite, &cfg.CreateExe, &cfg.LoadSys, &cfg.ProcInject)
		break
	}
	rows.Close()

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RulesSafeHighGet(commit transaction): %s\n", err)
		tx.Rollback()
		return cfg, err
	}
	return cfg, nil
}

// 设置系统增强防护配置
func RulesSafeHighSet(cfg SafeHighConfig) (err error) {
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RulesSafeHighGet: %s\n", err)
		return err
	}

	sql := fmt.Sprintf("update safe_high set Mode = %d, AddService = %d, AutoRun = %d, AddStart = %d, ReadWrite = %d, CreateExe = %d, LoadSys = %d, ProcInject = %d  where id = 1", cfg.Mode, cfg.AddService, cfg.AutoRun, cfg.AddStart, cfg.ReadWrite, cfg.CreateExe, cfg.LoadSys, cfg.ProcInject)
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("RulesSafeHighGet(): %s", err)
		tx.Rollback()
		return errors.New("错误:设置系统增强防护配置失败")
	}

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RulesSafeHighGet(commit transaction): %s\n", err)
		tx.Rollback()
		return err
	}

	// 更新内存
	rwLockRule.Lock()
	hMemRules.SafeHighCfg.Mode = cfg.Mode
	hMemRules.SafeHighCfg.AddService = cfg.AddService
	hMemRules.SafeHighCfg.AutoRun = cfg.AutoRun
	hMemRules.SafeHighCfg.AddStart = cfg.AddStart
	hMemRules.SafeHighCfg.ReadWrite = cfg.ReadWrite
	hMemRules.SafeHighCfg.CreateExe = cfg.CreateExe
	hMemRules.SafeHighCfg.LoadSys = cfg.LoadSys
	hMemRules.SafeHighCfg.ProcInject = cfg.ProcInject
	rwLockRule.Unlock()
	return nil
}

// 导出系统增强防护配置 ini
func RulesSafeHighSave() (saveString string, err error) {
	base, err := RulesSafeHighGet()
	if err != nil {
		return saveString, err
	}

	saveString = ""
	saveString += "### 配置文件\n\n"

	saveString += "[INFO]\n"
	saveString += "Name = SafeHigh\n"
	saveString += "\n"

	saveString += "[CONFIG]\n"
	saveString += "Mode = " + string(base.Mode) + "\n"
	saveString += "AddService = " + string(base.AddService) + "\n"
	saveString += "AutoRun = " + string(base.AutoRun) + "\n"
	saveString += "AddStart = " + string(base.AddStart) + "\n"
	saveString += "ReadWrite = " + string(base.ReadWrite) + "\n"
	saveString += "CreateExe = " + string(base.CreateExe) + "\n"
	saveString += "LoadSys = " + string(base.LoadSys) + "\n"
	saveString += "ProcInject = " + string(base.ProcInject) + "\n"
	saveString += "\n"

	return saveString, nil
}

// 添加系统目录及文件 WinDir
func RulesAddSafeBaseWinDir(fpath string, perm string) (err error) {
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RulesAddSafeBaseWinDir: %s\n", err)
		return err
	}

	absPath, _ := filepath.Abs(fpath)
	sql := fmt.Sprintf("insert into win_dir (id, path, perm) values (null, '%s', '%s')", absPath, perm)
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("RulesAddSafeBaseWinDir(): %s", err)
		tx.Rollback()
		return errors.New("错误:添加系统目录或文件失败")
	}

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RulesAddSafeBaseWinDir(commit transaction): %s\n", err)
		tx.Rollback()
		return err
	}

	// 更新内存
	rwLockRule.Lock()
	hMemRules.WinDir[absPath] = perm
	rwLockRule.Unlock()

	return nil
}

// 删除系统目录及文件 WinDir
func RulesDelSafeBaseWinDir(fpath string) (err error) {
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RulesDelSafeBaseWinDir: %s\n", err)
		return err
	}

	absPath, _ := filepath.Abs(fpath)
	sql := fmt.Sprintf("delete from win_dir where path = '%s'", absPath)
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("RulesDelSafeBaseWinDir(): %s", err)
		tx.Rollback()
		return errors.New("错误:删除系统目录及文件失败")
	}

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RulesDelSafeBaseWinDir(commit transaction): %s\n", err)
		tx.Rollback()
		return err
	}

	// 更新内存
	rwLockRule.Lock()
	delete(hMemRules.WinDir, absPath)
	rwLockRule.Unlock()

	return nil
}

// 查询系统目录及文件 WinDir
func RulesQuerySafeBaseWinDir() (files map[string]string, err error) {
	files = make(map[string]string)
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RulesQuerySafeBaseWinDir: %s\n", err)
		return files, err
	}

	sql := fmt.Sprintf("select path, perm from win_dir")
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("RulesQuerySafeBaseWinDir(): %s", err)
		return files, errors.New("错误:查询系统目录及文件失败")
	}
	defer rows.Close()

	var file, perm string
	for rows.Next() {
		rows.Scan(&file, &perm)
		files[file] = perm
	}
	rows.Close()

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RulesQuerySafeBaseWinDir(commit transaction): %s\n", err)
		tx.Rollback()
		return files, err
	}

	return files, nil
}

/////
// 添加系统启动文件 WinStart
func RulesAddSafeBaseWinStart(fpath string, perm string) (err error) {
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RulesAddSafeBaseWinStart: %s\n", err)
		return err
	}

	absPath, _ := filepath.Abs(fpath)
	sql := fmt.Sprintf("insert into win_start (id, path, perm) values (null, '%s', '%s')", absPath, perm)
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("RulesAddSafeBaseWinStart(): %s", err)
		tx.Rollback()
		return errors.New("错误:添加系统启动文件失败")
	}

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RulesAddSafeBaseWinStart(commit transaction): %s\n", err)
		tx.Rollback()
		return err
	}

	// 更新内存
	rwLockRule.Lock()
	hMemRules.WinStart[absPath] = perm
	rwLockRule.Unlock()
	return nil
}

// 删除系统启动文件 WinStart
func RulesDelSafeBaseWinStart(fpath string) (err error) {
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RulesDelSafeBaseWinStart: %s\n", err)
		return err
	}

	absPath, _ := filepath.Abs(fpath)
	sql := fmt.Sprintf("delete from win_start where path = '%s'", absPath)
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("RulesDelSafeBaseWinStart(): %s", err)
		tx.Rollback()
		return errors.New("错误:删除系统启动文件失败")
	}

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RulesDelSafeBaseWinStart(commit transaction): %s\n", err)
		tx.Rollback()
		return err
	}

	// 更新内存
	rwLockRule.Lock()
	delete(hMemRules.WinDir, absPath)
	rwLockRule.Unlock()
	return nil
}

// 查询系统启动文件 WinStart
func RulesQuerySafeBaseWinStart() (files map[string]string, err error) {
	files = make(map[string]string)
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RulesQuerySafeBaseWinStart: %s\n", err)
		return files, err
	}

	sql := fmt.Sprintf("select path, perm from win_start")
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("RulesQuerySafeBaseWinStart(): %s", err)
		return files, errors.New("错误:查询系统启动文件失败")
	}
	defer rows.Close()

	var file, perm string
	for rows.Next() {
		rows.Scan(&file, &perm)
		files[file] = perm
	}
	rows.Close()

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RulesQuerySafeBaseWinStart(commit transaction): %s\n", err)
		tx.Rollback()
		return files, err
	}

	return files, nil
}

////
// 添加系统关键进程 WinProc
func RulesAddSafeBaseWinProc(fpath string) (err error) {
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RulesAddSafeBaseWinProc: %s\n", err)
		return err
	}

	absPath, _ := filepath.Abs(fpath)
	sql := fmt.Sprintf("insert into win_proc (id, path) values (null, '%s')", absPath)
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("RulesAddSafeBaseWinProc(): %s", err)
		tx.Rollback()
		return errors.New("错误:添加系统关键进程件失败")
	}

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RulesAddSafeBaseWinProc(commit transaction): %s\n", err)
		tx.Rollback()
		return err
	}

	// 更新内存
	rwLockRule.Lock()
	hMemRules.WinProc[absPath] = 0
	rwLockRule.Unlock()

	return nil
}

// 删除系统关键进程 WinProc
func RulesDelSafeBaseWinProc(fpath string) (err error) {
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RulesDelSafeBaseWinProc: %s\n", err)
		return err
	}

	absPath, _ := filepath.Abs(fpath)
	sql := fmt.Sprintf("delete from win_proc where path = '%s'", absPath)
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("RulesDelSafeBaseWinProc(): %s", err)
		tx.Rollback()
		return errors.New("错误:删除系统关键进程失败")
	}

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RulesDelSafeBaseWinProc(commit transaction): %s\n", err)
		tx.Rollback()
		return err
	}

	// 更新内存
	rwLockRule.Lock()
	delete(hMemRules.WinProc, absPath)
	rwLockRule.Unlock()
	return nil
}

// 查询系统关键进程 WinProc
func RulesQuerySafeBaseWinProc() (files []string, err error) {
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RulesQuerySafeBaseWinProc: %s\n", err)
		return files, err
	}

	sql := fmt.Sprintf("select path from win_proc")
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("RulesQuerySafeBaseWinProc(): %s", err)
		return files, errors.New("错误:查询系统关键进程失败")
	}
	defer rows.Close()

	var file string
	for rows.Next() {
		rows.Scan(&file)
		files = append(files, file)
	}
	rows.Close()

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RulesQuerySafeBaseWinProc(commit transaction): %s\n", err)
		tx.Rollback()
		return files, err
	}

	return files, nil
}

/////////////////////////////////
// Account
/////////////////////////////////

// 获取账户安全配置
func RulesAccountGet() (cfg AccountConfig, err error) {
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RulesAccountGet: %s\n", err)
		return cfg, err
	}

	sql := fmt.Sprintf("select Mode, SafeLev, PwdComplex, PwdMinLen, PwdUsedMin, PwdUsedMax, PwdOldNum, AccountTimes, AccountMinute from safe_account where id = 1")
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("RulesAccountGet(): %s", err)
		return cfg, errors.New("错误:获取账户安全配置失败")
	}
	defer rows.Close()

	for rows.Next() {
		rows.Scan(&cfg.Mode, &cfg.SafeLev, &cfg.PwdComplex, &cfg.PwdMinLen, &cfg.PwdUsedMin, &cfg.PwdUsedMax, &cfg.PwdOldNum, &cfg.AccountTimes, &cfg.AccountMinute)
		break
	}
	rows.Close()

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RulesAccountGet(commit transaction): %s\n", err)
		tx.Rollback()
		return cfg, err
	}
	return cfg, nil
}

// 设置账户安全配置
func RulesAccountSet(cfg AccountConfig) (err error) {
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RulesAccountSet: %s\n", err)
		return err
	}

	sql := fmt.Sprintf("update safe_account set Mode = %d, SafeLev = %d, PwdComplex = %d, PwdMinLen = %d, PwdUsedMin = %d, PwdUsedMax = %d, PwdOldNum = %d, AccountTimes = %d, AccountMinute = %d where id = 1", cfg.Mode, cfg.SafeLev, cfg.PwdComplex, cfg.PwdMinLen, cfg.PwdUsedMin, cfg.PwdUsedMax, cfg.PwdOldNum, cfg.AccountTimes, cfg.AccountMinute)
	_, err = tx.Exec(sql)
	if err != nil {
		log.Printf("RulesAccountSet(): %s", err)
		tx.Rollback()
		return errors.New("错误:设置账户安全配置失败")
	}

	// 事务提交
	err = tx.Commit()
	if err != nil {
		log.Printf("RulesAccountSet(commit transaction): %s\n", err)
		tx.Rollback()
		return err
	}

	// 更新系统配置
	return nil
}

// 导出账户安全配置 ini
func RulesAccountSave() (saveString string, err error) {
	acc, err := RulesAccountGet()
	if err != nil {
		return saveString, err
	}

	saveString = ""
	saveString += "### 配置文件\n\n"

	saveString += "[INFO]\n"
	saveString += "Name = Account\n"
	saveString += "\n"

	saveString += "[CONFIG]\n"
	saveString += "Mode = " + string(acc.Mode) + "\n"
	saveString += "SafeLev = " + string(acc.SafeLev) + "\n"
	saveString += "PwdComplex = " + string(acc.PwdComplex) + "\n"
	saveString += "PwdMinLen = " + string(acc.PwdMinLen) + "\n"
	saveString += "PwdUsedMin = " + string(acc.PwdUsedMin) + "\n"
	saveString += "PwdUsedMax = " + string(acc.PwdUsedMax) + "\n"
	saveString += "PwdOldNum = " + string(acc.PwdOldNum) + "\n"
	saveString += "AccountTimes = " + string(acc.AccountTimes) + "\n"
	saveString += "AccountMinute = " + string(acc.AccountMinute) + "\n"
	saveString += "\n"

	return saveString, nil
}
