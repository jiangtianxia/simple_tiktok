/*
 * @Description:
 * @Author: liuxin
 * @Date: 2023-01-28 09:33:15
 * @LastEditTime: 2023-01-29 13:45:24
 * @LastEditors:
 */
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
)

//请求体
type RegisterRequire struct {
	Username string 
	Password    string 
}

type RegisterResponse struct {
	Identity uint64 `json:"identity"`
	Token    string `json:"token"`
}

// 注册得到token和id
func PostUserRegister(c *gin.Context, req *RegisterRequire) (*RegisterResponse, error) {
	//验证合法性
	if req.Username == "" {
		return nil, errors.New("用户名为空")
	}
	if req.Password == "" {
		return nil, errors.New("密码为空")
	}

	//判断用户名
	if mysql.IsExist(req.Username) {
		return nil, errors.New("用户名已存在")
	}

	redisErr := myRedis.RedisUserRegister(c, req.Username, map[string]interface{}{
		"identity": -1,
	})
	if redisErr != nil {
		logger.SugarLogger.Error()
		return nil, redisErr
	}
	
	//雪花算法生成id
	getid, err := utils.GetID()
	getpsw := utils.MakePassword(req.Password)
	if err != nil {
		return nil, err
	}

	ur := models.UserBasic{Identity: getid, Username: req.Username, Password: getpsw}
	var res1 = map[string]interface{}{
		"identity":       getid,
		"username":       ur.Username,
		"password":       getpsw,
	}
	
	// 新增缓存
	err = myRedis.RedisUserRegister(c, req.Username, res1)
	
	if err != nil {
		logger.SugarLogger.Error(err)
		return nil, err
	}
	
	// 使用缓存
	cathe := utils.RDB1.HGetAll(c, req.Username).Val()
	identity, _ := strconv.Atoi(cathe["identity"])
	
	res2 := models.UserBasic{Identity: (uint64)(identity), Username: cathe["username"], Password: cathe["password"]}
	
	//更新数据
	err = mysql.AddUserBasic(&res2)
	if err != nil {
		return nil, err
	}

	//给token
	token, err3 := utils.GenerateToken(ur.Identity, ur.Username)
	if err3 != nil {
		logger.SugarLogger.Error("Generate Token Error:" + err.Error())
		return nil, err3
	}

	//response
	resp := RegisterResponse{}
	resp.Token = token
	resp.Identity = ur.Identity

	return &resp, nil
}
