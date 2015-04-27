package main

import (
	"fmt"

	_ "../server/routers"
	"./tools/rules"
	"github.com/astaxie/beego"
)

func main() {
	err := rules.RulesInit()
	if err != nil {
		fmt.Println(err)
		return
	}

	rules.RulesMemPrint()
	beego.Run()
}
