package service

import (
	"github.com/gin-gonic/gin"
)

func GetFavoriteListByUserId(c *gin.Context, userId *uint64, loginUserId *uint64) (*[]Video, error) {
	ctx = c
	// 获取用户喜欢的视频列表
	// 遍历视频，调用之前写的函数获取视频的信息
	// 返回给controller层
	return nil, nil
}
