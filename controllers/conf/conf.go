package conf

import (
	"errors"

	"beego-test/models"

	"github.com/astaxie/beego"
)

/******************************************************************************
 **函数名称: Load
 **功    能: 加载配置信息
 **输入参数:
 **     path: 配置路径
 **输出参数: NONE
 **返    回:
 **     conf: 配置信息
 **     err: 错误描述
实现描述:
 **注意事项:
 **作    者: # guoshuangpeng@le.com # 2019-04-15 15:33:29 #
 ******************************************************************************/
func Load() (conf *models.ApiConf, err error) {
	conf = &models.ApiConf{}

	// 加载log相关
	conf.Log = beego.AppConfig.String("log_conf")
	if 0 == len(conf.Log) {
		return nil, errors.New("get log conf failed")
	}

	// 加载数据库相关
	conf.Mysql.Conn = beego.AppConfig.String("sql_conn")
	if 0 == len(conf.Mysql.Conn) {
		return nil, errors.New("get mysql addr failed")
	}

	return conf, err
}
