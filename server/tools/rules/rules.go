package rules

import (
	"../RootDir"
	"crypto/md5"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"github.com/larspensjo/config"
	_ "github.com/mattn/go-sqlite3"
)

// 全局变量定义
var (
	rwLockRule sync.RWMutex // 全局读写锁 - 内存中的规则
	//hMemRules  ruleMemHandle                    // 规则内存句柄
	configFile string  = "config.ini"     // 配置文件相对路径
	ruleDB     string  = "rule\\rules.db" // 规则数据库文件
	hDbRules   *sql.DB                    // 规则数据库句柄
)

func RulesInit() (err error) {
	err = RulesConnectDb()
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

// 添加白名单(可能是目录或末尾带*)
func RulesAddWhite(fpath string) (err error) {
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

// 添加黑名单(可能是目录或末尾带*)
func RulesAddBlack(fpath string) (err error) {
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
