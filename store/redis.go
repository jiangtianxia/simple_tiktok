package store

import (
	"fmt"
	"simple_tiktok/conf"

	"github.com/go-redis/redis/v9"
)

var rdb *redis.Client

func GetRDB() *redis.Client {
	return rdb
}

// 初始化redis
func InitRedis(c *conf.RedisConf) {
	rdb = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", c.RedisHost, c.RedisPort),
		Password:     c.RedisPassword,
		DB:           c.RedisDB,
		PoolSize:     1000,
		MinIdleConns: 100,
	})
}

func RedisClose() {
	rdb.Close()
}
