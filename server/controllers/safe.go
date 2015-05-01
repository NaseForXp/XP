package controllers

import (
	"encoding/json"
	"fmt"

	"../tools/rules"
	"../tools/xplog"
	"github.com/astaxie/beego"
)

type SafeController struct {
	beego.Controller
}

func (c *SafeController) Get() {
	fmt.Println("---Safe Get")
}

// 基本防护 - 设置 - 请求
type SafeBaseSetRequest struct {
	Mode       int // 模式：0:监视模式 1:防护模式
	WinDir     int // 系统文件及目录保护状态 0:关闭 1:开启
	WinStart   int // 系统启动文件保护状态   0:关闭 1:开启
	WinFormat  int // 防止格式化磁盘状态    0:关闭 1:开启
	WinProc    int // 防止系统关键进程被杀死 0:关闭 1:开启
	WinService int // 防止篡改系统服务      0:关闭 1:开启
}

// 基本防护 - 获取设置 - 响应
type SafeBaseGetResponse struct {
	Status     int    // 1:成功 其他:失败
	Errmsg     string // 错误原因
	Mode       int    // 模式：0:监视模式 1:防护模式
	WinDir     int    // 系统文件及目录保护状态 0:关闭 1:开启
	WinStart   int    // 系统启动文件保护状态   0:关闭 1:开启
	WinFormat  int    // 防止格式化磁盘状态    0:关闭 1:开启
	WinProc    int    // 防止系统关键进程被杀死 0:关闭 1:开启
	WinService int    // 防止篡改系统服务      0:关闭 1:开启
}

// 增强防护 - 设置 - 请求
type SafeHighSetRequest struct {
	Mode       int // 模式：0:监视模式 1:防护模式
	AddService int // 防止服务被添加       0:关闭 1:开启
	AutoRun    int // 防止自动运行恶意程序  0:关闭 1:开启
	AddStart   int // 防止添加开机启动项    0:关闭 1:开启
	ReadWrite  int // 防止磁盘直接读写      0:关闭 1:开启
	CreateExe  int // 防止创建EXE文件      0:关闭 1:开启
	LoadSys    int // 防止驱动被加载        0:关闭 1:开启
	ProcInject int // 防止进程被注入        0:关闭 1:开启
}

// 增强防护 - 获取设置 - 响应
type SafeHighGetResponse struct {
	Status     int    // 1:成功 其他:失败
	Errmsg     string // 错误原因
	Mode       int    // 模式：0:监视模式 1:防护模式
	AddService int    // 防止服务被添加       0:关闭 1:开启
	AutoRun    int    // 防止自动运行恶意程序  0:关闭 1:开启
	AddStart   int    // 防止添加开机启动项    0:关闭 1:开启
	ReadWrite  int    // 防止磁盘直接读写      0:关闭 1:开启
	CreateExe  int    // 防止创建EXE文件      0:关闭 1:开启
	LoadSys    int    // 防止驱动被加载        0:关闭 1:开启
	ProcInject int    // 防止进程被注入        0:关闭 1:开启
}

// 基本防护 - 设置 - 响应
type SafeSetResponse struct {
	Status int    // 1:成功 其他:失败
	Errmsg string // 错误原因
}

// 基本防护 - 导出 - ini配置文件的内容
type SafeSaveResponse struct {
	Status int    // 1:成功 其他:失败
	Errmsg string // 错误原因
	Config string // 配置内容 string
}

// 基本防护 - 设置
func (c *SafeController) SafeBaseSet() {
	var req SafeBaseSetRequest
	var res SafeSetResponse

	usertokey := c.GetString("UserTokey")
	data := c.GetString("data")

	fmt.Println("---SafeBaseSet")
	fmt.Println("request :", usertokey, " | ", data)

	if LoginCheckTokeyJson(usertokey) == false {
		res.Status = 2
		res.Errmsg = "错误:请登录后操作"
		goto End
	}

	if data == "" {
		res.Status = 2
		res.Errmsg = "错误:数据data为空"
	} else {
		err := json.Unmarshal([]byte(data), &req)
		if err != nil {
			res.Status = 2
			res.Errmsg = "错误:参数格式错误" + data
		} else {
			//正常
			var base rules.SafeBaseConfig
			base.Mode = req.Mode
			base.WinDir = req.WinDir
			base.WinStart = req.WinStart
			base.WinFormat = req.WinFormat
			base.WinProc = req.WinProc
			base.WinService = req.WinService

			err := rules.RulesSafeBaseSet(base)
			if err != nil {
				res.Status = 2
				res.Errmsg = err.Error()
			} else {
				// 成功
				res.Status = 1
				res.Errmsg = "基本防护设置成功"
			}
		}
	}
End:
	if res.Status == 1 {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "设置基本防护配置", data, "成功")
	} else {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "设置基本防护配置", data, "失败")
	}
	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["Safe_ret"] = string(jres)

	c.TplNames = "safecontroller/safe.tpl"
}

