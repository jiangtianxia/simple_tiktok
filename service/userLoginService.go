package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"simple_tiktok/dao/mysql"
	"simple_tiktok/dao/redis"
	"simple_tiktok/logger"
	rocket "simple_tiktok/rocketmq"
	"simple_tiktok/utils"
	"strconv"
)

// Login 用户登陆函数
func Login(c *gin.Context, username string, password string) (int64, string) {

	//验证输入
	if username == "" {
		return -1, "用户名为空"
	}
	if password == "" {
		return -1, "密码为空"
	}

	//判断用户1分钟内是否登录了5次。
	if times, _ := redis.ExceLoginBank(c, username); times == true {
		return -1, "登录限制"
	}

	var passwordModels string
	var id uint64

	//查询用户名是否存在
	//redis查询
	hashKey := viper.GetString("redis.KeyUserHashPrefix") + username
	// 判断是否有缓存
	if utils.RDB1.Exists(c, hashKey).Val() == 0 {
		// 查询数据库
		if err := mysql.IsExist(username); err != nil {

			if err.Error() == gorm.ErrRecordNotFound.Error() {
				return -1, "用户不存在"
			}
			return -1, "数据库错误"
		}
		passwordModels = mysql.QueryPassword(username)
		id = mysql.QueryIdentity(username)
	} else {
		passwordModels = utils.RDB1.HGetAll(c, hashKey).Val()["password"]
		idPre, _ := strconv.Atoi(utils.RDB0.HGetAll(c, hashKey).Val()["identity"])
		id = uint64(idPre)
	}

	//判断密码是否正确，不正确增加失败次数次数
	if !utils.ValidPassword(password, passwordModels) {
		//不正确增加失败次数次数
		if err := redis.AddLoginError(c, username); err != nil {
			return -1, "增加失败错误失败"
		}
		return -1, "密码错误"
	}

	var userLogin = map[string]interface{}{
		"identity": id,
		"username": username,
		"password": passwordModels,
	}

	// 判断是否有缓存
	if utils.RDB1.Exists(c, hashKey).Val() == 0 {
		//新增缓存
		err := redis.RedisAddUserInfo(c, hashKey, userLogin)
		if err != nil {
			logger.SugarLogger.Error(err)
			return -1, "新增失败"
		}
		//fmt.Println("数据库")
	}

	//token
	token, err := utils.GenerateToken(id, username) //参数有问题
	if err != nil {
		logger.SugarLogger.Error("Generate Token Error:" + err.Error())
		return -1, "token错误"
	}

	// 将数据发送到消息队列
	redisTopic := viper.GetString("rocketmq.redisTopic")
	Producer := viper.GetString("rocketmq.redisProducer")
	tag := viper.GetString("rocketmq.userLoginTag")
	data, _ := json.Marshal(userLogin)
	rocket.SendMsg(c, Producer, redisTopic, tag, data)
	return int64(id), token

}
