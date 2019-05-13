package models

import (
	"beego-test/lib"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	SECRET_KEY  = "HGyQn2UN5cRXcV9dOPYKqZNpSXxNxWNz" // 加密秘钥
	SECRET_SLAT = "9YFAhvR2sgfWu1xURVnGoQf0xRKTSmxs" // 加密盐
)

// 生成密码
func CreatePassword(pwd string) (string, error) {

	md5Pwd := lib.Md5Sum(pwd + SECRET_SLAT)

	orig := fmt.Sprintf("%d.%s", time.Now().Unix(), md5Pwd)

	return lib.DesCbcEncrypt([]byte(orig), []byte(SECRET_KEY))
}

// 验证密码
func MatchPassword(pwd string, encodePwd string) (bool, error) {

	md5Pwd := lib.Md5Sum(pwd + SECRET_SLAT)

	orig, err := lib.DesCbcDecrypt(encodePwd, []byte(SECRET_KEY))
	if nil != err {
		return false, err
	}

	strs := strings.Split(orig, ".")
	if len(strs) < 2 {
		return false, errors.New("password not match")
	}

	if strs[1] != md5Pwd {
		return false, errors.New("password not match")
	}

	return true, nil
}
