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
	"path/filepath"
	"strings"

	"../xplog"
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

	bret := RuleMatchFileCreate(uname, proc, fpath)
	if bret {
		return C.TRUE
	}
	return C.FALSE
}

//export Go_file_unlink
func Go_file_unlink(user_name, process, file_path *C.char) C.BOOLEAN {
	decoder := mahonia.NewDecoder("GBK")
	uname := C.GoString(user_name)
	proc := C.GoString(process)
	fpath := C.GoString(file_path)

	uname = decoder.ConvertString(uname)
	proc = decoder.ConvertString(proc)
	fpath = decoder.ConvertString(fpath)

	bret := RuleMatchFileUnlink(uname, proc, fpath)
	if bret {
		return C.TRUE
	}
	return C.FALSE
}

//export Go_file_read
func Go_file_read(user_name, process, file_path *C.char) C.BOOLEAN {
	decoder := mahonia.NewDecoder("GBK")
	uname := C.GoString(user_name)
	proc := C.GoString(process)
	fpath := C.GoString(file_path)

	uname = decoder.ConvertString(uname)
	proc = decoder.ConvertString(proc)
	fpath = decoder.ConvertString(fpath)

	bret := RuleMatchFileRead(uname, proc, fpath)
	if bret {
		return C.TRUE
	}
	return C.FALSE
}

//export Go_file_write
func Go_file_write(user_name, process, file_path *C.char) C.BOOLEAN {
	decoder := mahonia.NewDecoder("GBK")
	uname := C.GoString(user_name)
	proc := C.GoString(process)
	fpath := C.GoString(file_path)

	uname = decoder.ConvertString(uname)
	proc = decoder.ConvertString(proc)
	fpath = decoder.ConvertString(fpath)

	bret := RuleMatchFileWrite(uname, proc, fpath)
	if bret {
		return C.TRUE
	}
	return C.FALSE
}

//export Go_file_rename
func Go_file_rename(user_name, process, file_path *C.char, new_path *C.char) C.BOOLEAN {
	decoder := mahonia.NewDecoder("GBK")
	uname := C.GoString(user_name)
	proc := C.GoString(process)
	fpath := C.GoString(file_path)
	dpath := C.GoString(new_path)

	uname = decoder.ConvertString(uname)
	proc = decoder.ConvertString(proc)
	fpath = decoder.ConvertString(fpath)
	dpath = decoder.ConvertString(dpath)

	bret := RuleMatchFileRename(uname, proc, fpath, dpath)
	if bret {
		return C.TRUE
	}
	return C.FALSE
}

//export Go_disk_read
func Go_disk_read(user_name, process, file_path *C.char) C.BOOLEAN {
	decoder := mahonia.NewDecoder("GBK")
	uname := C.GoString(user_name)
	proc := C.GoString(process)
	fpath := C.GoString(file_path)

	uname = decoder.ConvertString(uname)
	proc = decoder.ConvertString(proc)
	fpath = decoder.ConvertString(fpath)

	bret := RuleMatchDiskReadWrite(uname, proc, fpath, "读磁盘")
	if bret {
		return C.TRUE
	}
	return C.FALSE
}

//export Go_disk_write
func Go_disk_write(user_name, process, file_path *C.char) C.BOOLEAN {
	decoder := mahonia.NewDecoder("GBK")
	uname := C.GoString(user_name)
	proc := C.GoString(process)
	fpath := C.GoString(file_path)

	uname = decoder.ConvertString(uname)
	proc = decoder.ConvertString(proc)
	fpath = decoder.ConvertString(fpath)

	bret := RuleMatchDiskReadWrite(uname, proc, fpath, "写磁盘")
	if bret {
		return C.TRUE
	}
	return C.FALSE
}

//export Go_disk_formate
func Go_disk_formate(user_name, process, file_path *C.char) C.BOOLEAN {
	decoder := mahonia.NewDecoder("GBK")
	uname := C.GoString(user_name)
	proc := C.GoString(process)
	fpath := C.GoString(file_path)

	uname = decoder.ConvertString(uname)
	proc = decoder.ConvertString(proc)
	fpath = decoder.ConvertString(fpath)

	bret := RuleMatchDiskFormat(uname, proc, fpath)
	if bret {
		return C.TRUE
	}
	return C.FALSE
}

