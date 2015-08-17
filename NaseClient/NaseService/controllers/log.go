package controllers

import (
	"encoding/json"
	"fmt"
	"os"

	"../tools/rules"
	"../tools/serial"
	"../tools/xplog"

	"github.com/astaxie/beego"
)

type LogController struct {
	beego.Controller
}

func (c *LogController) Get() {
	fmt.Println("---Log Get")
}

// 系统日志 - 查询 - 请求
type LogQuerySysRequest struct {
	KeyWord   string // 关键词
	TimeStart string // 起始时间
	TimeStop  string // 结束时间
	Start     int    // 日志起始位置
	Length    int    // 日志条数
}

// 系统日志 - 查询 - 响应
type LogQuerySysResponse struct {
	Status   int                    // 1:成功 其他:失败
	Errmsg   string                 // 错误原因
	LogArray []xplog.LogSysQueryRes // 日志数组
}

// 安全日志 - 查询 - 请求
type LogQueryEventRequest struct {
	KeyWord   string // 关键词
	TimeStart string // 起始时间
	TimeStop  string // 结束时间
	Start     int    // 日志起始位置
	Length    int    // 日志条数
}

// 安全日志 - 查询 - 响应
type LogQueryEventResponse struct {
	Status   int                      // 1:成功 其他:失败
	Errmsg   string                   // 错误原因
	LogArray []xplog.LogEventQueryRes // 日志数组
}

// 日志 - 数量 - 响应
type LogTotleResponse struct {
	Status int    // 1:成功 其他:失败
	Errmsg string // 错误原因
	Count  int    // 日志数量
}

// 日志 - 查询首页统计信息 - 响应
type LogHomeCountResponse struct {
	Status         int    // 1:成功 其他:失败
	Errmsg         string // 错误原因
	BaseMode       int    // 基础防护模式
	HighMode       int    // 增强防护模式
	Totle          int    // 总数
	White          int    // 白名单事件数量
	Black          int    // 黑名单事件数量
	BaseWinDir     int    // 基本防护-系统文件及目录保护
	BaseWinStart   int    // 基本防护-系统启动文件保护
	BaseWinFormat  int    // 基本防护-防止格式化系统磁盘
	BaseWinProc    int    // 基本防护-防止系统关键进程被杀死
	BaseWinService int    // 基本防护-防止篡改系统服务
	HighAddService int    // 增强防护-防止服务被添加
	HighAutoRun    int    // 增强防护-防止自动运行
	HighAddStart   int    // 增强防护-防止开机自启动
	HighReadWrite  int    // 增强防护-防止磁盘被直接读写
	HighCreateExe  int    // 增强防护-禁止创建.exe文件
	HighLoadSys    int    // 增强防护-防止驱动程序被加载
	HighProcInject int    // 增强防护-防止进程被注入
}

// 日志导出 - 请求
type LogExportRequest struct {
	SaveDir string // 保存目录
}

// 日志导出 - 响应
type LogExportResponse struct {
	Status   int    // 1:成功 其他:失败
	Errmsg   string // 错误原因
	SaveFile string // 保存的文件全路径
}

// 日志 - 查询 - 系统日志数量
func (c *LogController) LogSysTotle() {
	var res LogTotleResponse

	usertokey := c.GetString("UserTokey")

	fmt.Println("---LogSysTotle")
	fmt.Println("request :", usertokey)

	if LoginCheckTokeyJson(usertokey) == false {
		res.Status = 2
		res.Errmsg = "错误:请登录后操作"
		goto End
	} else {
		//正常
		tot, err := xplog.LogQuerySysTotle()
		if err != nil {
			res.Status = 2
			res.Errmsg = err.Error()
		} else {
			// 成功
			res.Status = 1
			res.Errmsg = "查询:系统日志数量成功"
			res.Count = tot
		}
	}
End:
	if res.Status == 1 {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "查询:系统日志数量", "", "成功")
	} else {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "查询:系统日志数量", "", "失败")
	}
	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["Log_ret"] = string(jres)

	c.TplNames = "logcontroller/log.tpl"
}

