package main

import (
	"math/rand"
	"simple_tiktok/conf"
	"simple_tiktok/logger"
	"simple_tiktok/router"
	"simple_tiktok/store"
	"simple_tiktok/utils"
	"time"

	"github.com/spf13/pflag"
)

var (
	DefaultConf     conf.DefaultConf
	SrvCnf          conf.ServerConf
	TiktokDBConf    conf.DbConf
	TiktokRedisConf conf.RedisConf
)

func init() {
	// 基础配置
	{
		pflag.UintVar(&SrvCnf.HTTPPort, "http_port", 8080, "server http port")
		pflag.UintVar(&SrvCnf.GRPCPort, "grpc_port", 8090, "server grpc port")
		pflag.StringVar(&SrvCnf.LogPath, "log_path", "./log", "log path")
		pflag.StringVar(&SrvCnf.LogFile, "log_file", "all.log", "log file")
		pflag.BoolVar(&SrvCnf.OpenCORS, "cors", true, "open cors")
		pflag.StringSliceVar(&SrvCnf.AllowOrigins, "origins", []string{""}, "allow origins")

		pflag.UintVar(&DefaultConf.MachineId, "machine_id", 2345576453432980, "snowflake machineID")

		pflag.StringVar(&DefaultConf.COSAddr, "cos_addr", "https://tiktok-1310814941.cos.ap-guangzhou.myqcloud.com", "cos addr")
		pflag.StringVar(&DefaultConf.COSSecretId, "cos_secret_id", "AKIDXwQpCd4OchWXR9ZEsOk5IRYS9ds9KVkA", "cos secret id")
		pflag.StringVar(&DefaultConf.COSSecretKey, "cos_secret_key", "UpyueCh2DVeieQErV2OaaM4bM5lGZuFl", "cos secret key")

		pflag.Int64Var(&DefaultConf.BucketRate, "bucket_rate", 1000, "bucket rate")
		pflag.Int64Var(&DefaultConf.BucketCapacity, "bucket_capacity", 5000, "bucket capacity")

		pflag.StringVar(&DefaultConf.Md5Salt, "md5_salt", "tiktokGi0I0R1tC#%", "md5 salt")
		pflag.StringVar(&DefaultConf.JwtKey, "jwt_key", "h2wnknlsd", "jwt key")
		pflag.IntVar(&DefaultConf.JwtExpire, "jwt_expire", 1200, "jwt expire")
		pflag.StringVar(&DefaultConf.HashSalt, "hash_salt", "tiktok123#%", "hash salt")

		pflag.StringVar(&DefaultConf.UploadBase, "upload_base", "./upload/", "upload base")
		pflag.StringVar(&DefaultConf.UploadAddr, "upload_addr", "http://127.0.0.1:8080", "upload addr")
	}

	// db
	{
		pflag.StringVar(&TiktokDBConf.Name, "tiktok_db", "tiktok", "tiktok db name")
		pflag.StringVar(&TiktokDBConf.Driver, "tiktok_driver", "mysql", "tiktok db driver")
		pflag.StringVar(&TiktokDBConf.Source, "tiktok_conn", "test:674092@tcp(159.75.164.227:3306)/tiktok?charset=utf8mb4&parseTime=True&loc=Local", "driver connection string")
	}

	// redis
	{
		pflag.StringVar(&TiktokRedisConf.Name, "tiktok_redis_name", "tiktok", "tiktok redis name")
		pflag.StringVar(&TiktokRedisConf.RedisHost, "tiktok_redis_host", "159.75.164.227", "tiktok redis host")
		pflag.StringVar(&TiktokRedisConf.RedisPassword, "tiktok_redis_pass", "579021", "tiktok redis password")
		pflag.IntVar(&TiktokRedisConf.RedisPort, "tiktok_redis_port", 6379, "tiktok redis port")
		pflag.IntVar(&TiktokRedisConf.RedisDB, "tiktok_redis_db", 0, "tiktok redis db")
	}

	pflag.Parse()
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	conf.SetGlobalConf(&DefaultConf)

	// 初始化日志
	logger.InitLogger(SrvCnf.LogPath, SrvCnf.LogFile)
	defer logger.SugarLogger.Sync()

	// 初始化雪花算法
	if err := utils.SnowflakeInit(uint16(DefaultConf.MachineId)); err != nil {
		panic(err)
	}

	// 初始化Mysql
	if err := store.InitDB(TiktokDBConf.Source); err != nil {
		panic(err)
	}
	if err := store.AutoMigrate(); err != nil {
		panic(err)
	}

	// 初始化redis
	store.InitRedis(&TiktokRedisConf)
	defer store.RedisClose()

	// 初始化令牌桶
	utils.InitCurrentLimit(DefaultConf.BucketRate, DefaultConf.BucketCapacity)

	// 初始化COS客户端
	store.InitCOS(DefaultConf.COSAddr, DefaultConf.COSSecretId, DefaultConf.COSSecretKey)

	// 初始化其他相关配置
	utils.SetMd5Salt(DefaultConf.Md5Salt)
	utils.InitJwt(DefaultConf.JwtKey, DefaultConf.JwtExpire)
	if err := utils.InitHashID(DefaultConf.HashSalt, 8); err != nil {
		panic(err)
	}

	router.Start(&SrvCnf)
}