//export Go_process_kill
func Go_process_kill(user_name, process, file_path *C.char) C.BOOLEAN {
	decoder := mahonia.NewDecoder("GBK")
	uname := C.GoString(user_name)
	proc := C.GoString(process)
	fpath := C.GoString(file_path)

	uname = decoder.ConvertString(uname)
	proc = decoder.ConvertString(proc)
	fpath = decoder.ConvertString(fpath)

	bret := RuleMatchProcessKill(uname, proc, fpath)
	if bret {
		return C.TRUE
	}
	return C.FALSE
}

//export Go_process_create_thread
func Go_process_create_thread(user_name, process, file_path *C.char) C.BOOLEAN {
	decoder := mahonia.NewDecoder("GBK")
	uname := C.GoString(user_name)
	proc := C.GoString(process)
	fpath := C.GoString(file_path)

	uname = decoder.ConvertString(uname)
	proc = decoder.ConvertString(proc)
	fpath = decoder.ConvertString(fpath)

	bret := RuleMatchProcessInject(uname, proc, fpath)
	if bret {
		return C.TRUE
	}
	return C.FALSE
}

//export Go_service_create
func Go_service_create(user_name, process, sname, binPath *C.char) C.BOOLEAN {
	decoder := mahonia.NewDecoder("GBK")
	uname := C.GoString(user_name)
	proc := C.GoString(process)
	name := C.GoString(sname)
	fpath := C.GoString(binPath)

	uname = decoder.ConvertString(uname)
	proc = decoder.ConvertString(proc)
	name = decoder.ConvertString(name)
	fpath = decoder.ConvertString(fpath)

	bret := RuleMatchServiceAdd(uname, proc, name, fpath)
	if bret {
		return C.TRUE
	}
	return C.FALSE
}

//export Go_service_delete
func Go_service_delete(user_name, process, sname *C.char) C.BOOLEAN {
	decoder := mahonia.NewDecoder("GBK")
	uname := C.GoString(user_name)
	proc := C.GoString(process)
	name := C.GoString(sname)

	uname = decoder.ConvertString(uname)
	proc = decoder.ConvertString(proc)
	name = decoder.ConvertString(name)

	bret := RuleMatchServiceChange(uname, proc, name, "删除服务")
	if bret {
		return C.TRUE
	}
	return C.FALSE
}

//export Go_service_change
func Go_service_change(user_name, process, sname *C.char) C.BOOLEAN {
	decoder := mahonia.NewDecoder("GBK")
	uname := C.GoString(user_name)
	proc := C.GoString(process)
	name := C.GoString(sname)

	uname = decoder.ConvertString(uname)
	proc = decoder.ConvertString(proc)
	name = decoder.ConvertString(name)

	bret := RuleMatchServiceChange(uname, proc, name, "修改服务")
	if bret {
		return C.TRUE
	}
	return C.FALSE
}

//export Go_driver_load
func Go_driver_load(user_name, process, sname, binPath *C.char) C.BOOLEAN {
	decoder := mahonia.NewDecoder("GBK")
	uname := C.GoString(user_name)
	proc := C.GoString(process)
	name := C.GoString(sname)
	fpath := C.GoString(binPath)

	uname = decoder.ConvertString(uname)
	proc = decoder.ConvertString(proc)
	name = decoder.ConvertString(name)
	fpath = decoder.ConvertString(fpath)

	bret := RuleMatchDriveLoad(uname, proc, name, fpath)
	if bret {
		return C.TRUE
	}
	return C.FALSE
}

