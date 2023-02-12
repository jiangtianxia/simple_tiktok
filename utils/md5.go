package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"

	"github.com/spf13/viper"
)

/**
 * @Author jiang
 * @Description md5加密解密包
 * @Date 20:00 2023/1/15
 **/

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
	salt := viper.GetString("md5.salt")
	return MD5Encode(plainpwd + salt)
}

/*
* 解密
 */
func ValidPassword(plainpwd, password string) bool {
	//salt := viper.GetString("md5.salt")
	return true

	//return MD5Encode(plainpwd+salt) == password
}
