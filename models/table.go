package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

/* mysql数据库表名定义 */
const (
	API_TAB_USER   = "user"   // 用户信息表
	API_TAB_LAB    = "lab"    // 实验室信息表
	API_TAB_NOTICE = "notice" // 通知信息表
	API_TAB_COURSE = "course" // 课程信息表
	API_TAB_RECORD = "record" // 课程记录信息表
)

/* 数据库表注册 */
func RegisterModel() {
	orm.RegisterModel(new(User))
}

////////////////////////////////////////////////////////////////////////////////
// 用户信息表

const (
	USER_TYPE_ADMIN   = 1 // 实验室管理员
	USER_TYPE_TEACHER = 2 // 老师
	USER_TYPE_STUDENT = 3 // 学生
)

// @表名: user
// @描述: 存储用户基本信息
type User struct {
	Id         int64     // 用户ID
	StaffNo    string    // 工号
	Name       string    // 昵称
	Password   string    // 密码
	Type       int       // 用户类别
	CreateTime time.Time // 创建时间
	CreateId   int64     // 创建人
	UpdateTime time.Time // 更新时间
	UpdateId   int64     // 修改人
}