//export Go_reg_set_value
func Go_reg_set_value(user_name, process, rpath, rvalue *C.char) C.BOOLEAN {
	decoder := mahonia.NewDecoder("GBK")
	uname := C.GoString(user_name)
	proc := C.GoString(process)
	regpath := C.GoString(rpath)
	regvalue := C.GoString(rvalue)

	uname = decoder.ConvertString(uname)
	proc = decoder.ConvertString(proc)
	regpath = decoder.ConvertString(regpath)
	regvalue = decoder.ConvertString(regvalue)

	bret := RuleMatchRegSetValue(uname, proc, regpath, regvalue)
	if bret {
		return C.TRUE
	}
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

	err = SewindowsSetModeIntercept()
	if err != nil {
		return err
	}

	err = SewindowsRegOps()
	if err != nil {
		return err
	}

	return nil
}

// 规则匹配 - 磁盘直接读写
func RuleMatchDiskReadWrite(uname, proc, file, opStr string) bool {
	proc, _ = filepath.Abs(proc)
	proc = strings.ToLower(proc)
	file, _ = filepath.Abs(file)
	file = strings.ToLower(file)

	rwLockRule.Lock()
	defer rwLockRule.Unlock()

	// 白名单放行
	_, ok := hMemRules.White[proc]
	if ok {
		//xplog.LogInsertEvent("白名单", "防护模式", uname, proc, file, opStr, "允许")
		return true
	}

	// 黑名单拒绝
	_, ok = hMemRules.Black[proc]
	if ok {
		xplog.LogInsertEvent("黑名单", "防护模式", uname, proc, file, opStr, "拒绝")
		return false
	}

	// 功能开启
	if hMemRules.SafeHighCfg.ReadWrite == 1 {
		if hMemRules.SafeHighCfg.Mode == 0 {
			xplog.LogInsertEvent("增强防护-防止磁盘被直接读写", "监视模式", uname, proc, file, opStr, "拒绝")
			return true
		} else {
			xplog.LogInsertEvent("增强防护-防止磁盘被直接读写", "防护模式", uname, proc, file, opStr, "拒绝")
			return false
		}
	}

	return true
}

// 规则匹配 - 磁盘格式化
func RuleMatchDiskFormat(uname, proc, file string) bool {
	proc, _ = filepath.Abs(proc)
	proc = strings.ToLower(proc)
	// 格式化路径已经是磁盘了
	file = strings.ToLower(file)

	opStr := "格式化磁盘"

	rwLockRule.Lock()
	defer rwLockRule.Unlock()

	// 白名单放行
	_, ok := hMemRules.White[proc]
	if ok {
		//xplog.LogInsertEvent("白名单", "防护模式", uname, proc, file, opStr, "允许")
		return true
	}

	// 黑名单拒绝
	_, ok = hMemRules.Black[proc]
	if ok {
		xplog.LogInsertEvent("黑名单", "防护模式", uname, proc, file, opStr, "拒绝")
		return false
	}

	// 功能开启
	if hMemRules.SafeBaseCfg.WinFormat == 1 {
		if hMemRules.SafeBaseCfg.Mode == 0 {
			xplog.LogInsertEvent("基本防护-防止格式化系统磁盘", "监视模式", uname, proc, file, opStr, "拒绝")
			return true
		} else {
			xplog.LogInsertEvent("基本防护-防止格式化系统磁盘", "防护模式", uname, proc, file, opStr, "拒绝")
			return false
		}
	}

	return true
}

// 文件读操作 - autorun.inf
func RuleMatchFileRead(uname, proc, file string) bool {
	proc, _ = filepath.Abs(proc)
	proc = strings.ToLower(proc)
	file, _ = filepath.Abs(file)
	file = strings.ToLower(file)

	rwLockRule.Lock()
	defer rwLockRule.Unlock()

	// 白名单放行
	_, ok := hMemRules.White[proc]
	if ok {
		//xplog.LogInsertEvent("白名单", "防护模式", uname, proc, file, "读文件", "允许")
		return true
	}

	// 黑名单拒绝
	_, ok = hMemRules.Black[proc]
	if ok {
		xplog.LogInsertEvent("黑名单", "防护模式", uname, proc, file, "读文件", "拒绝")
		return false
	}

	// 自动运行
	if hMemRules.SafeHighCfg.AutoRun == 1 {
		// 功能开启
		if strings.Index(file, "autorun.inf") > 0 {
			if hMemRules.SafeHighCfg.Mode == 0 {
				xplog.LogInsertEvent("增强防护-防止自动运行", "监视模式", uname, proc, file, "自动运行", "拒绝")
				return true
			} else {
				xplog.LogInsertEvent("增强防护-防止自动运行", "防护模式", uname, proc, file, "自动运行", "拒绝")
				return false
			}
		}
	}

	return true
}

