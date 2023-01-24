package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simple_tiktok/logger"
	"simple_tiktok/utils"
)

func UserInfo(c *gin.Context) {
	token, err := utils.GenerateToken("1", "zxy", "")
	if err != nil {
		logger.SugarLogger.Error("Generate Token Error:" + err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"token":   token,
	})
}
