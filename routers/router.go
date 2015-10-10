package routers

import (
	"log"
	"ollyster/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/olly", &controllers.MainController{}, "get:HelloSitepoint")
	log.Println("Routers started!")
}
