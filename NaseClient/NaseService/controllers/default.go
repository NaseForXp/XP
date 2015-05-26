package controllers

import (
	"fmt"

	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplNames = "index.tpl"
}

func (c *MainController) Test() {
	var xmldata map[string]string
	xmldata = make(map[string]string)
	xmldata["Type"] = "column2d" // pie3d
	xmldata["Width"] = "500"
	xmldata["Height"] = "300"
	xmldata["Caption"] = "总体概况"
	xmldata["SubCaption"] = "分类统计图"
	xmldata["XAxisName"] = "类别"
	xmldata["YAxisName"] = "数量"

	type LableValue struct {
		Label string
		Value int
	}
	var data []LableValue
	for i := 0; i < 10; i++ {
		var lv LableValue
		lv.Label = fmt.Sprintf("Lab%d", i+1)
		lv.Value = i + 123
		data = append(data, lv)
	}

	c.Data["LVData"] = data
	c.Data["XmlData"] = xmldata
	c.TplNames = "chars.html"
}
