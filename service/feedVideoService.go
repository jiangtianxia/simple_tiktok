package service

import (
	"simple_tiktok/dao/mysql"
	myredis "simple_tiktok/dao/redis"
	"simple_tiktok/logger"
	"simple_tiktok/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/spf13/viper"
)

/**
 * @Author
 * @Description 视频流接口
 * @Date 21:00 2023/2/11
 **/
func FeedVideo(c *gin.Context, user_id uint64, latestTime int64) ([]VideoInfo, int64, error) {

	hashKey := viper.GetString("redis.KeyVideoList")
	//先判断zset是否存在  zset没有设置过期时间
	if utils.RDB2.Exists(c, hashKey).Val() == 0 {
		return nil, 0, nil
	}

	// 使用缓存 查询出发布时间小于latestTime的30条记录  记录中包含视频的identity
	identityList := utils.RDB2.ZRevRangeByScore(c, hashKey, &redis.ZRangeBy{Min: strconv.Itoa(0), Max: strconv.FormatInt(latestTime-1, 10), Count: viper.GetInt64("feedVideoCnt")}).Val()
	var nextTime int64

	videoInfos := make([]VideoInfo, len(identityList))
	for i, identity := range identityList {
		if utils.RDB3.Exists(c, viper.GetString("redis.KeyVideoInfoHashPrefix")+identity).Val() == 0 {
			id, _ := strconv.Atoi(identity)
			video, err := mysql.FindVideoById(uint64(id))
			if err != nil {
				logger.SugarLogger.Error(err)
				return nil, 0, err
			}

			// 将数据库查询出的数据写入redis
			err = myredis.RedisAddVideoInfo(c, *video)
			if err != nil {
				logger.SugarLogger.Error(err)
				return nil, 0, err
			}
			// fmt.Println("数据库")
		}

		// 使用缓存
		video := utils.RDB3.HGetAll(c, viper.GetString("redis.KeyVideoInfoHashPrefix")+identity).Val()

		videoInfos[i].Id, _ = strconv.ParseUint(video["id"], 10, 64)
		userId, _ := strconv.ParseUint(video["author_id"], 10, 64)
		user, _ := mysql.FindUserByIdentity(userId)

		// 获取关注数
		followCount, err := GetFollowCount(c, strconv.Itoa(int(user.Identity)))
		if err != nil {
			logger.SugarLogger.Error("GetFollowCount Error：", err.Error())
			return nil, 0, err
		}

		// 获取粉丝数
		followerCount, err := GetFollowerCount(c, strconv.Itoa(int(user.Identity)))
		if err != nil {
			logger.SugarLogger.Error("GetFollowerCount Error：", err.Error())
			return nil, 0, err
		}

		// 判断是否关注该用户
		flag := false
		if user_id != 0 {
			flag, err = IsFollow(c, strconv.Itoa(int(user.Identity)), strconv.Itoa(int(user_id)))
			if err != nil {
				logger.SugarLogger.Error("IsFollow Error：", err.Error())
				return nil, 0, err
			}
		}

		// 获取赞数
		favoriteCount, err := getVideoFavoriteCount(c, videoInfos[i].Id)
		if err != nil {
			logger.SugarLogger.Error("getVideoFavoriteCount Error：", err.Error())
			return nil, 0, err
		}

		// 获取评论数
		commentCount, err := getVideoCommentCount(c, videoInfos[i].Id)
		if err != nil {
			logger.SugarLogger.Error("getVideoCommentCount Error：", err.Error())
			return nil, 0, err
		}

		// 判断使用者是否喜欢该视频
		bIsFavorite, err := judgeLoginUserLoveVideo(c, videoInfos[i].Id, user_id)
		isFavorite := *bIsFavorite
		if err != nil {
			logger.SugarLogger.Error("judgeLoginUserLoveVideo Error：", err.Error())
			return nil, 0, err
		}

		// 获取点赞数量，作品数和喜欢数
		totalFavourited, workCount, FavouriteCount, err := GetTotalFavouritedANDWorkCountANDFavoriteCount(user.Identity)
		if err != nil {
			logger.SugarLogger.Error("GetTotalFavouritedANDWorkCountANDFavoriteCount Error：", err.Error())
			return nil, 0, err
		}

		videoInfos[i].Author = Author{
			Id:              user.Identity,
			Name:            user.Username,
			FollowCount:     followCount,
			FollowerCount:   followerCount,
			IsFollow:        flag,
			Avatar:          viper.GetString("defaultAvatarUrl"),
			BackGroundImage: viper.GetString("defaultBackGroudImage"),
			Signature:       viper.GetString("defaultSignature"),
			TotalFavorited:  totalFavourited,
			WorkCount:       workCount,
			FavoriteCount:   FavouriteCount,
		}
		videoInfos[i].PlayUrl = viper.GetString("cos.addr") + video["play_url"]
		videoInfos[i].CoverUrl = viper.GetString("uploadAddr") + video["cover_url"]
		videoInfos[i].CommentCount = *commentCount
		videoInfos[i].FavoriteCount = *favoriteCount
		videoInfos[i].IsFavorite = isFavorite
		videoInfos[i].Title = video["title"]
		nextTime, _ = strconv.ParseInt(video["publish_time"], 10, 64)
	}

	return videoInfos, nextTime, nil
}
