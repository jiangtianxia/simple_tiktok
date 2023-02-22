package redis

import (
	"simple_tiktok/dao/mysql"
	"simple_tiktok/models"
	"simple_tiktok/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/spf13/viper"
)

/**
 * @Author jiang
 * @Description 使用有序集合 按照视频发布时间降序存储视频的identity 和 发布时间
 * @Date 19:00 2023/1/29
 **/
func RedisAddVideoList(ctx *gin.Context, publishTime int64, datavideoId string) error {
	// 1、获取key
	key := viper.GetString("redis.KeyVideoList")

	// 2、头部插入数据到key当中
	// LPUSH KEY_NAME VALUE1.. VALUEN
	err := utils.RDB2.ZAdd(ctx, key, redis.Z{Score: float64(publishTime), Member: datavideoId}).Err()
	return err
}

/**
 * @Author jiang
 * @Description 存储视频信息到hash集合
 * @Date 19:00 2023/1/29
 **/
func RedisAddVideoInfo(ctx *gin.Context, videoInfo models.VideoBasic) error {
	// 1、获取前缀，拼接key
	key := viper.GetString("redis.KeyVideoInfoHashPrefix") + strconv.Itoa(int(videoInfo.Identity))

	// 2、存到缓存
	value := map[string]interface{}{
		"id":           videoInfo.Identity,
		"author_id":    videoInfo.UserIdentity,
		"play_url":     videoInfo.PlayUrl,
		"cover_url":    videoInfo.CoverUrl,
		"title":        videoInfo.Title,
		"publish_time": videoInfo.PublishTime,
	}
	// 开启事务
	pipeline := utils.RDB3.Pipeline()
	pipeline.HSet(ctx, key, value)
	pipeline.Expire(ctx, key, time.Duration(viper.GetInt("redis.RedisExpireTime"))*time.Hour)

	_, err := pipeline.Exec(ctx)
	return err
}

/**
 * @Author jiang
 * @Description 使用列表按照用户发布时间降序存储视频的identity
 * @Date 19:00 2023/1/29
 **/
func RedisAddPublishList(ctx *gin.Context, userId string, videoId string) error {
	// 1、获取前缀，拼接key
	key := viper.GetString("redis.KeyPublishListPrefix") + userId

	// 2、判断key是否存在
	if utils.RDB4.Exists(ctx, key).Val() == 0 {
		// 1）不存在，查询数据库，将数据存到缓存，当中
		identity, _ := strconv.Atoi(userId)
		videoList, err := mysql.FindVideoByUserIdentity(uint64(identity))
		if err != nil {
			return err
		}

		// 头部插入数据到key当中
		// LPUSH KEY_NAME VALUE1.. VALUEN
		pipeline := utils.RDB4.Pipeline()
		for _, video := range videoList {
			pipeline.LPush(ctx, key, video.Identity)
		}
		pipeline.Expire(ctx, key, time.Duration(viper.GetInt("redis.RedisExpireTime"))*time.Hour)
		_, err = pipeline.Exec(ctx)
		return err
	}

	// 头部插入数据到key当中
	// LPUSH KEY_NAME VALUE1.. VALUEN
	pipeline := utils.RDB4.Pipeline()
	pipeline.LPush(ctx, key, videoId)
	pipeline.Expire(ctx, key, time.Duration(viper.GetInt("redis.RedisExpireTime"))*time.Hour)

	_, err := pipeline.Exec(ctx)
	return err
}
