package controllers

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"../tools/rules"
	"../tools/toolcenter"
	"../tools/xplog"
	"github.com/astaxie/beego"
)

// 登录请求
type LoginRequset struct {
	User         string // 用户名
	Password     string // 密码
	CenterIPPort string // 管理中心IP、端口
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

	if res.Status == 1 {
		xplog.LogInsertSys(req.User, "登录", res.Errmsg, "成功")
	} else {
		xplog.LogInsertSys(req.User, "登录", res.Errmsg, "失败")
	}

	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["login_ret"] = string(jres)

	c.TplNames = "login.tpl"
}

func LoginCheckIpFormat(ip string) bool {
	bTrue, _ := regexp.MatchString(`^([1-9]|[1-9][0-9]|1[0-9]?[0-9]?|2[0-5][0-5]){1}(\.([0-9]|[1-9][0-9]|1[0-9]?[0-9]?|2[0-5][0-5])){3}$`, ip)
	return bTrue
}

func LoginCheckPortFormat(port string) bool {
	bTrue, _ := regexp.MatchString(`^\d{1,5}$`, port)
	return bTrue
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

	// 管理员需要验证IP+端口
	if req.User == "Admin" {
		ipport := strings.Split(req.CenterIPPort, ":")
		if len(ipport) != 2 {
			res.Errmsg = "错误:管理中心IP端口格式错误:" + req.CenterIPPort
			return res
		}

		ip := ipport[0]
		port := ipport[1]

		if LoginCheckIpFormat(ip) == false {
			res.Errmsg = "错误:管理中心IP格式错误:" + ip
			return res
		}

		if LoginCheckPortFormat(port) == false {
			res.Errmsg = "错误:管理中心Port格式错误:" + port
			return res
		}

		rip, rport, err := toolcenter.CenterGetIpPort()
		if err != nil {
			res.Errmsg = err.Error()
			return res
		}

		if ip != rip || port != rport {
			err = toolcenter.CenterSetIpPort(ip, port)
			if err != nil {
				res.Errmsg = err.Error()
				return res
			}
		}

	}

	_, user_type, err := rules.RulesCheckUserPassword(req.User, req.Password)
	if err != nil {
		res.Errmsg = err.Error()
		return res
	}

	if req.User == "Admin" {
		// 将客户端信息发送给管理中心，独立线程，防止等待
		go toolcenter.CenterSendClientInfo()
		/*
			err = toolcenter.CenterSendClientInfo()
			if err != nil {
				res.Errmsg = "管理中心无法连接"
				res.Status = 2
				//return res
			}
		*/
	}

	// 生成用户令牌
	tokey := LoginCreateTokey()
	LoginAddTokey(req.User, tokey)

	res.Errmsg = "登录成功"
	res.Status = 1
	res.User = req.User
	res.Usertype = user_type
	res.UserTokey = tokey

	return res
}
