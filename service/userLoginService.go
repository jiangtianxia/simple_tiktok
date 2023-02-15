package service

import (
	"errors"
	"fmt"
	"simple_tiktok/dao/mysql"
	myRedis "simple_tiktok/dao/redis"
	"simple_tiktok/logger"
	"simple_tiktok/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

/**
 * @Author
 * @Description 用户登录接口
 * @Date
 **/
func Login(c *gin.Context, username string, password string) (LoginResponse, error) {

	//验证输入
	if username == "" {
		logger.SugarLogger.Error("用户名为空")
		return LoginResponse{}, errors.New("用户名为空")
	}
	if password == "" {
		logger.SugarLogger.Error("密码为空")
		return LoginResponse{}, errors.New("密码为空")
	}

	//判断用户1分钟内是否登录了5次。
	flag, time := myRedis.ExceLoginBank(c, username)
	if flag {
		return LoginResponse{}, fmt.Errorf("1分钟内连续登录5次, 登录限制, 请 %f 分钟后重试", time)
	}

	// 查询数据库
	if !mysql.UserIsExist(username) {
		return LoginResponse{}, errors.New("用户不存在")
	}

	user := mysql.QueryUserInfo(username)

	//判断密码是否正确，不正确增加失败次数次数
	if !utils.ValidPassword(password, user.Password) {
		//不正确增加失败次数次数
		if err := myRedis.AddLoginError(c, username); err != nil {
			logger.SugarLogger.Error("fail times Error:" + err.Error())
			return LoginResponse{}, errors.New("登录失败")
		}
		return LoginResponse{}, errors.New("密码错误, 1分钟内连续登录失败5次, 该账号会被锁10分钟")
	}

	//token
	token, err := utils.GenerateToken(user.Identity, username) //参数有问题
	if err != nil {
		logger.SugarLogger.Error("Generate Token Error:" + err.Error())
		return LoginResponse{}, errors.New("登录失败")
	}

	var userLogin = LoginResponse{
		Identity: user.Identity,
		Token:    token,
	}

	// 将数据添加到缓存
	var newCathe = map[string]interface{}{
		"identity": user.Identity,
		"username": user.Username,
	}
	hashKey := viper.GetString("redis.KeyUserHashPrefix") + strconv.Itoa(int(user.Identity))
	err = myRedis.RedisAddUserInfo(c, hashKey, newCathe)
	if err != nil {
		logger.SugarLogger.Error(err)
	}
	return userLogin, nil
}
