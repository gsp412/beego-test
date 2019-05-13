package controllers

import (
	"errors"
	"strconv"
	"time"
	"vcoin-api/src/lib/comm"

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

type UserGetRsp struct {
	Data    UserGetData `json:"data" description:"业务数据"`
	Code    int         `json:"code" description:"详细错误码"`
	Message string      `json:"message" description:"错误描述"`
}

type UserGetData struct {
	Pages    int           `json:"pages" description:"总页数"`
	PageSize int           `json:"page_size" description:"每页条目数"`
	Total    int64         `json:"total" description:"总条目数"`
	Len      int           `json:"len" description:"列表总长度"`
	List     []UserGetItem `json:"list" description:"列表数据"`
}

type UserGetItem struct {
	Id      int64  `json:"id"  description:"用户ID"`
	StaffNo string `json:"staff_no"  description:"用户工号"`
	Name    string `json:"name"  description:"用户名称"`
	Type    int    `json:"type"  description:"用户类型"`
}

// @Title 用户查询接口(仅管理员可调用)
// @Summary 用户查询接口(仅管理员可调用)
// @Description 用户查询接口(仅管理员可调用)
// @Param name query string false "用户名称：支持模糊查询"
// @Param type query number false "用户类型"
// @Param staff_no query string false "用户工号"
// @Param page query number false "分页号：默认第一页"
// @Param page_size query number false "分页大小: 默认10条"
// @Success 201 {object} controllers.UserGetRsp
// @Failure 400
// @router / [get]
func (this *UserController) Get() {
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

	q, code, err := this.ParseGetParam()
	if nil != err {
		logs.Error("parse get param failed, msg: %s", err.Error())
		this.ErrorMessage(code, err.Error())
		return
	}

	data, code, err := ctx.GetUserList(q)
	if nil != err {
		logs.Error("query user list failed, msg: %s", err.Error())
		this.ErrorMessage(code, "数据库异常")
		return
	}

	rsp := &UserGetRsp{
		Data: UserGetData{
			Pages:    data.Pages,
			PageSize: data.PageSize,
			Total:    data.Total,
			Len:      data.Len,
			List:     make([]UserGetItem, 0),
		},
		Code:    lib.OK,
		Message: "OK",
	}

	for _, v := range data.List {
		item := UserGetItem{
			Id:      v.Id,
			StaffNo: v.StaffNo,
			Name:    v.Name,
			Type:    v.Type,
		}
		rsp.Data.List = append(rsp.Data.List, item)
	}

	this.Data["json"] = rsp
	this.ServeJSON()
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

	_user := &lib.User{
		StaffNo: req.StaffNo,
	}

	err = ctx.Mysql.O.Read(_user, "staff_no")
	if nil != err && orm.ErrNoRows != err {
		logs.Error("query db exception, msg: %s", err.Error())
		this.ErrorMessage(lib.ERR_SYS_MYSQL, "数据库异常")
		return
	} else if nil == err {
		logs.Error("staff has exist, staff_no: %s", req.StaffNo)
		this.ErrorMessage(lib.ERR_PARAM_INVALID, "用户已存在")
		return
	}

	pwd, err := models.CreatePassword(req.StaffNo)
	if nil != err {
		logs.Error("create password error, msg: ", err.Error())
		this.ErrorMessage(lib.ERR_INTERNAL_SERVER_ERROR, "生成密码失败")
		return
	}

	user := &lib.User{
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
			Id: id,
		},
		Code:    lib.OK,
		Message: "OK",
	}

	this.Data["json"] = rsp
	this.ServeJSON()
}

type UserPutRsp struct {
	Code    int    `json:"code" description:"详细错误码"`
	Message string `json:"message" description:"错误描述"`
}

