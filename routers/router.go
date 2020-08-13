// @APIVersion 1.0.0
// @Title 试卷库后台管理系统
// @Controller
// @APICode 10000 操作成功
// @APICode 10001 操作失败，未知错误
// @APICode 10002 参数错误
// @APICode 13300 资源未找到

package routers

import (
	"knowtech/controllers"

	"github.com/astaxie/beego"
)

func init() {

	beego.Router("/", &controllers.HealthCheckController{}, "get:HealthCheck")

	beego.Router("/login", &controllers.LoginController{}, "*:LoginIn")
	beego.Router("/login_out", &controllers.LoginController{}, "*:LoginOut")
	beego.Router("/home", &controllers.HomeController{}, "*:Index")
	beego.Router("/home/start", &controllers.HomeController{}, "*:Start")

	beego.AutoRouter(&controllers.PaperController{})
	beego.AutoRouter(&controllers.ChapterController{})
	beego.AutoRouter(&controllers.QuestionController{})
	beego.AutoRouter(&controllers.CheckController{})
	beego.AutoRouter(&controllers.CollectController{})
	beego.AutoRouter(&controllers.AdminController{})
	beego.AutoRouter(&controllers.TempController{})
	beego.AutoRouter(&controllers.CustomController{})
	beego.AutoRouter(&controllers.RecommendController{})
}
