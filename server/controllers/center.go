package controllers

import (
	"encoding/json"
	"fmt"

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

	res.CenterIP = "192.168.1.100"
	res.CenterPort = 8080

	jres, err := json.Marshal(res)
	fmt.Println("response:", string(jres), err)
	c.Data["CenterAddr"] = string(jres)

	c.TplNames = "centercontroller/centeraddr.tpl"
}