// 文件写操作
func RuleMatchFileWrite(uname, proc, file string) bool {
	proc, _ = filepath.Abs(proc)
	proc = strings.ToLower(proc)
	file, _ = filepath.Abs(file)
	file = strings.ToLower(file)

	rwLockRule.Lock()
	defer rwLockRule.Unlock()

	// 白名单放行
	_, ok := hMemRules.White[proc]
	if ok {
		//xplog.LogInsertEvent("白名单", "防护模式", uname, proc, file, "写文件", "允许")
		return true
	}

	// 黑名单拒绝
	_, ok = hMemRules.Black[proc]
	if ok {
		xplog.LogInsertEvent("黑名单", "防护模式", uname, proc, file, "写文件", "拒绝")
		return false
	}

	// 自保护
	if strings.Index(file, naseClientDir) >= 0 {
		return false
	}

	// 系统文件及目录
	if hMemRules.SafeBaseCfg.WinDir == 1 {
		// 功能开启
		for p, _ := range hMemRules.WinDir {
			if strings.Index(file, p) == 0 {
				if hMemRules.SafeBaseCfg.Mode == 0 {
					xplog.LogInsertEvent("基本防护-系统文件及目录保护", "监视模式", uname, proc, file, "写文件", "拒绝")
					return true
				} else {
					xplog.LogInsertEvent("基本防护-系统文件及目录保护", "防护模式", uname, proc, file, "写文件", "拒绝")
					return false
				}
			}
		}
	}

	// 系统启动文件
	if hMemRules.SafeBaseCfg.WinStart == 1 {
		// 功能开启
		_, ok = hMemRules.WinStart[file]
		if ok {
			if hMemRules.SafeBaseCfg.Mode == 0 {
				xplog.LogInsertEvent("基本防护-系统启动文件保护", "监视模式", uname, proc, file, "写文件", "拒绝")
				return true
			} else {
				xplog.LogInsertEvent("基本防护-系统启动文件保护", "防护模式", uname, proc, file, "写文件", "拒绝")
				return false
			}
		}
	}

	return true
}

// 文件删除操作
func RuleMatchFileUnlink(uname, proc, file string) bool {
	proc, _ = filepath.Abs(proc)
	proc = strings.ToLower(proc)
	file, _ = filepath.Abs(file)
	file = strings.ToLower(file)

	rwLockRule.Lock()
	defer rwLockRule.Unlock()

	// 白名单放行
	_, ok := hMemRules.White[proc]
	if ok {
		//xplog.LogInsertEvent("白名单", "防护模式", uname, proc, file, "删除文件", "允许")
		return true
	}

	// 黑名单拒绝
	_, ok = hMemRules.Black[proc]
	if ok {
		xplog.LogInsertEvent("黑名单", "防护模式", uname, proc, file, "删除文件", "拒绝")
		return false
	}

	// 自保护
	if strings.Index(file, naseClientDir) >= 0 {
		return false
	}

	// 系统文件及目录
	if hMemRules.SafeBaseCfg.WinDir == 1 {
		// 功能开启
		for p, _ := range hMemRules.WinDir {
			if strings.Index(file, p) == 0 {
				if hMemRules.SafeBaseCfg.Mode == 0 {
					xplog.LogInsertEvent("基本防护-系统文件及目录保护", "监视模式", uname, proc, file, "删除文件", "拒绝")
					return true
				} else {
					xplog.LogInsertEvent("基本防护-系统文件及目录保护", "防护模式", uname, proc, file, "删除文件", "拒绝")
					return false
				}
			}
		}
	}

	// 系统启动文件
	if hMemRules.SafeBaseCfg.WinStart == 1 {
		// 功能开启
		_, ok = hMemRules.WinStart[file]
		if ok {
			if hMemRules.SafeBaseCfg.Mode == 0 {
				xplog.LogInsertEvent("基本防护-系统启动文件保护", "监视模式", uname, proc, file, "删除文件", "拒绝")
				return true
			} else {
				xplog.LogInsertEvent("基本防护-系统启动文件保护", "防护模式", uname, proc, file, "删除文件", "拒绝")
				return false
			}
		}
	}

	return true
}

