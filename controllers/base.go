package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"beego-test/lib"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/json-iterator/go"
)

type BaseController struct {
	beego.Controller
}

/* 异常消息回复 */
func (b *BaseController) SendJson(status int, m map[string]interface{}) {
	str, _ := jsoniter.Marshal(m)

	/* 应答 */
	b.Ctx.ResponseWriter.WriteHeader(status)
	b.Ctx.ResponseWriter.Write([]byte(str))

	b.StopRun()
}

/* 回复应答消息 */
func (b *BaseController) SendResponse(status int, code int, message string) {
	m := make(map[string]interface{})

	m["code"] = code
	m["message"] = message

	str, _ := jsoniter.Marshal(m)

	/* 应答 */
	b.Ctx.ResponseWriter.WriteHeader(status)
	b.Ctx.ResponseWriter.Write([]byte(str))

	b.StopRun()
}

/* 发送应答消息 */
func (b *BaseController) ErrorMessage(code int, message string) {
	if code >= lib.ERR_BAD_REQ && code < lib.ERR_AUTH {
		b.BadRequest(code, message)
		return
	} else if code >= lib.ERR_AUTH && code < lib.ERR_FORBIDDEN {
		b.Unauthorized(code, message)
		return
	} else if code >= lib.ERR_FORBIDDEN && code < lib.ERR_NOT_FOUND {
		b.Forbidden(code, message)
		return
	} else if code >= lib.ERR_NOT_FOUND && code < lib.ERR_METHOD_NOT_ALLOWED {
		b.NotFound(code, message)
		return
	} else if code >= lib.ERR_METHOD_NOT_ALLOWED && code < lib.ERR_GONE {
		b.MethodNotAllowed(code, message)
		return
	} else if code >= lib.ERR_GONE && code < lib.ERR_UNSUPPORTED_MEDIA_TYPE {
		b.Gone(code, message)
		return
	} else if code >= lib.ERR_UNSUPPORTED_MEDIA_TYPE && code < lib.ERR_UNPROCESSABLE_ENTITY {
		b.UnsupportedMediaType(code, message)
		return
	} else if code >= lib.ERR_UNPROCESSABLE_ENTITY && code < lib.ERR_TOO_MANY_REQ {
		b.UnprocessableEntity(code, message)
		return
	} else if code >= lib.ERR_TOO_MANY_REQ && code < lib.ERR_INTERNAL_SERVER_ERROR {
		b.TooManyRequests(code, message)
		return
	} else if code >= lib.ERR_INTERNAL_SERVER_ERROR && code < lib.ERR_SVC_UNAVAILABLE {
		b.InternalServerError(code, message)
		return
	}
	b.SendResponse(http.StatusOK, code, message)
	return
}

/* 发送应答消息 */
func (b *BaseController) ErrorMessageOld(code int, message string) {
	b.SendResponse(http.StatusOK, code, message)
	return
}

//1.服务端错误
// 1.1 Internal server error
//  @状态码: 500
//  @状态含义: Internal server error
//  @状态原因: 客户端请求有效, 服务器处理时发生了意外!
//  @错误码
func (b *BaseController) InternalServerError(code int, message string) {
	b.SendResponse(http.StatusInternalServerError, code, message)
}

//2.客户端错误
// 2.1 Bad request
//  @状态码: 400
//  @状态含义: Bad request
//  @状态原因: 服务器不理解客户端的请求, 未做任何处理!
//  @错误码
func (b *BaseController) BadRequest(code int, message string) {
	b.SendResponse(http.StatusBadRequest, code, message)
}

// 2.2 Unauthorized
//  @状态码: 401
//  @状态含义: Unauthorized
//  @状态原因: 用户未提供身份验证凭据, 或者没有通过身份验证!
//  @错误码
func (b *BaseController) Unauthorized(code int, message string) {
	b.SendResponse(http.StatusUnauthorized, code, message)
}

// 2.3 Forbidden
//  @状态码: 403
//  @状态含义: Forbidden
//  @状态原因: 用户通过了身份验证, 但是不具有访问资源所需的权限!
//  @错误码
func (b *BaseController) Forbidden(code int, message string) {
	b.SendResponse(http.StatusForbidden, code, message)
}

// 2.4 Not found
//  @状态码: 404
//  @状态含义: Not found
//  @状态原因: 所请求的资源不存在, 或不可用!
//  @错误码
func (b *BaseController) NotFound(code int, message string) {
	b.SendResponse(http.StatusNotFound, code, message)
}

