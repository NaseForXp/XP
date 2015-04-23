package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
)

// 登录请求
type LoginRequset struct {
	User       string // 用户名
	Password   string // 密码
	CenterIP   string // 管理中心IP
	CenterPort string // 管理中心端口
}

// 登录响应
type LoginResponse struct {
	Status   int    // 1:成功 其他:失败
	Errmsg   string // 错误原因
	User     string // 用户名
	Usertype int    // 用户类型 1:管理中心  2:admin  3：audit
}

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	fmt.Println("---Login Get")
}

func (c *LoginController) Login() {
	fmt.Println("---Login")
	var req LoginRequset
	var res LoginResponse

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
			res = LoginCheck(req)
		}
	}

	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["login_ret"] = string(jres)

	c.TplNames = "login.tpl"
}

func LoginCheck(req LoginRequset) (res LoginResponse) {
	fmt.Printf("Login: user=[%s], pwd=[%s], center=[%s:%s]\n", req.User, req.Password, req.CenterIP, req.CenterPort)
	res.Status = 2
	if req.User == "" {
		res.Errmsg = "错误:用户名不能为空"
		return res
	}

	if req.Password == "" {
		res.Errmsg = "错误:密码不能为空"
		return res
	}

	if req.CenterIP == "" {
		res.Errmsg = "错误:管理中心IP不能为空"
		return res
	}

	if req.CenterPort == "" {
		res.Errmsg = "错误:管理中心端口不能为空"
		return res
	}

	res.Errmsg = "登录成功"
	res.Status = 1
	res.User = req.User
	res.Usertype = 2
	return res
}
