package routers

import (
	"github.com/astaxie/beego"
	"influx_demo/controllers"
)

func init() {
	nsHome := beego.NewNamespace("/home",
		beego.NSRouter("/main", &controllers.MainController{}, "get:Main"),
	)
	nsInflux := beego.NewNamespace("/influx_db",
		beego.NSRouter("/write", &controllers.InfluxDbController{}, "post:Write"),
		beego.NSRouter("/query", &controllers.InfluxDbController{}, "get:Query"),
	)

	beego.AddNamespace(nsHome)
	beego.AddNamespace(nsInflux)
}