// 基本防护 - 获取设置
func (c *SafeController) SafeBaseGet() {
	var res SafeBaseGetResponse

	usertokey := c.GetString("UserTokey")

	fmt.Println("---SafeBaseGet")
	fmt.Println("request :", usertokey)

	if LoginCheckTokeyJson(usertokey) == false {
		res.Status = 2
		res.Errmsg = "错误:请登录后操作"
		goto End
	} else {
		base, err := rules.RulesSafeBaseGet()
		if err != nil {
			res.Status = 2
			res.Errmsg = err.Error()
		} else {
			//正常
			res.Status = 1
			res.Errmsg = "获取基本防护设置成功"

			res.Mode = base.Mode
			res.WinDir = base.WinDir
			res.WinStart = base.WinStart
			res.WinFormat = base.WinFormat
			res.WinProc = base.WinProc
			res.WinService = base.WinService
		}
	}
End:
	if res.Status == 1 {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "获取基本防护配置", "", "成功")
	} else {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "获取基本防护配置", "", "失败")
	}
	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["Safe_ret"] = string(jres)
	c.TplNames = "safecontroller/safe.tpl"
}

// 基本防护 - 导出
func (c *SafeController) SafeBaseSave() {
	var res SafeSaveResponse

	usertokey := c.GetString("UserTokey")

	fmt.Println("---SafeBaseSave")
	fmt.Println("request :", usertokey)

	if LoginCheckTokeyJson(usertokey) == false {
		res.Status = 2
		res.Errmsg = "错误:请登录后操作"
		goto End
	} else {
		saveString, err := rules.RulesSafeBaseSave()
		if err != nil {
			res.Status = 2
			res.Errmsg = err.Error()
		} else {
			//正常
			res.Status = 1
			res.Errmsg = "基本防护规则导出成功"
			res.Config = saveString
		}
	}
End:

	if res.Status == 1 {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "导出基本防护配置", "", "成功")
	} else {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "导出基本防护配置", "", "失败")
	}

	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["Safe_ret"] = string(jres)
	c.TplNames = "safecontroller/safe.tpl"
}

// 增强防护 - 设置
func (c *SafeController) SafeHighSet() {
	var req SafeHighSetRequest
	var res SafeSetResponse

	usertokey := c.GetString("UserTokey")
	data := c.GetString("data")

	fmt.Println("---SafeHighSet")
	fmt.Println("request :", usertokey, " | ", data)

	if LoginCheckTokeyJson(usertokey) == false {
		res.Status = 2
		res.Errmsg = "错误:请登录后操作"
		goto End
	}

	if data == "" {
		res.Status = 2
		res.Errmsg = "错误:数据data为空"
	} else {
		err := json.Unmarshal([]byte(data), &req)
		if err != nil {
			res.Status = 2
			res.Errmsg = "错误:参数格式错误" + data
		} else {
			//正常
			var high rules.SafeHighConfig
			high.Mode = req.Mode
			high.AddService = req.AddService
			high.AutoRun = req.AutoRun
			high.AddStart = req.AddStart
			high.ReadWrite = req.ReadWrite
			high.CreateExe = req.CreateExe
			high.LoadSys = req.LoadSys
			high.ProcInject = req.ProcInject

			err := rules.RulesSafeHighSet(high)
			if err != nil {
				res.Status = 2
				res.Errmsg = err.Error()
			} else {
				// 成功
				res.Status = 1
				res.Errmsg = "设置增强防护配置成功"
			}
		}
	}
End:
	if res.Status == 1 {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "设置增强防护配置", data, "成功")
	} else {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "设置增强防护配置", data, "失败")
	}
	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["Safe_ret"] = string(jres)

	c.TplNames = "safecontroller/safe.tpl"
}

// 增强防护 - 获取设置
func (c *SafeController) SafeHighGet() {
	var res SafeHighGetResponse

	usertokey := c.GetString("UserTokey")

	fmt.Println("---SafeHighGet")
	fmt.Println("request :", usertokey)

	if LoginCheckTokeyJson(usertokey) == false {
		res.Status = 2
		res.Errmsg = "错误:请登录后操作"
		goto End
	} else {
		high, err := rules.RulesSafeHighGet()
		if err != nil {
			res.Status = 2
			res.Errmsg = err.Error()
		} else {
			//正常
			res.Status = 1
			res.Errmsg = "获取增强防护配置成功"

			res.Mode = high.Mode
			res.AddService = high.AddService
			res.AutoRun = high.AutoRun
			res.AddStart = high.AddStart
			res.ReadWrite = high.ReadWrite
			res.CreateExe = high.CreateExe
			res.LoadSys = high.LoadSys
			res.ProcInject = high.ProcInject
		}
	}
End:
	if res.Status == 1 {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "获取增强防护配置", "", "成功")
	} else {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "获取增强防护配置", "", "失败")
	}
	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["Safe_ret"] = string(jres)
	c.TplNames = "safecontroller/safe.tpl"
}

// 增强防护 - 导出
func (c *SafeController) SafeHighSave() {
	var res SafeSaveResponse

	usertokey := c.GetString("UserTokey")

	fmt.Println("---SafeHighSave")
	fmt.Println("request :", usertokey)

	if LoginCheckTokeyJson(usertokey) == false {
		res.Status = 2
		res.Errmsg = "错误:请登录后操作"
		goto End
	} else {
		saveString, err := rules.RulesSafeHighSave()
		if err != nil {
			res.Status = 2
			res.Errmsg = err.Error()
		} else {
			//正常
			res.Status = 1
			res.Errmsg = "导出增强防护配置成功"
			res.Config = saveString
		}
	}
End:
	if res.Status == 1 {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "导出增强防护配置", "", "成功")
	} else {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "导出增强防护配置", "", "失败")
	}
	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["Safe_ret"] = string(jres)
	c.TplNames = "safecontroller/safe.tpl"
}
