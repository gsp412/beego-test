package lib

import (
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type Conf struct {
	Conn string
}

/* 连接池对象 */
type Pool struct {
	DbName string    // 注册数据库名称
	O      orm.Ormer // 注册数据库连接对象
}

/* 默认数据库名称 (保留) */
const DBNAME = "default"

/******************************************************************************
 **函数名称: RegisterDb
 **功    能: 注册数据库，并返回通用ormer，在使用事务时需要额外注册
 **输入参数:
 **    addr: 数据库连接地址
 **输出参数: NONE
 **返    回: 连接池数组
 **实现描述:
 **注意事项:
 **作    者: # guoshuangpeng@le.com # 2019-05-13 10:03:58 #
 ******************************************************************************/
func RegisterDb(addr string) {

	/* 注册数据库引擎 */
	// 参数1   driverName
	// 参数2   数据库类型
	// 这个用来设置 driverName 对应的数据库类型
	// mysql / sqlite3 / postgres 这三种是默认已经注册过的，所以可以无需设置
	//orm.RegisterDriver("mymysql", orm.DRMySQL)

	/* 注册数据库 */
	//orm.RegisterDataBase(DBNAME, "mysql", addr)

	/* 高级设置 设置最大空闲数、连接数 */
	maxIdle := 1000 // 设最大空闲连接
	maxConn := 2048 // 设置最大数据库连接 (go >= 1.2)
	orm.RegisterDataBase(DBNAME, "mysql", addr, maxIdle, maxConn)
	// orm.SetMaxIdleConns("default", maxIdle)  // 额外设置
	// orm.SetMaxOpenConns("default", maxConn)  // 额外设置

	/* 设置时区 */
	orm.DefaultTimeLoc, _ = time.LoadLocation("Asia/Shanghai")
}

/******************************************************************************
 **函数名称: GetMysqlPool
 **功    能: 获取通用ormer，在使用事务时需要额外注册
 **输入参数:
 **输出参数: NONE
 **返    回: 连接池数组
 **实现描述:
 **注意事项:
 **作    者: # guoshuangpeng@le.com # 2019-05-13 10:04:10 #
 ******************************************************************************/
func GetMysqlPool() *Pool {
	/* 数据库连接对象 */
	o := orm.NewOrm()
	o.Using(DBNAME)

	return &Pool{DbName: DBNAME, O: o}
}
