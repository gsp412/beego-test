package controllers

import (
	"beego-test/lib"
	"beego-test/models"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// token相关接口
type TokenController struct {
	BaseController
}

type TokenRsp struct {
	Data    TokenData `json:"data" description:"业务数据"`
	Code    int       `json:"code" description:"详细错误码"`
	Message string    `json:"message" description:"错误描述"`
}

type TokenData struct {
	Id      int64  `json:"id" description:"用户ID"`
	Name    string `json:"name" description:"用户名称"`
	StaffNo string `json:"staff_no" description:"工号"`
	Type    int    `json:"type" description:"用户类别：1：管理员 2：老师 3：学生"`
	Tk      string `json:"tk" description:"用户TK"`
}

// @Title 用户获取token接口
// @Summary 用户获取token接口
// @Description 用户获取token接口
// @Param staff_no query string true "用户工号"
// @Param password query string true "用户密码"
// @Success 200 {object} controllers.TokenRsp
// @Failure 400
// @router / [get]
func (this *TokenController) Get() {
	ctx := GetApiCntx()

	// 获取参数
	staffNo := lib.RemoveEmpty(this.GetString("staff_no"))
	if 0 == len(staffNo) {
		logs.Error("staff_no can not be empty")
		this.ErrorMessage(lib.ERR_PARAM_INVALID, "工号不能为空")
		return
	}
	password := lib.RemoveEmpty(this.GetString("password"))
	if 0 == len(password) {
		logs.Error("password can not be empty")
		this.ErrorMessage(lib.ERR_PARAM_INVALID, "密码不能为空")
		return
	}

	// 获取用户信息
	user := &lib.User{
		StaffNo: staffNo,
	}
	err := ctx.Mysql.O.Read(user, "staff_no")
	if nil != err && orm.ErrNoRows != err {
		logs.Error("query user db exception, msg: %s", err.Error())
		this.ErrorMessage(lib.ERR_SYS_MYSQL, "数据库异常")
		return
	} else if orm.ErrNoRows == err {
		logs.Error("user not exits, staff_no: %s", staffNo)
		this.ErrorMessage(lib.ERR_PARAM_INVALID, "用户名或密码错误")
		return
	}

	// 验证密码
	_, err = models.MatchPassword(password, user.Password)
	if nil != err {
		logs.Error("password not match, msg: %s", err.Error())
		this.ErrorMessage(lib.ERR_PARAM_INVALID, "用户名或密码错误")
		return
	}

	// 创建token
	tk, code, err := ctx.CreateToken(user.Id, user.Password)
	if nil != err {
		logs.Error("create token failed, msg: %s", err.Error())
		this.ErrorMessage(code, "创建token失败")
		return
	}

	rsp := &TokenRsp{
		Data: TokenData{
			Id:      user.Id,
			Name:    user.Name,
			StaffNo: user.StaffNo,
			Type:    user.Type,
			Tk:      tk,
		},
		Code:    lib.OK,
		Message: "成功",
	}

	this.Data["json"] = rsp
	this.ServeJSON()
}
