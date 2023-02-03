package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/spf13/viper"
	"simple_tiktok/dao/mysql"
	myredis "simple_tiktok/dao/redis"
	"simple_tiktok/logger"
	"simple_tiktok/utils"
	"strconv"
	"time"
)

type VideoInfo struct {
	Identity      uint64 `json:"id"`        // 视频唯一标识
	Author        Author `json:"author"`    // 作者信息
	PlayUrl       string `json:"play_url"`  // 视频路径
	CoverUrl      string `json:"cover_url"` // 封面路径
	FavoriteCount int    `json:"favorite_count"`
	CommentCount  int    `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Title         string `json:"title"` // 视频标题
}

type Author struct {
	Id            uint64 `json:"id"`
	Name          string `json:"name"`
	FollowCount   int    `json:"follow_count"`
	FollowerCount int    `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

func FeedVideo(c *gin.Context, latestTime time.Time) ([]VideoInfo, time.Time, error) {

	hashKey := viper.GetString("redis.KeyVideoList")
	//先判断zset是否存在  zset没有设置过期时间
	if utils.RDB2.Exists(c, hashKey).Val() == 0 {
		return nil, time.Time{}, nil
	}

	// 使用缓存 查询出发布时间小于latestTime的30条记录  记录中包含视频的identity
	identityList := utils.RDB2.ZRevRangeByScore(c, hashKey, &redis.ZRangeBy{Min: strconv.Itoa(0), Max: strconv.FormatInt(latestTime.Unix(), 10), Count: viper.GetInt64("feedVideoCnt")}).Val()
	var nextTime time.Time

	videoInfos := make([]VideoInfo, len(identityList))
	for i, identity := range identityList {
		if utils.RDB3.Exists(c, viper.GetString("redis.KeyVideoInfoHashPrefix")+identity).Val() == 0 {
			id, _ := strconv.Atoi(identity)
			video, err := mysql.FindVideoById(uint64(id))
			if err != nil {
				logger.SugarLogger.Error(err)
			}
			myredis.RedisAddVideoInfo(*video) // 将数据库查询出的数据写入redis
			fmt.Println("数据库")
		}

		timeTmp := "2006-01-02 15:04:05" // 定义时间模板
		loc, _ := time.LoadLocation("local")
		// 使用缓存
		video := utils.RDB3.HGetAll(c, viper.GetString("redis.KeyVideoInfoHashPrefix")+identity).Val()

		videoInfos[i].Identity, _ = strconv.ParseUint(video["id"], 10, 64)
		userId, _ := strconv.ParseUint(video["author_id"], 10, 64)
		user, _ := mysql.FindUserByIdentity(userId)
		videoInfos[i].Author = Author{
			Id:            user.Identity,
			Name:          user.Username,
			FollowCount:   0,
			FollowerCount: 0,
			IsFollow:      false,
		}
		videoInfos[i].PlayUrl = video["play_url"]
		videoInfos[i].CoverUrl = video["cover_url"]
		videoInfos[i].CommentCount = 0
		videoInfos[i].IsFavorite = false
		videoInfos[i].Title = video["title"]
		nextTime, _ = time.ParseInLocation(timeTmp, video["publish_time"], loc)
		fmt.Println("缓存")
	}

	return videoInfos, nextTime, nil
}
