package main

import (
	"fmt"
	"runtime"

	"beego-test/controllers"
	"beego-test/controllers/conf"
	"beego-test/models"
	_ "beego-test/routers"

	"github.com/astaxie/beego"
)

/* 初始化 */
func _init() *models.ApiCntx {
	runtime.GOMAXPROCS(runtime.NumCPU())

	/* > 加载EXCHANGE配置 */
	conf, err := conf.Load()
	if nil != err {
		fmt.Printf("Load configuration failed! errmsg:%s\n", err.Error())
		return nil
	}

	/* > 初始化AUTH环境 */
	ctx, err := controllers.Init(conf)
	if nil != err {
		fmt.Printf("Initialize context failed! errmsg:%s\n", err.Error())
		return nil
	}

	return ctx
}

/* 注册回调 */
func register(ctx *models.ApiCntx) {
}

/* 启动服务 */
func launch(ctx *models.ApiCntx) {
	controllers.Launch(ctx)

	beego.Run()
}

func main() {

	/* > 初始化 */
	ctx := _init()
	if nil == ctx {
		fmt.Printf("Initialize context failed!\n")
		return
	}

	/* > 注册回调 */
	register(ctx)

	/* > 启动服务 */
	launch(ctx)

}
