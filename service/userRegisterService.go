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

type UserRegister struct {
	identity uint64
	username string
	password string

	token string
	data  *RegisterResponse
}

// 注册得到token和id
func PostUserRegister(c *gin.Context,username, password string) (*RegisterResponse, error) {
	return NewPostUserRegister(c, username, password).Solve(c, username, password)
}

func NewPostUserRegister(c *gin.Context,username, password string) *UserRegister {
	return &UserRegister{username: username, password: password}
}

func (req *UserRegister) Solve(c *gin.Context, username, password string) (*RegisterResponse, error) {
	//验证
	if username == "" {
		return nil, errors.New("用户名为空")
	}
	if password == "" {
		return nil, errors.New("密码为空")
	}

	//更新数据
	if err := req.update(c, username, password); err != nil {
		return nil, err
	}
	//大概这里是要新增缓存的地方

	//response
	if err := req.registerResponse(); err != nil {
		return nil, err
	}
	return req.data, nil
}

func (req *UserRegister) update(c *gin.Context, username, password string) error {
	getid, err := utils.GetID()
	if err != nil {
		return err
	}
	
	//判断用户名
	ret1 := mysql.NewUserRegisterDAO()
	if ret1.IsExist(req.username) {
		return errors.New("用户名已存在")
	}

	// 添加缓存这里参考zxy的，暂未测试
	redisErr := myRedis.RedisUserRegister(c, username, map[string]interface{}{
		"identity": -1,
	})
	if redisErr != nil {
		logger.SugarLogger.Error(err)
		return redisErr
	}

	ur := models.UserBasic{Identity: getid, Username: utils.MakePassword(req.username), Password: utils.MakePassword(req.password)}
	var res1 = map[string]interface{}{
		"identity":       ur.Identity,
		"username":       ur.Username,
		"password":       ur.Password,
	}

	// 新增缓存
	err = myRedis.RedisAddUserInfo(c, username, res1)
	if err != nil {
		logger.SugarLogger.Error(err)
		return err
	}
	// 使用缓存
	cathe := utils.RDB0.HGetAll(c, username).Val()
	identity, _ := strconv.Atoi(cathe["identity"])

	res2 := models.UserBasic{Identity: (uint64)(identity), Username: cathe["username"], Password: cathe["password"]}

	//更新数据
	ret2 := mysql.NewUserBasicDAO()
	err = ret2.AddUserBasic(&res2)
	if err != nil {
		return err
	}

	//给token
	token, err3 := utils.GenerateToken(ur.Identity, ur.Username)
	if err3 != nil {
		logger.SugarLogger.Error("Generate Token Error:" + err.Error())
		return err3
	}
	req.token = token
	req.identity = ur.Identity
	return nil
}

func (req *UserRegister) registerResponse() error {
	req.data = &RegisterResponse{
		Identity: req.identity,
		Token:    req.token,
	}
	return nil
}

type RegisterResponse struct {
	Identity uint64 `json:"identity"`
	Token    string `json:"token"`
}
