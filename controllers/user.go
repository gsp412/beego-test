package controllers

import (
	"errors"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/json-iterator/go"

	"beego-test/lib"
	"beego-test/models"
)

// 用户相关接口
type UserController struct {
	BaseController
}

type UserPostRsp struct {
	Data    UserPostData `json:"data" description:"业务数据"`
	Code    int          `json:"code" description:"详细错误码"`
	Message string       `json:"message" description:"错误描述"`
}

type UserPostData struct {
	Id int64 `json:"id" description:"所添加的管理员ID"`
}

// @Title 用户添加接口(仅管理员可调用，初始密码为工号)
// @Summary 用户添加接口(仅管理员可调用，初始密码为工号)
// @Description 用户添加接口(仅管理员可调用，初始密码为工号)
// @Param body body controllers.UserReq true "请求参数"
// @Success 201 {object} controllers.UserPostRsp
// @Failure 400
// @router / [post]
func (this *UserController) Post() {
	ctx := GetApiCntx()

	admin, code, err := this.getAdmin()
	if nil != err {
		logs.Error("get admin failed, msg: %s", err.Error())
		this.ErrorMessage(code, err.Error())
		return
	} else if lib.USER_TYPE_ADMIN != admin.Type {
		logs.Error("no auth")
		this.ErrorMessage(lib.ERR_FORBIDDEN, "用户没有访问该接口权限")
		return
	}

	req, code, err := this.getUserReqParam()
	if nil != err {
		logs.Error("user post req format error, msg: ", err.Error())
		this.ErrorMessage(code, err.Error())
		return
	}

	pwd, err := models.CreatePassword(req.StaffNo)
	if nil != err {
		logs.Error("create password error, msg: ", err.Error())
		this.ErrorMessage(lib.ERR_INTERNAL_SERVER_ERROR, "生成密码失败")
		return
	}

	user := lib.User{
		StaffNo:    req.StaffNo, // 工号
		Name:       req.Name,    // 昵称
		Password:   pwd,         // 密码
		Type:       req.Type,    // 用户类别
		CreateTime: time.Now(),  // 创建时间
		CreateId:   admin.Id,    // 创建人
		UpdateTime: time.Now(),  // 更新时间
		UpdateId:   admin.Id,    // 修改人
	}

	id, err := ctx.Mysql.O.Insert(user)
	if nil != err {
		logs.Error("Insert user db exception, msg: ", err.Error())
		this.ErrorMessage(lib.ERR_SYS_MYSQL, "数据库异常")
		return
	}

	rsp := &UserPostRsp{
		Data: UserPostData{
			Id:      id,
		},
		Code:    lib.OK,
		Message: "OK",
	}

	this.Data["json"] = rsp
	this.ServeJSON()
}

// 通用请求
type UserReq struct {
	Name     string `json:"name" description:"用户昵称" required:"true"`
	StaffNo  string `json:"staff_no" description:"用户工号" required:"true"`
	Type     int    `json:"type" description:"用户类别" required:"true"`
}

// 获取管理员接口传入参数
func (this *UserController) getUserReqParam() (*UserReq, int, error) {
	ctx := GetApiCntx()

	req := &UserReq{}

	// 解析传入数据
	err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, req)
	if nil != err {
		logs.Error("Parameter format is invalid! body: ",
			this.Ctx.Input.RequestBody, " msg: ", err.Error())
		return nil, lib.ERR_PARAM_INVALID, errors.New("参数解析失败")
	}

	req.Name = lib.RemoveEmpty(req.Name)
	if 0 == len(req.Name) {
		logs.Error("name not exist")
		return nil, lib.ERR_PARAM_INVALID, errors.New("用户昵称不能为空")
	}
	req.StaffNo = lib.RemoveEmpty(req.StaffNo)
	if 0 == len(req.StaffNo) {
		logs.Error("StaffNo not exist")
		return nil, lib.ERR_PARAM_INVALID, errors.New("用户工号不能为空")
	}

	if lib.USER_TYPE_ADMIN != req.Type &&
		lib.USER_TYPE_TEACHER != req.Type && lib.USER_TYPE_STUDENT != req.Type {
		logs.Error("type not right")
		return nil, lib.ERR_PARAM_INVALID, errors.New("用户类型不合法")
	}

	user := &lib.User{
		StaffNo: req.StaffNo,
	}

	err = ctx.Mysql.O.Read(user, "staff_no")
	if nil != err && orm.ErrNoRows != err {
		logs.Error("query db exception, msg: %s", err.Error())
		return nil, lib.ERR_SYS_MYSQL, errors.New("数据库异常")
	} else if nil == err {
		logs.Error("staff has exist, staff_no: %s", req.StaffNo)
		return nil, lib.ERR_PARAM_INVALID, errors.New("用户已存在")
	}

	return req, lib.OK, nil
}