// 2.5 Method not allowed
//  @状态码: 405
//  @状态含义: Method not allowed
//  @状态原因: 用户已经通过身份验证, 但是所用的HTTP方法不在他的权限之内!
//  @错误码
func (b *BaseController) MethodNotAllowed(code int, message string) {
	b.SendResponse(http.StatusMethodNotAllowed, code, message)
}

// 2.6 Gone
//  @状态码: 410
//  @状态含义: Gone
//  @状态原因: 所请求的资源已从这个地址转移, 不再可用!
//  @错误码
func (b *BaseController) Gone(code int, message string) {
	b.SendResponse(http.StatusGone, code, message)
}

// 2.7 Unsupported media type
//  @状态码: 415
//  @状态含义: Unsupported media type
//  @状态原因: 客户端要求的返回格式不支持. 比如: API只能返回JSON格式, 但是客户端要求返回XML格式!
//  @错误码
func (b *BaseController) UnsupportedMediaType(code int, message string) {
	b.SendResponse(http.StatusUnsupportedMediaType, code, message)
}

// 2.8 Unprocessable entity
//  @状态码: 422
//  @状态含义: Unprocessable entity
//  @状态原因: 客户端上传的附件无法处理, 导致请求失败!
//  @错误码
func (b *BaseController) UnprocessableEntity(code int, message string) {
	b.SendResponse(http.StatusUnprocessableEntity, code, message)
}

// 2.9 Too many requests
//  @状态码: 429
//  @状态含义: Too many requests
//  @状态原因: 客户端的请求次数超过限额!
//  @错误码
func (b *BaseController) TooManyRequests(code int, message string) {
	b.SendResponse(http.StatusTooManyRequests, code, message)
}

////////////////////////////////////////////////////////////////////////////////
// 解析Get参数

func (b *BaseController) GetQueryParam(v interface{}) (code int, err error) {
	defer func() {
		// 将panic转换为error 返回
		if r := recover(); r != nil {
			code = lib.ERR_PARAM_INVALID
			err = errors.New("param unmarshal failed")
		}
	}()

	o := reflect.ValueOf(v)

	// 判断如果传入类型不是指针或者是空值或者值不可改则返回序列化失败
	if o.Kind() != reflect.Ptr || o.IsNil() || !o.Elem().CanSet() {
		return lib.ERR_PARAM_INVALID, errors.New("param unmarshal")
	}

	vt := reflect.TypeOf(v).Elem()
	vr := o.Elem()

	for i := 0; i < vt.NumField(); i++ {

		// 如果结构体字段为指针类型 但尚未初始化则返回失败
		if vr.Field(i).Kind() == reflect.Ptr && !vr.Field(i).Elem().CanSet() {
			errMsg := fmt.Sprintf("[%s] uninitialized", vt.Field(i).Name)
			return lib.ERR_PARAM_INVALID, errors.New(errMsg)
		}

		// 未设置注解的略过
		if "" == vt.Field(i).Tag.Get("json") {
			continue
		}

		// 1、读取参数
		_value := b.GetString(vt.Field(i).Tag.Get("json"), vt.Field(i).Tag.Get("def"))
		if "" != vt.Field(i).Tag.Get("request") && "" == _value {
			errMsg := fmt.Sprintf("[%s] param miss", vt.Field(i).Tag.Get("json"))
			return lib.ERR_PARAM_MISS, errors.New(errMsg)
		} else if "" == _value {
			// 如果未取到值则制空
			vr.Field(i).Set(reflect.Zero(vr.Field(i).Type()))
			continue
		}

		// 2、读取类型
		typ := vr.Field(i).Type().String()

		// 3、取实际值
		value, err := lib.ParamTypeConversion(_value, typ)
		if nil != err {
			return lib.ERR_PARAM_INVALID, err
		}

		// 4、传值
		if vr.Field(i).Kind() == reflect.Ptr {
			vr.Field(i).Elem().Set(reflect.ValueOf(value))
		} else {
			vr.Field(i).Set(reflect.ValueOf(value))
		}

	}

	return 0, nil
}

// 获取管理员信息
func (this *BaseController) getAdmin() (*lib.User, int, error) {
	_user := this.GetSession("user")
	if nil == _user {
		logs.Error("get user failed")
		return nil, lib.ERR_FORBIDDEN, errors.New("获取操作人信息失败")
	}
	user, ok := _user.(lib.User)
	if !ok {
		logs.Error("user type not match")
		return nil, lib.ERR_FORBIDDEN, errors.New("获取操作人信息失败")
	}
	return &user, lib.OK, nil
}