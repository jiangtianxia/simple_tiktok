package service

import (
	"errors"
	"simple_tiktok/dao/mysql"
	myRedis "simple_tiktok/dao/redis"
	"simple_tiktok/logger"
	"simple_tiktok/models"
	"simple_tiktok/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

/**
 * @Author: liuxin
 * @Description: 用户注册接口
 * @Date: 2023-01-28 09:33:15
 **/
func PostUserRegister(c *gin.Context, req *RegisterRequire) (RegisterResponse, error) {
	//验证合法性
	if req.Username == "" {
		return RegisterResponse{}, errors.New("用户名为空")
	}
	if req.Password == "" {
		return RegisterResponse{}, errors.New("密码为空")
	}

	//判断用户名
	if mysql.UserIsExist(req.Username) {
		return RegisterResponse{}, errors.New("用户名已存在")
	}

	//雪花算法生成id
	getid, err := utils.GetID()
	getpsw := utils.MakePassword(req.Password)
	if err != nil {
		return RegisterResponse{}, errors.New("注册失败")
	}

	ur := models.UserBasic{Identity: getid, Username: req.Username, Password: getpsw}
	// 更新数据
	err = mysql.AddUserBasic(ur)
	if err != nil {
		return RegisterResponse{}, err
	}

	//给token
	token, err := utils.GenerateToken(ur.Identity, ur.Username)
	if err != nil {
		logger.SugarLogger.Error("Generate Token Error:" + err.Error())
		return RegisterResponse{}, errors.New("注册失败")
	}

	//response
	resp := RegisterResponse{}
	resp.Token = token
	resp.Identity = ur.Identity

	// 添加缓存
	var newCathe = map[string]interface{}{
		"identity": ur.Identity,
		"username": ur.Username,
	}
	hashKey := viper.GetString("redis.KeyUserHashPrefix") + strconv.Itoa(int(ur.Identity))
	err = myRedis.RedisAddUserInfo(c, hashKey, newCathe)
	if err != nil {
		logger.SugarLogger.Error(err)
	}
	return resp, nil
}
