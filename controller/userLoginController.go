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

	userlogin, err := service.Login(c, username, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status_code": -1,
			"status_msg":  err.Error(),
		})
		return
	}

	////布隆过滤器
	//if !utils.Filter.Check(username) {
	//	c.JSON(http.StatusInternalServerError, gin.H{
	//		"status_code": -1,
	//		"status_msg":  "no such user",
	//	})
	//	return
	//}

	c.JSON(
		http.StatusOK,
		gin.H{
			"status_code": 0,
			"status_msg":  "string",
			"identity":    userlogin["identity"],
			"token":       userlogin["token"],
		},
	)

}
