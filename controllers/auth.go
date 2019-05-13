package controllers

import (
	"net/http"
	"regexp"
	"strings"

	"beego-test/lib"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
)

// 白名单
const WHITE_LIST = "/beego-test/v1/token"

func Auth(c *context.Context) {
	ctx := GetApiCntx()

	// 判断如果是白名单url，不做限制
	if IsWhiteList(c.Request.URL.Path) {
		logs.Info("path match white list, path: %s", c.Request.URL.Path)
		return
	}

	// 从 header 中读取tk
	tk := c.Request.Header.Get("tk")
	if 0 == len(tk) {
		logs.Error("tk not exist")
		SendError(c, lib.ERR_AUTH, "tk不能为空")
		return
	}

	user, code, err := ctx.VerifyToken(tk)
	if nil != err {
		logs.Error("verify user permissions failed, tk: %s, msg: %s", tk, err.Error())
		SendError(c, code, err.Error())
		return
	}

	// 存储session
	err = c.Input.CruSession.Set("user", *user)
	if nil != err {
		logs.Error("session 设置失败")
		SendError(c, lib.ERR_INTERNAL_SERVER_ERROR, "服务器内部错误")
		return
	}

	return
}

// 验证是否是白名单
// path: 待匹配路径
func IsWhiteList(path string) bool {

	list := strings.Split(WHITE_LIST, ",")

	for _, v := range list {
		match, _ := regexp.MatchString(v, path)
		if match {
			return true
		}
	}

	return false
}

////////////////////////////////////////////////////////////////////////////////
// 鉴权模块返回结果

type ErrRsp struct {
	Code    int    `json:"code" description:"详细错误码"`
	Message string `json:"message" description:"错误描述"`
}

func SendError(c *context.Context, code int, msg string) {
	/* 接收panic异常 打印日志 */
	defer func() {
		if r := recover(); nil != r {
			logs.Error(r)
		}
	}()

	status := http.StatusOK
	if code >= lib.ERR_BAD_REQ && code < lib.ERR_AUTH {
		status = http.StatusBadRequest
	} else if code >= lib.ERR_AUTH && code < lib.ERR_FORBIDDEN {
		status = http.StatusUnauthorized
	} else if code >= lib.ERR_FORBIDDEN && code < lib.ERR_NOT_FOUND {
		status = http.StatusForbidden
	} else if code >= lib.ERR_NOT_FOUND && code < lib.ERR_METHOD_NOT_ALLOWED {
		status = http.StatusNotFound
	} else if code >= lib.ERR_INTERNAL_SERVER_ERROR && code < lib.ERR_SVC_UNAVAILABLE {
		status = http.StatusInternalServerError
	}

	/* 返回JSON */
	rsp := &ErrRsp{
		Code:    code,
		Message: msg,
	}

	c.Output.Status = status
	c.Output.JSON(rsp, false, false)

	/* 终止进程 */
	c.Abort(status, msg)
}
