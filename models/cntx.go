package models

type ApiConf struct {
	Log   string // 日志配置
	Mysql Conf   // mysql 配置
}

type ApiCntx struct {
	Conf  *ApiConf // 配置信息
	Mysql *Pool    // mysql 连接池
}
