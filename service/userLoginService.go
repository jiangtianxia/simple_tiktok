package service

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"simple_tiktok/dao/mysql"
	"simple_tiktok/dao/redis"
	"simple_tiktok/logger"
	rocket "simple_tiktok/rocketmq"
	"simple_tiktok/utils"
)

// Login 用户登陆函数
func Login(c *gin.Context, username string, password string) (map[string]interface{}, error) {

	//验证输入
	if username == "" {
		return nil, errors.New("用户名为空")
	}
	if password == "" {
		return nil, errors.New("密码为空")
	}

	//判断用户1分钟内是否登录了5次。
	if times, _ := redis.ExceLoginBank(c, username); times == true {
		return nil, errors.New("登录限制")
	}

	// 查询数据库
	if err := mysql.IsExist(username); err != nil {

		if err.Error() == gorm.ErrRecordNotFound.Error() {
			logger.SugarLogger.Error("Generate Token Error:" + err.Error())
			return nil, err
		}
		logger.SugarLogger.Error("Generate Token Error:" + err.Error())
		return nil, err
	}

	user := mysql.QueryInfo(username)

	//判断密码是否正确，不正确增加失败次数次数
	if !utils.ValidPassword(password, user.Password) {
		//不正确增加失败次数次数
		if err := redis.AddLoginError(c, username); err != nil {
			logger.SugarLogger.Error("Generate Token Error:" + err.Error())
			return nil, err
		}
		return nil, errors.New("密码错误")
	}

	//token
	token, err := utils.GenerateToken(user.Identity, username) //参数有问题
	if err != nil {
		logger.SugarLogger.Error("Generate Token Error:" + err.Error())
		return nil, err
	}

	var userLogin = map[string]interface{}{
		"identity": user.Identity,
		"username": username,
		"password": user.Password,
		"token":    token,
	}

	// 将数据发送到消息队列
	redisTopic := viper.GetString("rocketmq.redisTopic1")
	Producer := viper.GetString("rocketmq.redisProducer")
	tag := viper.GetString("rocketmq.userLoginTag")
	data, _ := json.Marshal(userLogin)
	rocket.SendMsg(c, Producer, redisTopic, tag, data)
	return userLogin, nil
}
