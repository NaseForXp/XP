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
	digest := h.Sum([]byte(s))
	md5string = fmt.Sprintf("%x", digest)
	return md5string
}

// 用户密码验证
func RulesCheckUserPassword(user string, pwd string) (isOK bool, uid int, user_type int, err error) {
	db := hDbRules
	tx, err := db.Begin()
	if err != nil {
		log.Printf("RulesCheckUserPassword: %s\n", err)
		return false, uid, user_type, err
	}

	sql := fmt.Sprintf("select uid, user_type, password from user where uname = '%s'", user)
	rows, err := db.Query(sql)
	if err != nil {
		log.Printf("RulesCheckUserPassword(): %s", err)
		return false, uid, user_type, errors.New("错误:查询用户名密码失败")
	}
	defer rows.Close()

	var password string
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
		return false, uid, user_type, err
	}

	// 校验密码
	md5Pwd := RulesGetMd5String(pwdUsedInfo + pwd)
	if md5Pwd == password {
		return true, uid, user_type, nil
	}
	return false, uid, user_type, nil
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