// @Title 用户修改接口(仅管理员可调用)
// @Summary 用户修改接口(仅管理员可调用)
// @Description 用户修改接口(仅管理员可调用)
// @Param id path number true "ID"
// @Param body body controllers.UserReq true "请求参数"
// @Success 201 {object} controllers.UserPutRsp
// @Failure 400
// @router /:id([0-9]+) [put]
func (this *UserController) Put() {
	ctx := GetApiCntx()

	// 解析ID
	id, err := strconv.ParseInt(this.Ctx.Input.Param(":id"), 10, 64)
	if nil != err {
		logs.Error("id parser failed, msg:", err.Error())
		this.BadRequest(lib.ERR_PARAM_INVALID, "ID读取失败")
		return
	}

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

	user := &lib.User{
		Id:         id,
		Name:       req.Name,
		StaffNo:    req.StaffNo,
		Type:       req.Type,
		UpdateId:   admin.Id,
		UpdateTime: time.Now(),
	}

	num, err := ctx.Mysql.O.Update(user, "name", "staff_no", "type", "update_id", "update_time")
	if nil != err {
		logs.Error("update user db exception, msg: ", err.Error())
		this.ErrorMessage(lib.ERR_SYS_MYSQL, "数据库异常")
		return
	} else if num == 0 {
		this.ErrorMessage(lib.ERR_PARAM_INVALID, "用户不存在")
		return
	}

	rsp := &UserPutRsp{
		Code:    lib.OK,
		Message: "OK",
	}

	this.Data["json"] = rsp
	this.ServeJSON()
}

type UserDeleteRsp struct {
	Code    int    `json:"code" description:"详细错误码"`
	Message string `json:"message" description:"错误描述"`
}

// @Title 用户删除接口(仅管理员可调用)
// @Summary 用户删除接口(仅管理员可调用)
// @Description 用户删除接口(仅管理员可调用)
// @Param id path number true "ID"
// @Success 201 {object} controllers.UserPutRsp
// @Failure 400
// @router /:id([0-9]+) [delete]
func (this *UserController) Delete() {
	ctx := GetApiCntx()

	// 解析ID
	id, err := strconv.ParseInt(this.Ctx.Input.Param(":id"), 10, 64)
	if nil != err {
		logs.Error("id parser failed, msg:", err.Error())
		this.BadRequest(lib.ERR_PARAM_INVALID, "ID读取失败")
		return
	}

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

	// 删除ID
	user := &lib.User{
		Id: id,
	}

	num, err := ctx.Mysql.O.Delete(user)
	if nil != err {
		logs.Error("update user db exception, msg: ", err.Error())
		this.ErrorMessage(lib.ERR_SYS_MYSQL, "数据库异常")
		return
	} else if num == 0 {
		this.ErrorMessage(lib.ERR_PARAM_INVALID, "用户不存在")
		return
	}

	rsp := &UserDeleteRsp{
		Code:    lib.OK,
		Message: "OK",
	}

	this.Data["json"] = rsp
	this.ServeJSON()
}

type UserPasswordReq struct {
	OldPwd string `json:"old_pwd" description:"旧密码"`
	NewPwd string `json:"new_pwd" description:"新密码"`
}

type UserPasswordRsp struct {
	Code    int    `json:"code" description:"详细错误码"`
	Message string `json:"message" description:"错误描述"`
}

