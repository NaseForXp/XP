package rules

import (
	"fmt"
	"path/filepath"
	"strings"
)

type StringInt map[string]int
type StringStr map[string]string

type RuleMemHandle struct {
	White        StringInt      // 白名单的程序和目录
	Black        StringInt      // 黑名单的程序和目录
	WinDir       StringStr      // 受保护的系统目录
	WinStart     StringStr      // 受保护的系统启动项
	WinProc      StringInt      // 受保护的系统进程
	HighWinStart StringStr      // 增强防护 - 开机启动项
	SafeBaseCfg  SafeBaseConfig // 系统防护_基本防护配置
	SafeHighCfg  SafeHighConfig // 系统防护_增强防护配置
	AccountCfg   AccountConfig  // 账户安全配置
}

/* 全局变量定义 - rules.go中定义
var (
	rwLockRule sync.RWMutex  // 全局读写锁 - 内存中的规则
	hMemRules  RuleMemHandle // 规则内存句柄
)
*/

// 内存初始化
func RulesMemInit() (err error) {
	hMemRules.White = make(StringInt)
	hMemRules.Black = make(StringInt)
	hMemRules.WinDir = make(StringStr)
	hMemRules.WinStart = make(StringStr)
	hMemRules.WinProc = make(StringInt)
	hMemRules.HighWinStart = make(StringStr)

	// 获取白名单
	totCnt, err := RulesGetWhiteTotle()
	if err != nil {
		return err
	}

	if totCnt > 0 {
		ws, err := RulesQueryWhite(0, totCnt)
		if err != nil {
			return err
		}

		rwLockRule.Lock()
		for _, f := range ws {
			f, _ = filepath.Abs(strings.ToLower(f))
			hMemRules.White[f] = 0
		}
		rwLockRule.Unlock()
	}

	// 获取黑名单
	totCnt, err = RulesGetBlackTotle()
	if err != nil {
		return err
	}

	if totCnt > 0 {
		bs, err := RulesQueryBlack(0, totCnt)
		if err != nil {
			return err
		}

		rwLockRule.Lock()
		for _, f := range bs {
			f, _ = filepath.Abs(strings.ToLower(f))
			hMemRules.Black[f] = 0
		}
		rwLockRule.Unlock()
	}

	// 获取受保护系统目录及文件
	windir, err := RulesQuerySafeBaseWinDir()
	if err != nil {
		return err
	}

	rwLockRule.Lock()
	for f, perm := range windir {
		f, _ = filepath.Abs(strings.ToLower(f))
		hMemRules.WinDir[f] = perm

	}
	rwLockRule.Unlock()

	// 获取受保护系统启动文件
	winstart, err := RulesQuerySafeBaseWinStart()
	if err != nil {
		return err
	}

	rwLockRule.Lock()
	for f, perm := range winstart {
		f, _ = filepath.Abs(strings.ToLower(f))
		hMemRules.WinStart[f] = perm
	}
	rwLockRule.Unlock()

	// 获取受保护系统关键进程
	winsproc, err := RulesQuerySafeBaseWinProc()
	if err != nil {
		return err
	}

	rwLockRule.Lock()
	for _, f := range winsproc {
		f, _ = filepath.Abs(strings.ToLower(f))
		hMemRules.WinProc[f] = 0
	}
	rwLockRule.Unlock()

	// 获取开机启动注册表项
	startRegs, err := RulesQuerySafeHighWinStart()
	if err != nil {
		return err
	}

	rwLockRule.Lock()
	for f, _ := range startRegs {
		//f = strings.Replace(f, "\\\\", "\\", -1)
		hMemRules.HighWinStart[strings.ToUpper(f)] = ""
	}
	rwLockRule.Unlock()

	// 获取系统防护 - 基本防护 设置
	hMemRules.SafeBaseCfg, err = RulesSafeBaseGet()
	if err != nil {
		return err
	}

	// 获取系统防护 - 基本防护 设置
	hMemRules.SafeHighCfg, err = RulesSafeHighGet()
	if err != nil {
		return err
	}

	// 获取账户防护 设置
	hMemRules.AccountCfg, err = RulesAccountGet()
	if err != nil {
		return err
	}

	return err
}

func RulesMemGetHomeMode() (baseMode, highMode int) {
	rwLockRule.RLock()
	baseMode = hMemRules.SafeBaseCfg.Mode
	highMode = hMemRules.SafeHighCfg.Mode
	rwLockRule.RUnlock()
	return baseMode, highMode
}

func RulesMemPrint() {
	fmt.Println("SafeBaseCfg:", hMemRules.SafeBaseCfg)
	fmt.Println("SafeHighCfg:", hMemRules.SafeHighCfg)
	fmt.Println("AccountCfg:", hMemRules.AccountCfg)
	fmt.Println("Black:", hMemRules.Black)
	fmt.Println("White:", hMemRules.White)
	fmt.Println("WinDir:", hMemRules.WinDir)
	fmt.Println("WinStart:", hMemRules.WinStart)
	fmt.Println("WinProc:", hMemRules.WinProc)
	fmt.Println("HighWinStart:", hMemRules.HighWinStart)
}
