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
var RDB *redis.Client

// 初始化Redis连接
func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConn"),
	})

	fmt.Println("redis inited ......")
}

func ReidsClose() {
	RDB.Close()
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

	fmt.Println("TokenBucket inited")
}
