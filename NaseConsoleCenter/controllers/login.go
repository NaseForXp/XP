package controllers

import (
	"encoding/json"
	"fmt"

	"../et99"
	"../tools"
	"github.com/astaxie/beego"
)

// 登录请求
type LoginRequset struct {
	User     string // 用户名
	Password string // 密码
}

// 登录响应
type LoginResponse struct {
	Status    int    // 1:成功 其他:失败
	Errmsg    string // 错误原因
	User      string // 用户名
	Usertype  int    // 用户类型 1:管理中心  2:admin  3：audit
	UserTokey string // 随机字符串，验证是否已经登录
}

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	fmt.Println("---Login Get")
}

func (c *LoginController) Login() {
	fmt.Println("---Login")
	fmt.Println("request :", c.GetString("data"))
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
	res.Status = 2
	if req.User == "" {
		res.Errmsg = "错误:用户名不能为空"
		return res
	}

	if req.Password == "" {
		res.Errmsg = "错误:密码不能为空"
		return res
	}

	// 验证USBkey
	err := et99.Et99_check_center_login()
	if err != nil {
		res.Errmsg = err.Error()
		return res
	}

	_, user_type, err := tools.RulesCheckUserPassword(req.User, req.Password)
	if err != nil {
		res.Errmsg = err.Error()
		return res
	}

	// 生成用户令牌
	tokey := tools.LoginCreateTokey()
	tools.LoginAddTokey(req.User, tokey)

	res.Errmsg = "登录成功"
	res.Status = 1
	res.User = req.User
	res.Usertype = user_type
	res.UserTokey = tokey

	return res
}