// 文件移动操作
func RuleMatchFileRename(uname, proc, file, new_file string) bool {
	proc, _ = filepath.Abs(proc)
	proc = strings.ToLower(proc)
	file, _ = filepath.Abs(file)
	file = strings.ToLower(file)
	new_file, _ = filepath.Abs(new_file)
	new_file = strings.ToLower(new_file)

	rwLockRule.Lock()
	defer rwLockRule.Unlock()

	// 白名单放行
	_, ok := hMemRules.White[proc]
	if ok {
		//xplog.LogInsertEvent("白名单", "防护模式", uname, proc, file+"->"+new_file, "移动文件", "允许")
		return true
	}

	// 黑名单拒绝
	_, ok = hMemRules.Black[proc]
	if ok {
		xplog.LogInsertEvent("黑名单", "防护模式", uname, proc, file+"->"+new_file, "移动文件", "拒绝")
		return false
	}

	// 自保护
	if strings.Index(file, naseClientDir) >= 0 {
		return false
	}

	// 系统文件及目录
	if hMemRules.SafeBaseCfg.WinDir == 1 {
		// 功能开启
		for p, _ := range hMemRules.WinDir {
			if strings.Index(file, p) == 0 {
				if hMemRules.SafeBaseCfg.Mode == 0 {
					xplog.LogInsertEvent("基本防护-系统文件及目录保护", "监视模式", uname, proc, file+"->"+new_file, "移动文件", "拒绝")
					return true
				} else {
					xplog.LogInsertEvent("基本防护-系统文件及目录保护", "防护模式", uname, proc, file+"->"+new_file, "移动文件", "拒绝")
					return false
				}
			}
		}

		// 判断目标文件是否允许
		for p, _ := range hMemRules.WinDir {
			if strings.Index(new_file, p) == 0 {
				if hMemRules.SafeBaseCfg.Mode == 0 {
					xplog.LogInsertEvent("基本防护-系统文件及目录保护", "监视模式", uname, proc, file+"->"+new_file, "移动文件", "拒绝")
					return true
				} else {
					xplog.LogInsertEvent("基本防护-系统文件及目录保护", "防护模式", uname, proc, file+"->"+new_file, "移动文件", "拒绝")
					return false
				}
			}
		}
	}

	// 系统启动文件
	if hMemRules.SafeBaseCfg.WinStart == 1 {
		// 功能开启
		_, ok = hMemRules.WinStart[file]
		if ok {
			if hMemRules.SafeBaseCfg.Mode == 0 {
				xplog.LogInsertEvent("基本防护-系统启动文件保护", "监视模式", uname, proc, file+"->"+new_file, "移动文件", "拒绝")
				return true
			} else {
				xplog.LogInsertEvent("基本防护-系统启动文件保护", "防护模式", uname, proc, file+"->"+new_file, "移动文件", "拒绝")
				return false
			}
		}

		// 判断目标文件是否允许
		_, ok = hMemRules.WinStart[new_file]
		if ok {
			if hMemRules.SafeBaseCfg.Mode == 0 {
				xplog.LogInsertEvent("基本防护-系统启动文件保护", "监视模式", uname, proc, file+"->"+new_file, "移动文件", "拒绝")
				return true
			} else {
				xplog.LogInsertEvent("基本防护-系统启动文件保护", "防护模式", uname, proc, file+"->"+new_file, "移动文件", "拒绝")
				return false
			}
		}
	}

	// exe文件
	lens := len(new_file)
	if lens > 4 {
		if new_file[lens-4:] == ".exe" {
			if hMemRules.SafeHighCfg.Mode == 0 {
				xplog.LogInsertEvent("增强防护-禁止创建.exe文件", "监视模式", uname, proc, file+"->"+new_file, "移动文件", "拒绝")
				return true
			} else {
				xplog.LogInsertEvent("增强防护-禁止创建.exe文件", "防护模式", uname, proc, file+"->"+new_file, "移动文件", "拒绝")
				return false
			}
		}
	}

	return true
}

