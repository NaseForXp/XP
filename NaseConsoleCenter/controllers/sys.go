package controllers

import (
	"encoding/json"

	"../tools"
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
	tools.Println("---Login Get")
}

func (c *SysController) SysChangePassword() {
	var req SysChangePasswordRequset
	var res SysChangePasswordResponse

	usertokey := c.GetString("UserTokey")
	data := c.GetString("data")

	tools.Println("---SysChangePassword")
	tools.Println("request :", usertokey, " | ", data)

	if tools.LoginCheckTokeyJson(usertokey) == false {
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
			err := tools.RulesChangeUserPassword(req.User, req.OldPwd, req.NewPwd)
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

End:
	if res.Status == 1 {
		//
	} else {
		//
	}
	jres, err := json.Marshal(res)
	tools.Println("response:", string(jres), err)
	c.Data["Sys_ret"] = string(jres)

	c.TplNames = "sys.tpl"
}
