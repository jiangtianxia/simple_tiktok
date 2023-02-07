package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/**
 * @Author jiang
 * @Description 初始化配置
 * @Date 11:00 2023/1/15
 **/
func InitConfig() {
	/*
	* 发布模式
	 */
	// // 设置文件名，文件后缀默认为yml
	// // app为config下的名字
	// viper.SetConfigName("app")

	// // // 配置文件所在的文件夹
	// // viper.AddConfigPath(work + "/config")
	// viper.AddConfigPath("/build/config")

	// 获取当前工作目录
	work, _ := os.Getwd()

	// 设置文件名，文件后缀默认为yml
	// app为config下的名字
	viper.SetConfigName("app")

	// 配置文件所在的文件夹
	viper.AddConfigPath(work + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("config app inited ......")
	// // 获取全部文件内容
	// fmt.Println("config app:", viper.Get("app"))
	// fmt.Println("confg mysql:", viper.Get("mysql"))
}

/**
 * @Author jiang
 * @Description Mysql初始化
 * @Date 11:00 2023/1/15
 **/
var DB *gorm.DB

// 初始化数据库
func InitMysql() {

	/*
	* 发布模式
	 */
	// dns := viper.GetString("mysql.username") + ":" +
	// 	viper.GetString("mysql.password") + "@tcp(" + viper.GetString("mysql.addr") + ":" +
	// 	viper.GetString("mysql.port") + ")/" +
	// 	viper.GetString("mysql.database") + viper.GetString("mysql.base")
	// // 先休眠20秒，等数据库启动完成再连接
	// time.Sleep(time.Second * 20)
	// DB, _ = gorm.Open(mysql.Open(dns),
	// 	&gorm.Config{})

	// 自定义日志模板，打印SQL语句
	// os.Stdout 标准输出，控制台打印
	// 以\r\n来作为打印间隔
	// log.LastFlags 前面这串：2022/12/30 12:00
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢SQL阈值，超过1秒的sql查询会被记录到日志中
			LogLevel:      logger.Info, // 级别
			Colorful:      true,        // 彩色
		},
	)

	dns := viper.GetString("mysql.username") + ":" +
		viper.GetString("mysql.password") + "@tcp(" + viper.GetString("mysql.addr") + ":" +
		viper.GetString("mysql.port") + ")/" +
		viper.GetString("mysql.database") + viper.GetString("mysql.base")
	DB, _ = gorm.Open(mysql.Open(dns),
		&gorm.Config{Logger: newLogger})

	sqlDB, err := DB.DB()
	if err != nil {
		fmt.Println(err)
	}

	// 设置最大连接数
	sqlDB.SetMaxIdleConns(1000)
	// 设置最大的空闲连接数
	sqlDB.SetMaxOpenConns(100000)
	fmt.Println("mysql inited ......")
}

/**
 * @Author jiang
 * @Description redis初始化
 * @Date 11:00 2023/1/15
 **/
var RDB0 *redis.Client
var RDB1 *redis.Client
var RDB2 *redis.Client
var RDB3 *redis.Client
var RDB4 *redis.Client
var RDB5 *redis.Client
var RDB6 *redis.Client
var RDB7 *redis.Client
var RDB8 *redis.Client
var RDB9 *redis.Client
var RDB12 *redis.Client
var RDB13 *redis.Client

// 初始化Redis连接
func InitRedis() {
	// RDB0存储其他key
	RDB0 = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           0,
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})

	// RDB1存储用户信息hash集合
	RDB1 = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           1,
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})

	// RDB2存储上传视频列表
	RDB2 = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           2,
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})

	// RDB3存储视频信息hash集合
	RDB3 = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           3,
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})

	// RDB4存储用户发布视频列表
	RDB4 = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           4,
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})

	// RDB5存储视频的点赞用户sorted set
	RDB5 = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           5,
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})

	// RDB6存储用户点赞的视频列表
	RDB6 = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           6,
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})

	// RDB7存储评论信息hash集合
	RDB7 = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           7,
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})

	// RDB8存储视频的评论列表
	RDB8 = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           8,
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})

	// RDB9存储布隆过滤器
	RDB9 = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           9,
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})

	// RDB12存储发送信息sorted set
	RDB12 = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           12,
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})

	// RDB13存储聊天信息hash集合
	RDB13 = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           13,
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})

	fmt.Println("redis inited ......")
}

func ReidsClose() {
	RDB0.Close()
	RDB1.Close()
	RDB2.Close()
	RDB3.Close()
	RDB4.Close()
	RDB5.Close()
	RDB6.Close()
	RDB7.Close()
	RDB8.Close()
	RDB9.Close()
	RDB12.Close()
	RDB13.Close()
}

/**
 * @Author jiang
 * @Description 令牌桶初始化
 * @Date 9:00 2023/1/15
 **/
var Bucket TokenBucket

func InitCurrentLimit() {
	rate := viper.GetInt64("currentLimit.tokenBucket.rate")
	capacity := viper.GetInt64("currentLimit.tokenBucket.capacity")

	Bucket.Set(rate, capacity)

	fmt.Println("TokenBucket inited ...... ")
}