// 文件创建操作
func RuleMatchFileCreate(uname, proc, file string) bool {
	proc, _ = filepath.Abs(proc)
	proc = strings.ToLower(proc)
	file, _ = filepath.Abs(file)
	file = strings.ToLower(file)

	rwLockRule.Lock()
	defer rwLockRule.Unlock()

	// 白名单放行
	_, ok := hMemRules.White[proc]
	if ok {
		//xplog.LogInsertEvent("白名单", "防护模式", uname, proc, file, "创建文件", "允许")
		return true
	}

	// 黑名单拒绝
	_, ok = hMemRules.Black[proc]
	if ok {
		xplog.LogInsertEvent("黑名单", "防护模式", uname, proc, file, "创建文件", "拒绝")
		return false
	}

	// 自保护
	if strings.Index(file, naseClientDir) >= 0 {
		return false
	}

	// 系统文件及目录
	if hMemRules.SafeBaseCfg.WinDir == 1 {
		// 功能开启
		for p, _ := range hMemRules.WinDir {
			if strings.Index(file, p) == 0 {
				if hMemRules.SafeBaseCfg.Mode == 0 {
					xplog.LogInsertEvent("基本防护-系统文件及目录保护", "监视模式", uname, proc, file, "创建文件", "拒绝")
					return true
				} else {
					xplog.LogInsertEvent("基本防护-系统文件及目录保护", "防护模式", uname, proc, file, "创建文件", "拒绝")
					return false
				}
			}
		}
	}

	// 系统启动文件
	if hMemRules.SafeBaseCfg.WinStart == 1 {
		// 功能开启
		_, ok = hMemRules.WinStart[file]
		if ok {
			if hMemRules.SafeBaseCfg.Mode == 0 {
				xplog.LogInsertEvent("基本防护-系统启动文件保护", "监视模式", uname, proc, file, "创建文件", "拒绝")
				return true
			} else {
				xplog.LogInsertEvent("基本防护-系统启动文件保护", "防护模式", uname, proc, file, "创建文件", "拒绝")
				return false
			}
		}
	}

	// exe文件
	lens := len(file)
	if lens > 4 {
		if file[lens-4:] == ".exe" {
			if hMemRules.SafeHighCfg.Mode == 0 {
				xplog.LogInsertEvent("增强防护-禁止创建.exe文件", "监视模式", uname, proc, file, "创建文件", "拒绝")
				return true
			} else {
				xplog.LogInsertEvent("增强防护-禁止创建.exe文件", "防护模式", uname, proc, file, "创建文件", "拒绝")
				return false
			}
		}
	}
	return true
}

