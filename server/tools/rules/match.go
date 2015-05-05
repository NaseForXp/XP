package rules

/*

#include <string.h>
#include <stdlib.h>
#include <tchar.h>
#include <locale.h>
#include <windows.h>
#include "go_sewindows.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/mahonia"
)

//export Go_file_create
func Go_file_create(user_name, process, file_path *C.char) C.BOOLEAN {
	decoder := mahonia.NewDecoder("GBK")
	uname := C.GoString(user_name)
	proc := C.GoString(process)
	fpath := C.GoString(file_path)

	uname = decoder.ConvertString(uname)
	proc = decoder.ConvertString(proc)
	fpath = decoder.ConvertString(fpath)

	fmt.Println("Go_file_create :", uname, proc, fpath)
	return C.FALSE
}

// Sewindows驱动初始化
func SewindowsInit() (err error) {
	var ret C.int
	ret = C.C_SewinInit()

	if ret != 0 {
		return errors.New("错误:驱动初始化失败")
	}
	return nil
}

// Sewindows驱动设置模式为监视模式
func SewindowsSetModeNotify() (err error) {
	bret := C.C_SewinSetModeNotify()

	if bret != C.TRUE {
		return errors.New("错误:驱动设置监视模式失败")
	}
	return nil
}

// Sewindows驱动设置模式为防护模式
func SewindowsSetModeIntercept() (err error) {
	bret := C.C_SewinSetModeIntercept()

	if bret != C.TRUE {
		return errors.New("错误:驱动设置防护模式失败")
	}
	return nil
}

// Sewindows驱动注册回调函数
func SewindowsRegOps() (err error) {
	bret := C.C_SewinRegOps()
	if bret != C.TRUE {
		return errors.New("错误:驱动注册处理函数失败")
	}
	return nil
}

func RuleMatchInit() (err error) {
	err = SewindowsInit()
	if err != nil {
		return err
	}

	err = SewindowsSetModeNotify()
	if err != nil {
		return err
	}

	err = SewindowsRegOps()
	if err != nil {
		return err
	}

	return nil
}

// 文件写操作
func RuleMatchFileWrite(uname, proc, file string) bool {
	// 系统文件及目录
	for p, _ := range hMemRules.WinDir {
		if strings.Index(file, p) == 0 {
			// 访问保护目录，要拒绝
			// Log
			return false
		}
	}

	// 系统启动文件
	_, ok := hMemRules.WinStart[file]
	if ok {
		// Log
		return false
	}
	return true
}

// 文件删除操作
func RuleMatchFileUnlink(uname, proc, file string) bool {
	// 系统文件及目录
	for p, _ := range hMemRules.WinDir {
		if strings.Index(file, p) == 0 {
			// 访问保护目录，要拒绝
			// Log
			return false
		}
	}

	// 系统启动文件
	_, ok := hMemRules.WinStart[file]
	if ok {
		// Log
		return false
	}
	return true
}

// 文件移动操作
func RuleMatchFileRename(uname, proc, file, new_file string) bool {
	// 系统文件及目录
	for p, _ := range hMemRules.WinDir {
		if strings.Index(file, p) == 0 {
			// 访问保护目录，要拒绝
			// Log
			return false
		}
		if strings.Index(new_file, p) == 0 {
			// 访问保护目录，要拒绝
			// Log
			return false
		}
	}

	// 系统启动文件
	_, ok := hMemRules.WinStart[file]
	if ok {
		// Log
		return false
	}

	_, ok = hMemRules.WinStart[new_file]
	if ok {
		// Log
		return false
	}

	// exe文件
	lens := len(new_file)
	if lens > 4 {
		if file[lens-4:] == ".exe" {
			// Log
			return false
		}
	}
	return true
}

// 文件创建操作
func RuleMatchFileCreate(uname, proc, file string) bool {
	// 系统文件及目录
	for p, _ := range hMemRules.WinDir {
		if strings.Index(file, p) == 0 {
			// 访问保护目录，要拒绝
			// Log
			return false
		}
	}

	// 系统启动文件
	_, ok := hMemRules.WinStart[file]
	if ok {
		// Log
		return false
	}

	// exe文件
	lens := len(file)
	if lens > 4 {
		if file[lens-4:] == ".exe" {
			// Log
			return false
		}
	}
	return true
}

// 规则匹配 - 文件操作
func RuleMatchFile(uname, proc, file, new_file, op string) bool {
	uname, _ = filepath.Abs(strings.ToLower(uname))
	proc, _ = filepath.Abs(strings.ToLower(proc))
	file, _ = filepath.Abs(strings.ToLower(file))
	new_file, _ = filepath.Abs(strings.ToLower(new_file))

	rwLockRule.RLock()
	defer rwLockRule.RUnlock()

	// 白名单放行
	_, ok := hMemRules.White[proc]
	if ok {
		return true
	}

	// 黑名单拒绝
	_, ok = hMemRules.Black[proc]
	if ok {
		// Log
		//xplog.LogInsertEvent("黑名单", "防护模式", uname, proc, file, op, "拒绝")
		return false
	}

	switch op {
	case "FILE_WRITE":
		return RuleMatchFileWrite(uname, proc, file)
	case "FILE_UNLINK":
		return RuleMatchFileUnlink(uname, proc, file)
	case "FILE_RENAME":
		return RuleMatchFileRename(uname, proc, file, new_file)
	case "FILE_FORMAT":
		break
	case "FILE_IO":
		break
	case "FILE_CREATE":
		return RuleMatchFileCreate(uname, proc, file)
	}

	return true
}

// 规则匹配 - 进程操作
func RuleMatchProcess(uname, process, dst_proc, op string) bool {
	uname, _ = filepath.Abs(strings.ToLower(uname))
	process, _ = filepath.Abs(strings.ToLower(process))
	dst_proc, _ = filepath.Abs(strings.ToLower(dst_proc))

	rwLockRule.RLock()
	defer rwLockRule.RUnlock()

	// 白名单放行
	_, ok := hMemRules.White[process]
	if ok {
		return true
	}

	// 黑名单拒绝
	_, ok = hMemRules.Black[process]
	if ok {
		// Log
		return false
	}
	switch op {
	case "PROC_KILL":
		// 系统关键进程
		_, ok := hMemRules.WinProc[dst_proc]
		if ok {
			// Log
			return false
		}
	case "PROC_INJECT":
		// log
		return false
	}

	return true
}

// 规则匹配 - 服务
func RuleMatchService(uname, process, service_name, binPath, op string) bool {
	uname, _ = filepath.Abs(strings.ToLower(uname))
	process, _ = filepath.Abs(strings.ToLower(process))
	binPath, _ = filepath.Abs(strings.ToLower(binPath))

	rwLockRule.RLock()
	defer rwLockRule.RUnlock()

	// 白名单放行
	_, ok := hMemRules.White[process]
	if ok {
		return true
	}

	// 黑名单拒绝
	_, ok = hMemRules.Black[process]
	if ok {
		// Log
		return false
	}

	switch op {
	case "SRV_CREATE":
		// Log
		return false
	case "SRV_DEL":
		// Log change
		return false
	case "SRV_CHANGE":
		// Log change
		return false
	}

	return true
}

// 规则匹配 - 驱动
func RuleMatchDrive(uname, process, service_name, binPath, op string) bool {
	uname, _ = filepath.Abs(strings.ToLower(uname))
	process, _ = filepath.Abs(strings.ToLower(process))
	binPath, _ = filepath.Abs(strings.ToLower(binPath))

	rwLockRule.RLock()
	defer rwLockRule.RUnlock()

	// 白名单放行
	_, ok := hMemRules.White[process]
	if ok {
		return true
	}

	// 黑名单拒绝
	_, ok = hMemRules.Black[process]
	if ok {
		// Log
		return false
	}

	switch op {
	case "DRIVE_LOAD":
		// Log
		return false
	}

	return true
}