// @Title 用户修改密码接口
// @Summary 用户修改密码接口
// @Description 用户修改密码接口
// @Param body body controllers.UserPasswordReq true "请求参数"
// @Success 201 {object} controllers.UserPasswordRsp
// @Failure 400
// @router /password [post]
func (this *UserController) Password() {
	ctx := GetApiCntx()

	admin, code, err := this.getAdmin()
	if nil != err {
		logs.Error("get admin failed, msg: %s", err.Error())
		this.ErrorMessage(code, err.Error())
		return
	}

	// 获取参数
	req := &UserPasswordReq{}

	err = jsoniter.Unmarshal(this.Ctx.Input.RequestBody, req)
	if nil != err {
		logs.Error("Parameter format is invalid! body: ",
			this.Ctx.Input.RequestBody, " msg: ", err.Error())
		this.ErrorMessage(lib.ERR_PARAM_INVALID, "参数解析失败")
		return
	}

	req.OldPwd = lib.RemoveEmpty(req.OldPwd)
	if 0 == len(req.OldPwd) {
		logs.Error("OldPwd not exist")
		this.ErrorMessage(lib.ERR_PARAM_INVALID, "旧密码不能为空")
		return
	}
	req.NewPwd = lib.RemoveEmpty(req.NewPwd)
	if 0 == len(req.NewPwd) {
		logs.Error("NewPwd not exist")
		this.ErrorMessage(lib.ERR_PARAM_INVALID, "新密码不能为空")
		return
	}

	// 验证密码
	_, err = models.MatchPassword(req.OldPwd, admin.Password)
	if nil != err {
		logs.Error("password not match, msg: %s", err.Error())
		this.ErrorMessage(comm.ERR_PARAM_INVALID, "密码错误")
		return
	}

	pwd, err := models.CreatePassword(req.NewPwd)
	if nil != err {
		logs.Error("create password error, msg: ", err.Error())
		this.ErrorMessage(lib.ERR_INTERNAL_SERVER_ERROR, "生成密码失败")
		return
	}

	user := &lib.User{
		Id:         admin.Id,
		Password:   pwd,
		UpdateId:   admin.Id,
		UpdateTime: time.Now(),
	}

	num, err := ctx.Mysql.O.Update(user, "password", "update_id", "update_time")
	if nil != err {
		logs.Error("update user db exception, msg: ", err.Error())
		this.ErrorMessage(lib.ERR_SYS_MYSQL, "数据库异常")
		return
	} else if num == 0 {
		this.ErrorMessage(lib.ERR_PARAM_INVALID, "用户不存在")
		return
	}

	rsp := &UserPasswordRsp{
		Code:    lib.OK,
		Message: "OK",
	}

	this.Data["json"] = rsp
	this.ServeJSON()
}

////////////////////////////////////////////////////////////////////////////////
// 通用请求
type UserReq struct {
	Name    string `json:"name" description:"用户昵称" required:"true"`
	StaffNo string `json:"staff_no" description:"用户工号" required:"true"`
	Type    int    `json:"type" description:"用户类别" required:"true"`
}

// 获取管理员接口传入参数
func (this *UserController) getUserReqParam() (*UserReq, int, error) {

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

	return req, lib.OK, nil
}

// 解析Get参数
func (o *UserController) ParseGetParam() (map[string]interface{}, int, error) {
	q := make(map[string]interface{})

	// 用户名称
	name := o.GetString("name")
	if "" != name {
		q["name"] = name
	}

	// 用户工号
	staffNo := o.GetString("staff_no")
	if "" != name {
		q["staff_no"] = staffNo
	}

	// 用户类型
	_type := o.GetString("type")
	if "" != _type {
		typ, err := strconv.Atoi(_type)
		if nil != err {
			return nil, lib.ERR_PARAM_INVALID, errors.New("用户类型错误")
		}
		q["type"] = typ
	}

	// 页号
	_page := o.GetString("page")
	if "" != _page {
		page, err := strconv.Atoi(_page)
		if nil != err {
			return nil, lib.ERR_PARAM_INVALID, errors.New("页号错误")
		}
		q["page"] = page
	} else {
		q["page"] = lib.DEFAULT_PAGE // 默认第1页
	}

	// 每页最大条目数
	_pageSize := o.GetString("page_size")
	if "" != _pageSize {
		pageSize, err := strconv.Atoi(_pageSize)
		if nil != err {
			return nil, lib.ERR_PARAM_INVALID, errors.New("页面大小错误")
		}
		q["page_size"] = int(pageSize)
	} else {
		q["page_size"] = lib.DEFAULT_PAGE_SIZE // 默认10条
	}

	return q, lib.OK, nil
}
