package service

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"simple_tiktok/dao/redis"
	"simple_tiktok/logger"
	"simple_tiktok/models"
	"simple_tiktok/utils"
)

// UserService 用户注册服务
type UserService struct {
	UserName string `form:"user_name" json:"user_name" binding:"required,min=3,max=15" example:"FanOne"`
	Password string `form:"password" json:"password" binding:"required,min=5,max=16" example:"FanOne666"`
}

// Login 用户登陆函数
func Login(c *gin.Context, username string, password string) (string, string) {

	//验证
	if username == "" {
		return "-1", "用户名为空"
	}
	if password == "" {
		return "-1", "密码为空"
	}

	if err := models.IsExist(username); err != nil {
		// 如果查询不到，返回相应的错误
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return "-1", "用户不存在"
		}
		return "-1", "数据库错误"
	}

	//判断用户1分钟内是否登录了5次。
	if times, _ := redis.ExceLoginBank(c, models.QueryIdentity(password)); times == true {
		return "-1", "登录限制"
	}

	//判断密码是否正确，不正确增加失败次数次数
	if !utils.ValidPassword(password, models.QueryPassword(username)) {
		//不正确增加失败次数次数
		if err := redis.AddLoginError(c, models.QueryIdentity(username)); err != nil {
			return "-1", "增加失败错误失败"
		}
		return "-1", "密码错误"
	} else {
		//token
		token, err := utils.GenerateToken(models.QueryIdentity(username), username, password) //参数有问题
		if err != nil {
			logger.SugarLogger.Error("Generate Token Error:" + err.Error())
			return "-1", "token错误"
		}
		return models.QueryIdentity(username), token
	}

}
