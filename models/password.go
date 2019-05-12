package models

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
)

const SECRET_KEY = "HGyQn2UN5cRXcV9dOPYKqZNpSXxNxWNz"

// 生成密码
func CreatePassword(pwd string) (string, error) {

	md5Pwd := Md5Sum(pwd)

	orig := fmt.Sprintf("%d.%s", time.Now().Unix(), md5Pwd)

	return DesCbcEncrypt([]byte(orig), []byte(SECRET_KEY))
}

// 验证密码
func MatchPassword(pwd string, encodePwd string) (bool, error) {

	md5Pwd := Md5Sum(pwd)

	orig, err := DesCbcDecrypt(encodePwd, []byte(SECRET_KEY))
	if nil != err  {
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


///////////////////////////////////////////////////////////////////////////////
// 加密算法

/******************************************************************************
 **函数名称: Md5Sum
 **功    能: MD5加密处理
 **输入参数:
 **     s: 被加密处理的字串
 **输出参数: NONE
 **返    回: 加密字串
 **实现描述:
 **注意事项:
 **作    者: # Qifeng.zou # 2018.04.26 15:35:08 #
 ******************************************************************************/
func Md5Sum(str string) string {
	h := md5.New()

	h.Write([]byte(str))

	return string(hex.EncodeToString(h.Sum(nil)))
}

/******************************************************************************
 **函数名称: DesEncrypt
 **功    能: Des加密
 **输入参数:
 **     orig: 被加密处理的字串
 **     key: 加密秘钥
 **输出参数: NONE
 **返    回: 加密字串
 **实现描述:
 **注意事项:
 **作    者: # Qifeng.zou # 2018.09.28 16:53:28 #
 ******************************************************************************/
func DesCbcEncrypt(orig, key []byte) (string, error) {

	block, err := des.NewCipher(key)
	if nil != err {
		return "", err
	}

	blockSize := block.BlockSize()
	orig = PKCS5Padding(orig, blockSize)
	//获取CBC加密模式
	iv := key[:blockSize] //用密钥作为向量(不建议这样使用)
	mode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(orig))
	mode.CryptBlocks(crypted, orig)
	return  fmt.Sprintf("%X", crypted), nil
}

/******************************************************************************
 **函数名称: DesCbcDecrypt
 **功    能: DES解密
 **输入参数:
 **     crypted: 被加密处理的字串
 **     key: 加密秘钥
 **输出参数: NONE
 **返    回: 加密字串
 **实现描述:
 **注意事项:
 **作    者: # Qifeng.zou # 2018.09.28 16:55:15 #
 ******************************************************************************/
func DesCbcDecrypt(crypted string, key []byte) (string, error) {

	defer func() {
		recover()
	}()

	data, err := hex.DecodeString(crypted)
	if err != nil {
		return "", err
	}

	block, err := des.NewCipher(key)
	if nil != err {
		return "", err
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	orig := make([]byte, len(data))
	blockMode.CryptBlocks(orig, data)
	orig = PKCS5UnPadding(orig)

	return string(orig), nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(orig []byte) []byte {
	length := len(orig)
	unpadding := int(orig[length-1])
	return orig[:(length - unpadding)]
}
