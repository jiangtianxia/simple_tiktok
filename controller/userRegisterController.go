package controller

import (
	"net/http"
	"simple_tiktok/service"

	"github.com/gin-gonic/gin"
)

// 用户注册 /register
func UserRegister(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	req := service.RegisterRequire{
		Username: string(username),
		Password: string(password),
	}
	registerResponse, err := service.PostUserRegister(c, &req)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status_code": -1,
			"status_msg":  err.Error(),
			"user_id": -1,
			"token": "-1",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg": "success",
		"user_id": &registerResponse.Identity,
		"token": &registerResponse.Token,
	})
}