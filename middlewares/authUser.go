package middlewares

import (
	"errors"
	"simple_tiktok/dao/mysql"
	"simple_tiktok/logger"
	"simple_tiktok/utils"
	"time"
)

// 验证用户token信息
func AuthUserCheck(token string) (*utils.UserClaims, error) {
	userClaim, err := utils.AnalyseToken(token)
	if err != nil {
		logger.SugarLogger.Error(err)
		return userClaim, err
	}

	if userClaim == nil || userClaim.Identity == 0 || userClaim.Issuer != "simple_tiktok" || userClaim.Username == "" {
		logger.SugarLogger.Error("Unauthorized User")
		return userClaim, errors.New("unauthorized user")
	}

	// 判断token是否过期
	if time.Now().Unix() > userClaim.ExpiresAt {
		logger.SugarLogger.Error("Token Expired")
		return userClaim, errors.New("token expired")
	}

	// 判断是否存在该用户
	_, err = mysql.FindUserByIdentity(userClaim.Identity)
	if err != nil {
		logger.SugarLogger.Error("User Invalid", userClaim)
		return userClaim, err
	}

	return userClaim, nil
}
