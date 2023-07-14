package utils

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	myKey     []byte
	jwtExpire int
)

func InitJwt(s string, expire int) {
	myKey = []byte(s)
	jwtExpire = expire
}

type UserClaims struct {
	Identity uint64 `json:"identity"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// 生成token
func GenerateToken(identity uint64, username string) (string, error) {
	userClaim := &UserClaims{
		Identity: identity,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(
				time.Duration(jwtExpire) * time.Hour).Unix(), // 过期时间
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
