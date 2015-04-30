package main

import (
	"fmt"

	_ "../server/routers"
	"./tools/rules"
	"./tools/xplog"
	"github.com/astaxie/beego"
)

func main() {
	err := xplog.LogInit()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = rules.RulesInit()
	if err != nil {
		fmt.Println(err)
		xplog.LogFini()
		return
	}

	rules.RulesMemPrint()
	beego.Run()
}
