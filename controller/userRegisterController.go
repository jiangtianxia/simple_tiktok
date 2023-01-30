/*
 * @Description:
 * @Author: liuxin
 * @Date: 2023-01-28 09:17:53
 * @LastEditTime: 2023-01-29 13:47:58
 * @LastEditors:
 */
package controller

import (
	"net/http"
	"simple_tiktok/service"

	"github.com/gin-gonic/gin"
)

func UserRegister(c *gin.Context) {
	username := c.Query("username")
	getpassword, _ := c.Get("password")
	password, ok := getpassword.(string)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"message":     "密码格式错误",
			"status_code": -1,
			"status_msg":  ok,
		})
		return
	}
	registerResponse, err := service.PostUserRegister(username, password)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status_code": -1,
			"status_msg":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"Identity":    &registerResponse.Identity,
	})
}
