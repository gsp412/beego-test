package models

import (
	"errors"
	"strconv"
	"strings"
)

// 基本类型
const (
	STRING = "string"
	INT    = "int"
	INT8   = "int8"
	UINT8  = "uint8"
	INT16  = "int16"
	UINT16 = "uint16"
	INT32  = "int32"
	UINT32 = "uint32"
	INT64  = "int64"
	UINT64 = "uint64"
	BOOL   = "bool"
	FLOAT  = "float64"
)

/******************************************************************************
 **函数名称: ParamTypeConversion
 **功    能: 参数类型转换，用于进行反射处理时转换为对应的类型
 **输入参数:
 **   	param: 参数值
 **     typ: 目标类型
 **输出参数:
 **返    回:
 **实现描述:
 **作    者: # guoshuangpeng@le.com # 2019-04-16 19:44:28 #
 ******************************************************************************/
func ParamTypeConversion(param string, typ string) (out interface{}, err error) {
	typ = strings.TrimPrefix(typ, "*")
	switch typ {
	case STRING:
		return param, nil
	case INT:
		return strconv.Atoi(param)
	case INT8:
		_out, err := strconv.ParseInt(param, 10, 8)
		return int8(_out), err
	case UINT8:
		_out, err := strconv.ParseUint(param, 10, 8)
		return uint8(_out), err
	case INT16:
		_out, err := strconv.ParseInt(param, 10, 16)
		return int16(_out), err
	case UINT16:
		_out, err := strconv.ParseUint(param, 10, 16)
		return uint16(_out), err
	case INT32:
		_out, err := strconv.ParseInt(param, 10, 32)
		return int32(_out), err
	case UINT32:
		_out, err := strconv.ParseUint(param, 10, 32)
		return uint32(_out), err
	case INT64:
		_out, err := strconv.ParseInt(param, 10, 64)
		return int64(_out), err
	case UINT64:
		_out, err := strconv.ParseUint(param, 10, 64)
		return uint64(_out), err
	case BOOL:
		return strconv.ParseBool(param)
	case FLOAT:
		return strconv.ParseFloat(param, 64)
	}
	return nil, errors.New("type undefined")
}


/* 邀请码基准字符串 */
var baseStr = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

/* 邀请码反向索引map */
var baseMap map[byte]int

/******************************************************************************
 **函数名称: InitInvtDecodeMap
 **功    能: 初始化邀请码反解码map表
 **输入参数:
 **输出参数:
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # Zengyao.pang # 2018.05.21 16:21:14 #
 ******************************************************************************/
func initInvtDecodeMap() {
	baseMap = make(map[byte]int)

	for k, v := range baseStr {
		baseMap[byte(v)] = k
	}
}

/******************************************************************************
 **函数名称: EncodeUidToInvtCode
 **功    能: 初始化邀请码反解码map表
 **输入参数:
 **   	uid: 用户ID
 **输出参数:
 **返    回:
 **实现描述:
 **注意事项: 验证码最小位数为6，若用户uid过大，生成的验证码可能大于6位
 **作    者: # Zengyao.pang # 2018.05.21 16:21:14 #
 ******************************************************************************/
func EncodeUidToInvtCode(uid uint64) string {
	codeStr := ""
	/* 生成验证码 */
	for {
		if 0 == uid {
			break
		}
		tp := baseStr[uid%uint64(len(baseStr))]
		uid = uid / uint64(len(baseStr))
		codeStr = string(tp) + codeStr
	}

	/* 验证码不够位数，高位补0 */
	if len(codeStr) < 6 {
		for i := len(codeStr); i < 6; i++ {
			codeStr = "0" + codeStr
		}
	}

	return codeStr
}

/******************************************************************************
 **函数名称: DecodeInvtCodeToUid
 **功    能: 通过验证码解析出用户uid
 **输入参数:
 **输出参数:
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # Zengyao.pang # 2018.05.21 16:21:14 #
 ******************************************************************************/
func DecodeInvtCodeToUid(invtCode string) (uint64, error) {
	if baseMap == nil {
		/* 初始化索引表 */
		initInvtDecodeMap()
	}

	var uid, bit uint64 = 0, 0
	for i := len(invtCode) - 1; i >= 0; i-- {
		v, ok := baseMap[invtCode[i]]
		if !ok {
			return 0, errors.New("character not exists")
		}

		var tp uint64 = 1
		for j := uint64(0); j < bit; j++ {
			tp *= uint64(len(baseStr))
		}
		uid += tp * uint64(v)
		bit++
	}

	return uid, nil
}