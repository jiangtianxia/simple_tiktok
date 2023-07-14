package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

var salt string

func SetMd5Salt(s string) {
	salt = s
}

/*
* 小写
 */
func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	tempStr := h.Sum(nil)

	return hex.EncodeToString(tempStr)
}

/*
* 大写
 */
func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode(data))
}

/*
* 加密
 */
func MakePassword(plainpwd string) string {
	return MD5Encode(plainpwd + salt)
}

/*
* 解密
 */
func ValidPassword(plainpwd, password string) bool {
	return MD5Encode(plainpwd+salt) == password
}
