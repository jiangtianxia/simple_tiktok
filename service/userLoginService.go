package service

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"simple_tiktok/dao/redis"
	"simple_tiktok/logger"
	"simple_tiktok/models"
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
	if utils.RDB.Exists(c, hashKey).Val() == 0 {
		// 查询数据库
		if err := models.IsExist(username); err != nil {

			if err.Error() == gorm.ErrRecordNotFound.Error() {
				return -1, "用户不存在"
			}
			return -1, "数据库错误"
		}
		passwordModels = models.QueryPassword(username)
		id = models.QueryIdentity(username)
	} else {
		passwordModels = utils.RDB.HGetAll(c, hashKey).Val()["password"]
		idPre, _ := strconv.Atoi(utils.RDB.HGetAll(c, hashKey).Val()["identity"])
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

	var res = map[string]interface{}{
		"identity": id,
		"username": username,
		"password": passwordModels,
	}

	// 判断是否有缓存
	if utils.RDB.Exists(c, hashKey).Val() == 0 {
		//新增缓存
		err := redis.RedisAddUserInfo(c, hashKey, res)
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
	return int64(id), token

}