// 规则匹配 - 进程杀死
func RuleMatchProcessKill(uname, process, dst_proc string) bool {
	process, _ = filepath.Abs(process)
	process = strings.ToLower(process)
	dst_proc, _ = filepath.Abs(dst_proc)
	dst_proc = strings.ToLower(dst_proc)

	rwLockRule.Lock()
	defer rwLockRule.Unlock()

	// 白名单放行
	_, ok := hMemRules.White[process]
	if ok {
		//xplog.LogInsertEvent("白名单", "防护模式", uname, process, dst_proc, "进程杀死", "允许")
		return true
	}

	// 黑名单拒绝
	_, ok = hMemRules.Black[process]
	if ok {
		xplog.LogInsertEvent("黑名单", "防护模式", uname, process, dst_proc, "进程杀死", "拒绝")
		return false
	}

	// 自保护
	if strings.Index(dst_proc, naseClientDir) >= 0 {
		return false
	}

	// 系统关键进程
	if hMemRules.SafeBaseCfg.WinProc == 1 {
		// 功能开启
		_, ok = hMemRules.WinProc[dst_proc]
		if ok {
			if hMemRules.SafeBaseCfg.Mode == 0 {
				xplog.LogInsertEvent("基本防护-防止系统关键进程被杀死", "监视模式", uname, process, dst_proc, "进程杀死", "拒绝")
				return true
			} else {
				xplog.LogInsertEvent("基本防护-防止系统关键进程被杀死", "防护模式", uname, process, dst_proc, "进程杀死", "拒绝")
				return false
			}
		}
	}

	return true
}

// 规则匹配 - 进程被注入
func RuleMatchProcessInject(uname, process, dst_proc string) bool {
	process, _ = filepath.Abs(process)
	process = strings.ToLower(process)
	dst_proc, _ = filepath.Abs(dst_proc)
	dst_proc = strings.ToLower(dst_proc)

	rwLockRule.Lock()
	defer rwLockRule.Unlock()

	// 白名单放行
	_, ok := hMemRules.White[process]
	if ok {
		//xplog.LogInsertEvent("白名单", "防护模式", uname, process, dst_proc, "进程注入", "允许")
		return true
	}

	// 黑名单拒绝
	_, ok = hMemRules.Black[process]
	if ok {
		xplog.LogInsertEvent("黑名单", "防护模式", uname, process, dst_proc, "进程注入", "拒绝")
		return false
	}

	// 自保护
	if strings.Index(dst_proc, naseClientDir) >= 0 {
		return false
	}

	// 账户设置放行：特殊情况
	if strings.Index(dst_proc, seceditExe) >= 0 {
		return true
	}

	// 进程被注入
	if hMemRules.SafeHighCfg.ProcInject == 1 {
		// 功能开启
		if hMemRules.SafeHighCfg.Mode == 0 {
			xplog.LogInsertEvent("增强防护-防止进程被注入", "监视模式", uname, process, dst_proc, "进程注入", "拒绝")
			return true
		} else {
			xplog.LogInsertEvent("增强防护-防止进程被注入", "防护模式", uname, process, dst_proc, "进程注入", "拒绝")
			return false
		}
	}

	return true
}

// 规则匹配 - 服务篡改（Del,change）
func RuleMatchServiceChange(uname, process, service_name, op string) bool {
	process, _ = filepath.Abs(process)
	process = strings.ToLower(process)

	rwLockRule.Lock()
	defer rwLockRule.Unlock()

	// 白名单放行
	_, ok := hMemRules.White[process]
	if ok {
		//xplog.LogInsertEvent("白名单", "防护模式", uname, process, service_name, op, "允许")
		return true
	}

	// 黑名单拒绝
	_, ok = hMemRules.Black[process]
	if ok {
		xplog.LogInsertEvent("黑名单", "防护模式", uname, process, service_name, op, "拒绝")
		return false
	}

	// 防止篡改系统服务
	if hMemRules.SafeBaseCfg.WinService == 1 {
		// 功能开启
		if hMemRules.SafeBaseCfg.Mode == 0 {
			xplog.LogInsertEvent("基本防护-防止篡改系统服务", "监视模式", uname, process, service_name, op, "拒绝")
			return true
		} else {
			xplog.LogInsertEvent("基本防护-防止篡改系统服务", "防护模式", uname, process, service_name, op, "拒绝")
			return false
		}
	}
	return true
}

