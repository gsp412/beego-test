package models

import (
	"beego-test/lib"
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
)

// 验证Token
func (ctx ApiCntx) VerifyToken(tk string) (*lib.User, int, error) {

	strs := strings.Split(tk, ".")
	if len(strs) < 2 {
		return nil, lib.ERR_AUTH, errors.New("tk参数错误")
	}

	id, err := lib.DecodeInvtCodeToUid(strs[0])
	if nil != err {
		return nil, lib.ERR_AUTH, errors.New("tk参数错误")
	}

	user := &lib.User{
		Id: id,
	}

	if err := ctx.Mysql.O.Read(user); nil != err && orm.ErrNoRows != err {
		return nil, lib.ERR_SYS_MYSQL, errors.New("数据库异常")
	} else if orm.ErrNoRows == err {
		return nil, lib.ERR_AUTH, errors.New("用户不存在")
	}

	checkPwd := lib.Md5Sum(user.Password)

	orig, err := lib.DesCbcDecrypt(strs[2], []byte(user.Password))
	if nil != err {
		return nil, lib.ERR_AUTH, errors.New("tk验证失败")
	}

	origs := strings.Split(orig, ".")
	if len(origs) < 2 {
		return nil, lib.ERR_AUTH, errors.New("tk参数错误")
	}
	tm, err := strconv.ParseInt(origs[0], 10, 64)
	if nil != err {
		return nil, lib.ERR_AUTH, errors.New("tk验证失败")
	}

	if (time.Now().Unix() - tm) > lib.TK_TIME_OUT {
		return nil, lib.ERR_AUTH, errors.New("tk过期")
	}

	if checkPwd != origs[1] {
		return nil, lib.ERR_AUTH, errors.New("tk过期")
	}

	return user, lib.OK, nil
}


// 创建Token
func (ctx ApiCntx) CreateToken(uid int64, pwd string) (string, int, error) {

	forPwd := lib.Md5Sum(pwd)

	orig := fmt.Sprintf("%d.%s", time.Now().Unix(), forPwd)

	crypted, err := lib.DesCbcEncrypt([]byte(orig), []byte(pwd))
	if nil != err {
		return "", lib.ERR_INTERNAL_SERVER_ERROR, err
	}

	invite := lib.EncodeUidToInvtCode(uid)

	return fmt.Sprintf("%s.%s", invite, crypted), lib.OK, nil
}