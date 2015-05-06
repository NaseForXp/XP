package controllers

import (
	"encoding/json"
	"fmt"

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
			array, err := xplog.LogQuerySys(req.KeyWord, req.TimeStart, req.TimeStop)
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
			array, err := xplog.LogQueryEvent(req.KeyWord, req.TimeStart, req.TimeStop)
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
