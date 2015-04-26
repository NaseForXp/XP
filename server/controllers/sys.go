package controllers

import (
	"encoding/json"
	"fmt"

	"../tools/rules"
	"github.com/astaxie/beego"
)

// 修改密码
type SysChangePasswordRequset struct {
	User   string // 用户名
	OldPwd string // 旧密码
	NewPwd string // 密码
}

// 修改密码响应
type SysChangePasswordResponse struct {
	Status int    // 1:成功 其他:失败
	Errmsg string // 错误原因
}

type SysController struct {
	beego.Controller
}

func (c *SysController) Get() {
	fmt.Println("---Login Get")
}

func (c *SysController) SysChangePassword() {
	fmt.Println("---SysChangePassword")
	fmt.Println("request :", c.GetString("data"))
	var req SysChangePasswordRequset
	var res SysChangePasswordResponse

	data := c.GetString("data")
	if data == "" {
		res.Status = 2
		res.Errmsg = "错误:数据data为空"
	} else {
		err := json.Unmarshal([]byte(data), &req)
		if err != nil {
			res.Status = 2
			res.Errmsg = "错误:参数格式错误" + data
		} else {
			//res = LoginCheck(req)
			//正常
			err := rules.RulesChangeUserPassword(req.User, req.OldPwd, req.NewPwd)
			if err != nil {
				res.Status = 2
				res.Errmsg = err.Error()
			} else {
				// 成功
				res.Status = 1
				res.Errmsg = "密码修改成功"
			}
		}
	}

	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["SysChangePassword_ret"] = string(jres)

	c.TplNames = "syscontroller/syschangepassword.tpl"
}