// 规则匹配 - 服务添加（Add）
func RuleMatchServiceAdd(uname, process, service_name, binPath string) bool {
	process, _ = filepath.Abs(process)
	process = strings.ToLower(process)
	binPath, _ = filepath.Abs(binPath)
	binPath = strings.ToLower(binPath)

	rwLockRule.Lock()
	defer rwLockRule.Unlock()

	logdst := "[" + service_name + "]" + binPath

	// 白名单放行
	_, ok := hMemRules.White[process]
	if ok {
		//xplog.LogInsertEvent("白名单", "防护模式", uname, process, logdst, "添加服务", "允许")
		return true
	}

	// 黑名单拒绝
	_, ok = hMemRules.Black[process]
	if ok {
		xplog.LogInsertEvent("黑名单", "防护模式", uname, process, logdst, "添加服务", "拒绝")
		return false
	}

	// 防止服务被添加
	if hMemRules.SafeHighCfg.AddService == 1 {
		// 功能开启
		if hMemRules.SafeHighCfg.Mode == 0 {
			xplog.LogInsertEvent("增强防护-防止服务被添加", "监视模式", uname, process, logdst, "添加服务", "拒绝")
			return true
		} else {
			xplog.LogInsertEvent("增强防护-防止服务被添加", "防护模式", uname, process, logdst, "添加服务", "拒绝")
			return false
		}
	}
	return true
}

// 规则匹配 - 驱动加载
func RuleMatchDriveLoad(uname, process, service_name, binPath string) bool {
	process, _ = filepath.Abs(process)
	process = strings.ToLower(process)
	binPath, _ = filepath.Abs(binPath)
	binPath = strings.ToLower(binPath)

	rwLockRule.Lock()
	defer rwLockRule.Unlock()

	logdst := "[" + service_name + "]" + binPath

	// 白名单放行
	_, ok := hMemRules.White[process]
	if ok {
		//xplog.LogInsertEvent("白名单", "防护模式", uname, process, logdst, "驱动加载", "允许")
		return true
	}

	// 黑名单拒绝
	_, ok = hMemRules.Black[process]
	if ok {
		xplog.LogInsertEvent("黑名单", "防护模式", uname, process, logdst, "驱动加载", "拒绝")
		return false
	}

	// 防止驱动程序被加载
	if hMemRules.SafeHighCfg.LoadSys == 1 {
		// 功能开启
		if hMemRules.SafeHighCfg.Mode == 0 {
			xplog.LogInsertEvent("增强防护-防止驱动程序被加载", "监视模式", uname, process, logdst, "驱动加载", "拒绝")
			return true
		} else {
			xplog.LogInsertEvent("增强防护-防止驱动程序被加载", "防护模式", uname, process, logdst, "驱动加载", "拒绝")
			return false
		}
	}
	return true
}

// 规则匹配 - 注册表设置 - 开机启动
func RuleMatchRegSetValue(uname, process, regpath, regvalue string) bool {
	process, _ = filepath.Abs(process)
	process = strings.ToLower(process)
	regpath = strings.ToUpper(regpath)

	rwLockRule.Lock()
	defer rwLockRule.Unlock()

	logdst := "[" + regvalue + "]" + regpath
	// 防止驱动程序被加载
	if hMemRules.SafeHighCfg.AddStart == 1 {
		// 白名单放行
		_, ok := hMemRules.White[process]
		if ok {
			//xplog.LogInsertEvent("白名单", "防护模式", uname, process, logdst, "设置开机启动", "允许")
			return true
		}

		// 黑名单拒绝
		_, ok = hMemRules.HighWinStart[regpath]
		if ok {
			xplog.LogInsertEvent("黑名单", "防护模式", uname, process, logdst, "设置开机启动", "拒绝")
			return false
		}

		for r, _ := range hMemRules.HighWinStart {
			if strings.Index(regpath, r) == 0 {
				// 访问启动项注册表
				if hMemRules.SafeHighCfg.Mode == 0 {
					xplog.LogInsertEvent("增强防护-防止开机自启动", "监视模式", uname, process, logdst, "设置开机启动", "拒绝")
					return true
				} else {
					xplog.LogInsertEvent("增强防护-防止开机自启动", "防护模式", uname, process, logdst, "设置开机启动", "拒绝")
					return false
				}
			}
		}
	}
	return true
}
