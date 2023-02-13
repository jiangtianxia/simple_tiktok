package main

import (
	"fmt"
	"simple_tiktok/logger"
	rocket "simple_tiktok/rocketmq"
	"simple_tiktok/router"
	"simple_tiktok/utils"

	"github.com/apache/rocketmq-client-go/v2/rlog"
	"github.com/spf13/viper"
)

func main() {
	/*
	* 初始化
	 */
	// 初始化配置
	utils.InitConfig()

	// 初始化雪花算法
	if err := utils.SnowflakeInit(viper.GetUint16("snowflake.machineID")); err != nil {
		fmt.Println("init Snowflake failed, err:", err)
		return
	}
	fmt.Println("Snowflake inited ...... ")

	// 初始化日志
	logger.InitLogger()
	defer logger.SugarLogger.Sync() // 刷新流，写日志到文件中

	// 初始化Mysql和Redis
	utils.InitMysql()
	utils.InitRedis()
	defer utils.ReidsClose()

	// 初始化令牌桶
	utils.InitCurrentLimit()

	// 初始化布隆过滤器
	//utils.InitBloomFilter()

	// 初始化rocketmq
	rocket.InitRocketmq()
	rlog.SetLogLevel("error") // 控制台只打印rocketmq的error日志
	//初始化Comment的消息队列，并开启消费
	//rocket.Initcommentactionmq() //初始化CommentAction的消息队列，并开启消费

	// 初始化COS客户端
	utils.InitCos()

	// 初始化熔断器
	utils.InitCircuitBreaker()

	logger.SugarLogger.Info("初始化配置完成")

	// 配置路由
	r := router.Router()
	logger.SugarLogger.Info("配置路由完成")
	r.Run("127.0.0.1:8080")
}
