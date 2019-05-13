package models

import "beego-test/lib"

type ApiConf struct {
	Log   string   // 日志配置
	Mysql lib.Conf // mysql 配置
}

type ApiCntx struct {
	Conf  *ApiConf  // 配置信息
	Mysql *lib.Pool // mysql 连接池
}
