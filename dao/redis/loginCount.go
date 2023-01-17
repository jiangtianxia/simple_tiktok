package redis

import (
	"fmt"
	"simple_tiktok/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

/**
 * @Author jiang
 * @Description 记录登录失败次数防爆破
 * @Date 10:00 2023/1/17
 **/

const (
	LoginErrorKeyPrefix = "login_error:"
	LoginBankKeyPrefix  = "login_bank:"
)

// 获取ip
func GetRequestIP(ctx *gin.Context) string {
	reqIP := ctx.ClientIP()
	if reqIP == "::1" {
		reqIP = "127.0.0.1"
	}
	return reqIP
}

// 判断是否存在login_bank key，即判断用户是否有登录限制
func ExceLoginBank(ctx *gin.Context, userid string) (bool, float64) {
	key := LoginBankKeyPrefix + userid + GetRequestIP(ctx)
	fmt.Println(key)

	flag := utils.RDB.TTL(ctx, key).Val()
	if flag == -2 {
		return false, 0
	}
	fmt.Println(flag.Minutes())
	return true, flag.Minutes()
}

// 登录失败次数增加
// 1、如果不存在，则创建key，并设置过期时间为1分钟
// 2、如果存在，且val<=4，则val++
// 3、如果存在，且val==5，则创建login_bank key 过期时间为10分钟，并删除login_error key
func AddLoginError(ctx *gin.Context, userid string) error {
	key := LoginErrorKeyPrefix + userid + GetRequestIP(ctx)
	fmt.Println(key)
	bankkey := LoginBankKeyPrefix + userid + GetRequestIP(ctx)
	fmt.Println(bankkey)

	// 1、判断key是否存在
	flag := utils.RDB.Get(ctx, key).Val()
	if flag == "" {
		// 1、如果不存在，则创建key，并设置过期时间为1分钟
		return utils.RDB.Set(ctx, key, 1, time.Second*60).Err()
	}
	val, _ := strconv.Atoi(flag)

	if val <= 4 {
		// 2、如果存在，且val<=4，则val++
		return utils.RDB.Incr(ctx, key).Err()
	}

	// 3、如果存在，且val==5，则创建login_bank key 过期时间为10分钟，并删除login_error key
	// 事务操作
	pipeline := utils.RDB.Pipeline()
	pipeline.Del(ctx, key)
	pipeline.Set(ctx, bankkey, 1, time.Second*60*10)

	_, err := pipeline.Exec(ctx)
	return err
}
