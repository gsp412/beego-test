package lib

// 业务常量配置
const (
	TK_TIME_OUT = 7200 // tk超时时间，单位：秒
)

// 正常
//  @状态码: 200
//  @状态含义: Normal
//  @状态原因: 无异常
//  @错误码
const (
	OK = 0 // 正常
)

// 1.服务端错误码
// 1.1 Internal server error
//  @状态码: 500
//  @状态含义: Internal server error
//  @状态原因: 客户端请求有效, 服务器处理时发生了意外!
//  @错误码
const (
	ERR_INTERNAL_SERVER_ERROR = 500001 // 服务器内部错误
	ERR_SYS_TOO_BUSY          = 500002 // 系统繁忙
	ERR_SYS_RPC               = 500003 // RPC异常
	ERR_SYS_DB                = 500004 // 数据库异常
	ERR_SYS_MYSQL             = 500005 // MYSQL异常
	ERR_SYS_REDIS             = 500006 // REDIS异常
	ERR_SYS_MONGO             = 500007 // MONGO异常
)

// 1.2 Service unavailable
//  @状态码: 503
//  @状态含义: Service unavailable
//  @状态原因: 服务器无法处理请求, 一般用于网站维护状态!
//  @错误码
const (
	ERR_SVC_UNAVAILABLE = 503001 // 服务不可达
)

// 2.客户端端错误码
// 2.1 Bad request
//  @状态码: 400
//  @状态含义: Bad request
//  @状态原因: 服务器不理解客户端的请求, 未做任何处理!
//  @错误码
const (
	ERR_BAD_REQ       = 400001 // 非法请求
	ERR_PARAM_MISS    = 400002 // 参数缺失
	ERR_PARAM_INVALID = 400003 // 参数非法
	ERR_HEAD_INVALID  = 400004 // 报头非法
	ERR_BODY_INVALID  = 400005 // 报体非法
)

// 2.2 Unauthorized
//  @状态码: 401
//  @状态含义: Unauthorized
//  @状态原因: 用户未提供身份验证凭据, 或者没有通过身份验证!
//  @错误码
const (
	ERR_AUTH = 401001 // 鉴权失败
)

// 2.3 Forbidden
//  @状态码: 403
//  @状态含义: Forbidden
//  @状态原因: 用户通过了身份验证, 但是不具有访问资源所需的权限!
//  @错误码
const (
	ERR_FORBIDDEN = 403001 // 访问受限
)

// 2.4 Not found
//  @状态码: 404
//  @状态含义: Not found
//  @状态原因: 所请求的资源不存在, 或不可用!
//  @错误码
const (
	ERR_NOT_FOUND = 404001 // 资源不存在
)

// 2.5 Method not allowed
//  @状态码: 405
//  @状态含义: Method not allowed
//  @状态原因: 用户已经通过身份验证, 但是所用的HTTP方法不在他的权限之内!
//  @错误码
const (
	ERR_METHOD_NOT_ALLOWED = 405001 // 无本HTTP访问权限
)

// 2.6 Gone
//  @状态码: 410
//  @状态含义: Gone
//  @状态原因: 所请求的资源已从这个地址转移, 不再可用!
//  @错误码
const (
	ERR_GONE = 410001 // 资源不再可用
)

// 2.7 Unsupported media type
//  @状态码: 415
//  @状态含义: Unsupported media type
//  @状态原因: 客户端要求的返回格式不支持. 比如: API只能返回JSON格式, 但是客户端要求返回XML格式!
//  @错误码
const (
	ERR_UNSUPPORTED_MEDIA_TYPE = 415001 // 不支持的返回格式
)

// 2.8 Unprocessable entity
//  @状态码: 422
//  @状态含义: Unprocessable entity
//  @状态原因: 客户端上传的附件无法处理, 导致请求失败!
//  @错误码
const (
	ERR_UNPROCESSABLE_ENTITY = 422001 // 不支持的返回格式
)

// 2.9 Too many requests
//  @状态码: 429
//  @状态含义: Too many requests
//  @状态原因: 客户端的请求次数超过限额!
//  @错误码
const (
	ERR_TOO_MANY_REQ = 429001 // 请求次数超过限制
)
