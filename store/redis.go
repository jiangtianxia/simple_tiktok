package store

import (
	"context"
	"fmt"
	"simple_tiktok/conf"
	"time"

	"github.com/go-redis/redis/v9"
)

const (
	VideoHashPrefix   = "tiktok:videoInfo:"
	UserHashPrefix    = "tiktok:userInfo:"
	FavoriteSetPrefix = "tiktok:favorite:"
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

func NewRedisUserHash(username string, workCount, favoriteCount int64) map[string]interface{} {
	return map[string]interface{}{
		"username":       username,
		"work_count":     workCount,
		"favorite_count": favoriteCount,
	}
}

func NewRedisVideoHash(title string, playUrl string, coverUrl string, authorId uint, commentCount int64, favoriteCount int64) map[string]interface{} {
	return map[string]interface{}{
		"title":          title,
		"play_url":       playUrl,
		"cover_url":      coverUrl,
		"author_id":      authorId,
		"comment_count":  commentCount,
		"favorite_count": favoriteCount,
	}
}

func CreateHashData(data map[uint]map[string]interface{}, prefix string) error {
	pipeline := rdb.Pipeline()
	for id, value := range data {
		key := fmt.Sprintf("%s%d", prefix, id)
		pipeline.HSet(context.Background(), key, value)
		pipeline.Expire(context.Background(), key, 7*24*time.Hour)
	}
	_, err := pipeline.Exec(context.Background())
	return err
}

func GetHashData(ids []uint, prefix string) ([]bool, []map[string]string, error) {
	pipeline := rdb.Pipeline()
	ctx := context.Background()
	cmds := make([]*redis.MapStringStringCmd, len(ids))
	for i, id := range ids {
		key := fmt.Sprintf("%s%d", prefix, id)
		cmds[i] = pipeline.HGetAll(ctx, key)
	}

	_, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, nil, err
	}

	exists := make([]bool, len(ids))
	result := make([]map[string]string, len(ids))
	for i, cmd := range cmds {
		data, err := cmd.Result()
		if err != nil {
			return nil, nil, err
		}
		if len(data) == 0 {
			exists[i] = false
			continue
		}

		exists[i] = true
		result[i] = data
	}
	return exists, result, nil
}
