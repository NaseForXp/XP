package main

import (
	_ "../server/routers"
	"./tools/rules"
	"fmt"
	"github.com/astaxie/beego"
)

func main() {
	err := rules.RulesInit()
	if err != nil {
		fmt.Println(err)
		return
	}
	beego.Run()
}
