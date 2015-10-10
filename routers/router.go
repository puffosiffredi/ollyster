package routers

import (
	"log"
	"ollyster/controllers"

	"github.com/astaxie/beego"
)

func init() {
	//beego.Router("/", &controllers.MainController{})
	beego.Router("/:node([0-9]*)", &controllers.MainController{}, "get,post:HelloSitepoint")
	log.Println("Routers started! With website!")
}
