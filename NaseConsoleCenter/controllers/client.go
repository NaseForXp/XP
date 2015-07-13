package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
)

type ClientController struct {
	beego.Controller
}

// 登录时候将客户端信息发送给管理中心
type ClientInfomationRequest struct {
	IP   string // 客户端IP
	Port string // 客户端端口
}

// 响应
type ClientInfomationResponse struct {
	Status int    // 1:成功 其他:失败
	Errmsg string // 错误原因
}

func (c *ClientController) ClientAdd() {
	fmt.Println("---ClientAdd")
	fmt.Println("request :", c.GetString("data"))

	var req ClientInfomationRequest
	var res ClientInfomationResponse

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
			// 写入数据库
			res.Errmsg = "成功"
			res.Status = 1
		}
	}

	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["client_ret"] = string(jres)

	c.TplNames = "client.tpl"
}
