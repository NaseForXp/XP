package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"../tools/toolcenter"
	"github.com/astaxie/beego"
)

type CenterController struct {
	beego.Controller
}

func (c *CenterController) Get() {
	fmt.Println("---Center Get")
}

func (c *CenterController) CenterGetaddr() {
	fmt.Println("---CenterGetaddr")

	var res struct {
		CenterIP   string
		CenterPort int
	}

	ip, port, err := toolcenter.CenterGetIpPort()
	if err != nil {
		res.CenterIP = ""
		res.CenterPort = 0
	} else {
		res.CenterIP = ip
		res.CenterPort, _ = strconv.Atoi(port)
	}

	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["CenterAddr"] = string(jres)

	c.TplNames = "centercontroller/centeraddr.tpl"
}
