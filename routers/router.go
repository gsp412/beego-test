// @APIVersion 1.0.0
// @Title 实验室管理系统API
// @Description 提供实验室管理系统相关API
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"beego-test/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("beego-test/v1",

		beego.NSBefore(controllers.Auth),

		beego.NSNamespace("/token",
			beego.NSInclude(
				&controllers.TokenController{},
			),
		),

		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
