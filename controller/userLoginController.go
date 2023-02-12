package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"simple_tiktok/service"
)

func Userlogin(c *gin.Context) {
	//获取参数
	username := c.Query("username")
	password := c.Query("password")

	login, err := service.Login(c, username, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": -1,
			"status_msg":  err.Error(),
		})
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status_code": 0,
			"status_msg":  "string",
			"identity":    login["identity"],
			"token":       login["token"],
		},
	)

}
