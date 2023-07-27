package middlewares

import (
	"net/http"
	"simple_tiktok/model"
	"simple_tiktok/store"
	"simple_tiktok/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// TokenValidator 是一个 Gin 中间件函数，用于验证 HTTP 请求中的 token 是否有效
func TokenValidator() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取 Authorization 请求头
		token := ctx.GetHeader("Authorization")

		userInfo := utils.UserClaims{}
		if token != "" {
			// 验证token
			userClaim, err := utils.AnalyzeToken(token)
			if err != nil {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				ctx.Abort()
				return
			}

			if userClaim == nil || userClaim.Id <= 0 || userClaim.Issuer != utils.GetIssuer() || userClaim.Username == "" {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				ctx.Abort()
				return
			}

			// 判断token是否过期
			if time.Now().Unix() > userClaim.ExpiresAt {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
				ctx.Abort()
				return
			}

			// 判断是否存在该用户
			var cnt int64
			if err := store.GetDB().Model(&model.UserBasic{}).
				Where("id = ? AND username = ?", userClaim.Id, userClaim.Username).Count(&cnt).Error; err != nil || cnt == 0 {
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				ctx.Abort()
				return
			}

			userInfo.Id = userClaim.Id
			userInfo.Username = userClaim.Username
		}

		ctx.Set("tokenId", userInfo.Id)
		ctx.Set("tokenUsername", userInfo.Username)
		ctx.Next()
	}
}
