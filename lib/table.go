package lib

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
	orm.RegisterModel(new(User), new(Lab), new(Notice), new(Course), new(Record))
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

////////////////////////////////////////////////////////////////////////////////
// 实验室信息表
// @表名: lab
// @描述: 存储实验室基本信息
type Lab struct {
	Id         int64     // Id
	Name       string    // 名称
	Location   string    // 地址
	CreateTime time.Time // 创建时间
	CreateId   int64     // 创建人
	UpdateTime time.Time // 更新时间
	UpdateId   int64     // 修改人
}

////////////////////////////////////////////////////////////////////////////////
// 通知信息表
// @表名: notice
// @描述: 存储通知信息
type Notice struct {
	Id          int64     // Id
	Title       string    // 标题
	Description string    // 描述信息
	CreateTime  time.Time // 创建时间
	CreateId    int64     // 创建人
	UpdateTime  time.Time // 更新时间
	UpdateId    int64     // 修改人
}

////////////////////////////////////////////////////////////////////////////////
// 课程信息表

const (
	COURSE_STATUS_RESERVING = 1 // 预约中
	COURSE_STATUS_SUCC      = 2 // 预约成功
	COURSE_STATUS_FAIL      = 3 // 预约失败
	COURSE_STATUS_CANCEL    = 4 // 取消
)

// @表名: course
// @描述: 存储课程信息
type Course struct {
	Id         int64     // Id
	Name       string    // 课程名称
	UserId     int64     // 教师ID
	LabId      int64     // 实验室ID
	StartTime  time.Time // 开始时间
	EndTime    time.Time // 结束时间
	Total      int       // 课程容纳人数
	Count      int       // 课程当前人数
	Status     int       // 预约状态
	CreateTime time.Time // 创建时间
	CreateId   int64     // 创建人
	UpdateTime time.Time // 更新时间
	UpdateId   int64     // 修改人
}

////////////////////////////////////////////////////////////////////////////////
// 学生课堂记录信息表

const (
	RECORD_STATUS_RESERVING = 1 // 预约中
	RECORD_STATUS_SUCC      = 2 // 预约成功
	RECORD_STATUS_FAIL      = 3 // 预约失败
	RECORD_STATUS_CANCEL    = 4 // 取消
)

// @表名: record
// @描述: 存储学生课堂记录信息信息
type Record struct {
	Id         int64     // Id
	UserId     int64     // 学生ID
	CourseId   int64     // 课程ID
	Status     int       // 预约状态
	Score      float64   // 得分
	CreateTime time.Time // 创建时间
	CreateId   int64     // 创建人
	UpdateTime time.Time // 更新时间
	UpdateId   int64     // 修改人
}
