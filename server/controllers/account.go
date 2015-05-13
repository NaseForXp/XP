package controllers

import (
	"encoding/json"
	"fmt"

	"../tools/rules"
	"../tools/xplog"
	"github.com/astaxie/beego"
)

type AccountController struct {
	beego.Controller
}

func (c *AccountController) Get() {
	fmt.Println("---Safe Get")
}

// 账户安全 - 设置 - 请求
type AccountSetRequest struct {
	Mode                  int // 模式：0:关闭 1:开启
	SafeLev               int // 账户策略设置 0:自定义 1:低级 2:中级 3:高级
	PasswordComplexity    int // 密码复杂度  0:关闭 1:开启
	MinimumPasswordLength int // 密码最小长度(字符个数)
	MinimumPasswordAge    int // 最短使用期限(天)
	MaximumPasswordAge    int // 最长使用期限(天)
	PasswordHistorySize   int // 强制密码历史次数(次)
	LockoutBadCount       int // 账户锁定次数(无效登录次数)
	LockoutDuration       int // 账户锁定时长(分钟)
}

// 账户安全 - 设置 - 响应
type AccountSetResponse struct {
	Status int    // 1:成功 其他:失败
	Errmsg string // 错误原因
}

// 账户安全 - 获取设置 - 响应
type AccountGetResponse struct {
	Status                int    // 1:成功 其他:失败
	Errmsg                string // 错误原因
	Mode                  int    // 模式：0:关闭 1:开启
	SafeLev               int    // 账户策略设置 0:自定义 1:低级 2:中级 3:高级
	PasswordComplexity    int    // 密码复杂度  0:关闭 1:开启
	MinimumPasswordLength int    // 密码最小长度(字符个数)
	MinimumPasswordAge    int    // 最短使用期限(天)
	MaximumPasswordAge    int    // 最长使用期限(天)
	PasswordHistorySize   int    // 强制密码历史次数(次)
	LockoutBadCount       int    // 账户锁定次数(无效登录次数)
	LockoutDuration       int    // 账户锁定时长(分钟)
}

// 账户安全 - 导出 - ini配置文件的内容
type AccountSaveResponse struct {
	Status int    // 1:成功 其他:失败
	Errmsg string // 错误原因
	Config string // 配置内容 string
}

// 账户安全 - 设置
func (c *AccountController) AccountSet() {
	var req AccountSetRequest
	var res AccountSetResponse

	usertokey := c.GetString("UserTokey")
	data := c.GetString("data")

	fmt.Println("---AccountSet")
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
			var account rules.AccountConfig

			account.Mode = req.Mode
			account.SafeLev = req.SafeLev
			account.PasswordComplexity = req.PasswordComplexity
			account.MinimumPasswordLength = req.MinimumPasswordLength
			account.MinimumPasswordAge = req.MinimumPasswordAge
			account.MaximumPasswordAge = req.MaximumPasswordAge
			account.PasswordHistorySize = req.PasswordHistorySize
			account.LockoutBadCount = req.LockoutBadCount
			account.LockoutDuration = req.LockoutDuration

			err := rules.RulesAccountSet(account)
			if err != nil {
				res.Status = 2
				res.Errmsg = err.Error()
			} else {
				// 成功
				res.Status = 1
				res.Errmsg = "设置账户安全配置成功"
			}
		}
	}
End:
	if res.Status == 1 {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "设置账户安全配置", data, "成功")
	} else {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "设置账户安全配置", data, "失败")
	}
	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["Account_ret"] = string(jres)
	c.TplNames = "accountcontroller/account.tpl"
}

// 账户安全 - 获取设置
func (c *AccountController) AccountGet() {
	var res AccountGetResponse

	usertokey := c.GetString("UserTokey")

	fmt.Println("---AccountGet")
	fmt.Println("request :", usertokey)

	if LoginCheckTokeyJson(usertokey) == false {
		res.Status = 2
		res.Errmsg = "错误:请登录后操作"
		goto End
	} else {
		account, err := rules.RulesAccountGet()
		if err != nil {
			res.Status = 2
			res.Errmsg = err.Error()
		} else {
			//正常
			res.Status = 1
			res.Errmsg = "获取账户安全配置成功"

			res.Mode = account.Mode
			res.SafeLev = account.SafeLev
			res.PasswordComplexity = account.PasswordComplexity
			res.MinimumPasswordLength = account.MinimumPasswordLength
			res.MinimumPasswordAge = account.MinimumPasswordAge
			res.MaximumPasswordAge = account.MaximumPasswordAge
			res.PasswordHistorySize = account.PasswordHistorySize
			res.LockoutBadCount = account.LockoutBadCount
			res.LockoutDuration = account.LockoutDuration
		}
	}
End:
	if res.Status == 1 {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "获取账户安全配置", "", "成功")
	} else {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "获取账户安全配置", "", "失败")
	}
	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["Account_ret"] = string(jres)
	c.TplNames = "accountcontroller/account.tpl"
}

// 账户安全 - 导出
func (c *AccountController) AccountSave() {
	var res AccountSaveResponse

	usertokey := c.GetString("UserTokey")

	fmt.Println("---AccountSave")
	fmt.Println("request :", usertokey)

	if LoginCheckTokeyJson(usertokey) == false {
		res.Status = 2
		res.Errmsg = "错误:请登录后操作"
		goto End
	} else {
		saveString, err := rules.RulesAccountSave()
		if err != nil {
			res.Status = 2
			res.Errmsg = err.Error()
		} else {
			//正常
			res.Status = 1
			res.Errmsg = "导出账户安全配置成功"
			res.Config = saveString
		}
	}
End:
	if res.Status == 1 {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "导出账户安全配置", "", "成功")
	} else {
		xplog.LogInsertSys(LoginGetUserByTokey(usertokey), "导出账户安全配置", "", "失败")
	}
	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["Account_ret"] = string(jres)
	c.TplNames = "accountcontroller/account.tpl"
}
