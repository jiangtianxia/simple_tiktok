package utils

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

/**
 * @Author jiang
 * @Description jwt生成token包
 * @Date 20:00 2023/1/15
 **/

var myKey = []byte(viper.GetString("jwt.key"))

type UserClaims struct {
	Identity string `json:"identity"`
	Username string `json:"username"`
	Usericon string `json:"usericon"`
	jwt.StandardClaims
}

// 生成token
func GenerateToken(identity, username, usericon string) (string, error) {
	userClaim := &UserClaims{
		Identity: identity,
		Username: username,
		Usericon: usericon,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(
				time.Duration(viper.GetInt("jwt.expire")) * time.Hour).Unix(), // 过期时间
			Issuer: "simple_tiktok", // 签发人
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// 解析token
func AnalyseToken(tokenString string) (*UserClaims, error) {
	UserClaims := new(UserClaims)

	claims, err := jwt.ParseWithClaims(tokenString, UserClaims, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !claims.Valid {
		return nil, fmt.Errorf("analyse Token Error:%v", err)
	}
	return UserClaims, nil
}