// 日志 - 查询 - 安全日志数量
func (c *LogController) LogEventTotle() {
	var res LogTotleResponse

	usertokey := c.GetString("UserTokey")

	fmt.Println("---LogEventTotle")
	fmt.Println("request :", usertokey)

	if LoginCheckTokeyJson(usertokey) == false {
		res.Status = 2
		res.Errmsg = "错误:请登录后操作"
		goto End
	} else {
		//正常
		tot, err := xplog.LogQueryEventTotle()
		if err != nil {
			res.Status = 2
			res.Errmsg = err.Error()
		} else {
			// 成功
			res.Status = 1
			res.Errmsg = "查询:安全日志数量成功"
			res.Count = tot
		}
	}
End:
	if res.Status == 1 {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "查询:安全日志数量成功", "", "成功")
	} else {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "查询:安全日志数量成功", "", "失败")
	}
	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["Log_ret"] = string(jres)

	c.TplNames = "logcontroller/log.tpl"
}

// 日志 - 查询 - 系统日志
func (c *LogController) LogSysQuery() {
	var req LogQuerySysRequest
	var res LogQuerySysResponse

	usertokey := c.GetString("UserTokey")
	data := c.GetString("data")

	fmt.Println("---LogSysQuery")
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
			array, err := xplog.LogQuerySys(req.KeyWord, req.TimeStart, req.TimeStop, req.Start, req.Length)
			if err != nil {
				res.Status = 2
				res.Errmsg = err.Error()
			} else {
				// 成功
				res.Status = 1
				res.Errmsg = "查询:系统日志成功"
				res.LogArray = array
			}
		}
	}
End:
	if res.Status == 1 {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "查询:系统日志", data, "成功")
	} else {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "查询:系统日志", data, "失败")
	}

	jres, err := json.Marshal(res)
	fmt.Println("response:", err)
	c.Data["Log_ret"] = string(jres)

	c.TplNames = "logcontroller/log.tpl"
}

// 日志 - 查询 - 安全日志
func (c *LogController) LogEventQuery() {
	var req LogQueryEventRequest
	var res LogQueryEventResponse

	usertokey := c.GetString("UserTokey")
	data := c.GetString("data")

	fmt.Println("---LogEventQuery")
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
			array, err := xplog.LogQueryEvent(req.KeyWord, req.TimeStart, req.TimeStop, req.Start, req.Length)
			if err != nil {
				res.Status = 2
				res.Errmsg = err.Error()
			} else {
				// 成功
				res.Status = 1
				res.Errmsg = "查询:安全日志成功"
				res.LogArray = array
			}
		}
	}
End:
	if res.Status == 1 {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "查询:安全日志", data, "成功")
	} else {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "查询:安全日志", data, "失败")
	}
	jres, err := json.Marshal(res)
	fmt.Println("response:", err)
	c.Data["Log_ret"] = string(jres)

	c.TplNames = "logcontroller/log.tpl"
}

// 日志 - 查询 - 首页统计信息
func (c *LogController) LogHomeCount() {
	var res LogHomeCountResponse

	usertokey := c.GetString("UserTokey")

	fmt.Println("---LogHomeCount")
	fmt.Println("request :", usertokey)

	if LoginCheckTokeyJson(usertokey) == false {
		res.Status = 2
		res.Errmsg = "错误:请登录后操作"
		goto End
	} else {
		//正常
		homeCnt, err := xplog.LogQueryHomeCount()
		if err != nil {
			res.Status = 2
			res.Errmsg = err.Error()
		} else {
			// 成功
			res.Totle = homeCnt.Totle
			res.White = homeCnt.White
			res.Black = homeCnt.Black
			res.BaseWinDir = homeCnt.BaseWinDir
			res.BaseWinStart = homeCnt.BaseWinStart
			res.BaseWinFormat = homeCnt.BaseWinFormat
			res.BaseWinProc = homeCnt.BaseWinProc
			res.BaseWinService = homeCnt.BaseWinService
			res.HighAddService = homeCnt.HighAddService
			res.HighAutoRun = homeCnt.HighAutoRun
			res.HighAddStart = homeCnt.HighAddStart
			res.HighReadWrite = homeCnt.HighReadWrite
			res.HighCreateExe = homeCnt.HighCreateExe
			res.HighLoadSys = homeCnt.HighLoadSys
			res.HighProcInject = homeCnt.HighProcInject
			res.BaseMode, res.HighMode = rules.RulesMemGetHomeMode()
			res.Status = 1
			res.Errmsg = "查询:首页统计信息成功"
		}
	}
End:
	if res.Status == 1 {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "查询:首页统计信息", "", "成功")
	} else {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "查询:首页统计信息", "", "失败")
	}
	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["Log_ret"] = string(jres)

	c.TplNames = "logcontroller/log.tpl"
}

