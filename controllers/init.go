package controllers

import (
	"beego-test/lib"
	"beego-test/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/json-iterator/go/extra"
)

/* 全局对象 */
var gApiCtx = &models.ApiCntx{}

func GetApiCntx() *models.ApiCntx {
	return gApiCtx
}

/******************************************************************************
 **函数名称: Init
 **功    能: 初始化处理
 **输入参数:
 **     cf: 配置信息
 **输出参数: NONE
 **返    回: 错误描述
 **实现描述:
 **注意事项:
 **作    者: # guoshuangpeng@le.com # 2019-05-13 09:59:44 #
 *****************************************************************************/
func Init(cf *models.ApiConf) (ctx *models.ApiCntx, err error) {
	ctx = GetApiCntx()
	ctx.Conf = cf

	/* 1.初始化日志 */
	logs.SetLogger(logs.AdapterFile, ctx.Conf.Log)
	//beego.BeeLogger.DelLogger("console")
	logs.EnableFuncCallDepth(true)

	/* 2.注册Mysql */
	lib.RegisterDb(cf.Mysql.Conn)
	lib.RegisterModel()            // 注册定义的Model
	ctx.Mysql = lib.GetMysqlPool() // 获取ORM

	if beego.BConfig.RunMode == beego.DEV {
		orm.Debug = true
	}

	/* jsoniter 全局设置 */
	// 开启容忍字符串数字互转
	extra.RegisterFuzzyDecoders()

	// 开启容忍空数组作为对象
	extra.RegisterFuzzyDecoders()

	return ctx, nil
}

/******************************************************************************
 **函数名称: Launch
 **功    能: 启动程序
 **输入参数: NONE
 **输出参数: NONE
 **返    回: 错误描述
 **实现描述: 启动后台工作协程
 **注意事项:
 **作    者: # guoshuangpeng@le.com # 2019-05-13 09:59:57 #
 ******************************************************************************/
func Launch(ctx *models.ApiCntx) (err error) {

	return nil
}
