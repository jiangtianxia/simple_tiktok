package utils

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	myKey     []byte
	jwtExpire int
	issuer    string
)

func InitJwt(s string, expire int) {
	myKey = []byte(s)
	jwtExpire = expire
	issuer = "simple_tiktok"
}

func GetIssuer() string {
	return issuer
}

type UserClaims struct {
	Id       uint   `json:"identity"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// 生成token
func GenerateToken(id uint, username string) (string, error) {
	userClaim := &UserClaims{
		Id:       id,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(
				time.Duration(jwtExpire) * time.Hour).Unix(), // 过期时间
			Issuer: issuer, // 签发人
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
func AnalyzeToken(tokenString string) (*UserClaims, error) {
	UserClaims := new(UserClaims)

	claims, err := jwt.ParseWithClaims(tokenString, UserClaims, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !claims.Valid {
		return nil, fmt.Errorf("analyze Token Error:%v", err)
	}
	return UserClaims, nil
}
