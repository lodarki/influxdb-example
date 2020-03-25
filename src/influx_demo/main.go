package main

import (
	_ "influx_demo/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}