// 日志 - 查询 - 首页统计信息FushionCharts画图
func (c *LogController) LogHomeCountCharts() {
	type LableValue struct {
		Label string
		Value int
	}

	var res LogHomeCountResponse
	var data []LableValue
	var homeCnt xplog.LogHomeCount
	var err error

	usertokey := c.GetString("UserTokey")

	fmt.Println("---LogHomeCount")
	fmt.Println("request :", usertokey)

	if LoginCheckTokeyJson(usertokey) == false {
		res.Status = 2
		res.Errmsg = "错误:请登录后操作"
		goto End
	} else {
		//正常
		homeCnt, err = xplog.LogQueryHomeCount()
		if err != nil {
			res.Status = 2
			res.Errmsg = err.Error()
		} else {
			// 成功
			res.Status = 1
			res.Errmsg = "查询:首页统计信息成功"
			res.Totle = homeCnt.Totle
			res.BaseMode, res.HighMode = rules.RulesMemGetHomeMode()
		}
	}
End:
	if res.Status == 1 {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "查询:首页统计信息", "", "成功")
	} else {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "查询:首页统计信息", "", "失败")
	}

	err = serial.ClientVerifyLicense()
	if err == nil {
		c.Data["IsReg"] = "已经注册"
	} else {
		c.Data["IsReg"] = "尚未注册"
	}
	c.Data["EventTotle"] = res.Totle
	c.Data["BaseMode"] = "监视模式"
	if res.BaseMode == 1 {
		c.Data["BaseMode"] = "防护模式"
	}
	c.Data["HighMode"] = "监视模式"
	if res.HighMode == 1 {
		c.Data["HighMode"] = "防护模式"
	}

	//data = append(data, LableValue{"白名单", homeCnt.White})
	data = append(data, LableValue{"黑名单", homeCnt.Black})
	data = append(data, LableValue{"系统文件及目录保护", homeCnt.BaseWinDir})
	data = append(data, LableValue{"系统启动文件保护", homeCnt.BaseWinStart})
	data = append(data, LableValue{"防止格式化系统磁盘", homeCnt.BaseWinFormat})
	data = append(data, LableValue{"防止系统关键进程被杀死", homeCnt.BaseWinProc})
	data = append(data, LableValue{"防止篡改系统服务", homeCnt.BaseWinService})
	data = append(data, LableValue{"防止服务被添加", homeCnt.HighAddService})
	data = append(data, LableValue{"防止自动运行", homeCnt.HighAutoRun})
	data = append(data, LableValue{"防止开机自启动", homeCnt.HighAddStart})
	data = append(data, LableValue{"防止磁盘被直接读写", homeCnt.HighReadWrite})
	data = append(data, LableValue{"禁止创建.exe文件", homeCnt.HighCreateExe})
	data = append(data, LableValue{"防止驱动程序被加载", homeCnt.HighLoadSys})
	data = append(data, LableValue{"防止进程被注入", homeCnt.HighProcInject})
	c.Data["Data"] = data

	c.TplNames = "logcontroller/homecount.html"
}

// 日志 - 导出到文件
func (c *LogController) LogExport() {
	var req LogExportRequest
	var res LogExportResponse

	usertokey := c.GetString("UserTokey")
	data := c.GetString("data")

	fmt.Println("---LogExport")
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
			//saveName, err := xplog.LogExport(req.SaveDir)
			sysRoot := os.Getenv("SystemRoot")
			saveName, err := xplog.LogExport(sysRoot[0:3])
			if err != nil {
				res.Status = 2
				res.Errmsg = err.Error()
			} else {
				// 成功
				res.Status = 1
				res.Errmsg = "日志导出成功:" + saveName
				res.SaveFile = saveName
			}
		}
	}
End:
	if res.Status == 1 {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "日志导出", res.SaveFile, "成功")
	} else {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "日志导出", res.SaveFile, "失败")
	}

	jres, err := json.Marshal(res)
	fmt.Println("response:", err)
	c.Data["Log_ret"] = string(jres)

	c.TplNames = "logcontroller/log.tpl"
}
