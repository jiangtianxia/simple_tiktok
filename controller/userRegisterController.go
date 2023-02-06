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
	var req service.RegisterRequire
	c.ShouldBindJSON(&req)

	registerResponse, err := service.PostUserRegister(c, &req)

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