package service

import (
	"context"
	"fmt"
	"simple_tiktok/dao/mysql"
	Myredis "simple_tiktok/dao/redis"
	"simple_tiktok/logger"
	"simple_tiktok/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/spf13/viper"
)

// 获取作者用户名
func getAuthorName(ctx *gin.Context, authorId *uint64) (*string, error) {
	authorName := new(string)
	key := fmt.Sprintf("%s%d", viper.GetString("redis.KeyUserHashPrefix"), *authorId)
	n, err := utils.RDB1.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if n == 0 {
		authorName, err = mysql.QueryAuthorName(authorId)
		if err != nil {
			return nil, err
		}
		err = Myredis.RedisAddUserInfo(ctx, key, map[string]interface{}{
			"identity": *authorId,
			"username": *authorName,
		})
		if err != nil {
			return nil, err
		}
		return authorName, nil
	}
	*authorName, _ = utils.RDB1.HGet(ctx, key, "username").Result()
	// 没有username字段
	if *authorName == "" {
		authorName, err = mysql.QueryAuthorName(authorId)
		if err != nil {
			return nil, err
		}
		err = Myredis.RedisAddUserInfo(ctx, key, map[string]interface{}{
			"identity": *authorId,
			"username": *authorName,
		})
		if err != nil {
			return nil, err
		}
	}
	return authorName, nil
}

// 获取视频赞数的函数
func getVideoFavoriteCount(ctx *gin.Context, videoId uint64) (*int64, error) {
	favoriteCount := new(int64)
	key := fmt.Sprintf("%s%d", viper.GetString("redis.KetFavoriteSetPrefix"), videoId)
	// 先从RDB5中查看键值对是否存在
	n, err := utils.RDB5.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	// redis中有key，使用ZCount计数
	if n != 0 {
		// 0分为不喜欢，1分为喜欢
		*favoriteCount, err = utils.RDB5.ZCount(ctx, key, "1", "2").Result()
		if err != nil {
			return nil, err
		}
		return favoriteCount, nil
	}
	// 如果redis中没有key，调用mysql的函数获得状态
	videoFavoriteList, err := mysql.QueryVideoFavoriteCount(&videoId)
	if err != nil {
		return nil, err
	}
	*favoriteCount = 0
	for i := range *videoFavoriteList {
		if (*videoFavoriteList)[i].Status == 0 {
			continue
		}
		err = Myredis.RedisAddZSetRDB5(ctx, key, fmt.Sprintf("%d", (*videoFavoriteList)[i].UserIdentity), 1)
		if err != nil {
			return nil, err
		}
		*favoriteCount++
	}
	return favoriteCount, nil
}

// 获取视频评论数的函数
func getVideoCommentCount(ctx *gin.Context, videoId uint64) (*int64, error) {
	commentCount := new(int64)
	key := fmt.Sprintf("%s%d", viper.GetString("redis.KeyCommentListPrefix"), videoId)
	// 先从RDB8中查看键值对是否存在
	n, err := utils.RDB8.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	// redis中有key，查找缓存列表中有多少个元素
	if n != 0 {
		*commentCount, err = utils.RDB8.LLen(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		return commentCount, nil
	}
	// 如果redis中没有key，调用mysql的函数获取所有评论对象，将评论id缓存进RDB8 -- 将评论信息缓存进RDB7
	commentList, err := mysql.QueryVideoCommentInfo(&videoId)
	if err != nil {
		return nil, err
	}
	for i := range *commentList {
		// 缓存评论id
		err = Myredis.RedisAddListRBD8(ctx, key, fmt.Sprintf("%d", (*commentList)[i].Identity))
		if err != nil {
			return nil, err
		}
		// 缓存评论信息
		commentKey := fmt.Sprintf("%s%d", viper.GetString("redis.KeyCommentInfoHashPrefix"), (*commentList)[i].Identity)
		err = Myredis.RedisSetHashRDB7(ctx, commentKey, map[string]interface{}{
			"video_identity": (*commentList)[i].VideoIdentity,
			"user_identity":  (*commentList)[i].UserIdentity,
			"text":           (*commentList)[i].Text,
			"comment_time":   (*commentList)[i].CommentTime,
		})
		if err != nil {
			return nil, err
		}
	}
	// 获取评论数
	*commentCount = int64(len(*commentList))
	return commentCount, nil
}

// 判断登录的用户是否喜欢指定视频
func judgeLoginUserLoveVideo(ctx *gin.Context, videoId uint64, loginUserId uint64) (*bool, error) {
	key := fmt.Sprintf("%s%d", viper.GetString("redis.KeyFavoriteUserSortSetPrefix"), videoId)
	isFavorite := new(bool)
	n, err := utils.RDB5.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if n != 0 {
		scroe, err := utils.RDB5.ZScore(ctx, key, fmt.Sprintf("%d", loginUserId)).Result()
		if err != nil {
			// sorted set 中不存在这个member，没有关联过
			if err.Error() == "redis: nil" {
				scroe = 0
			} else {
				return nil, err
			}
		}
		*isFavorite = false
		if scroe > 0 {
			*isFavorite = true
		}
		return isFavorite, nil
	}
	videoFavoriteList, err := mysql.QueryVideoFavoriteCount(&videoId)
	if err != nil {
		return nil, err
	}
	// 默认不喜欢
	*isFavorite = false
	// 缓存RDB5
	for i := range *videoFavoriteList {
		if (*videoFavoriteList)[i].Status == 0 {
			continue
		}
		err = Myredis.RedisAddZSetRDB5(ctx, key, fmt.Sprintf("%d", (*videoFavoriteList)[i].UserIdentity), 1)
		if err != nil {
			return nil, err
		}
		// 搜索到这个人喜欢
		if (*videoFavoriteList)[i].UserIdentity == loginUserId {
			*isFavorite = true
		}
	}
	return isFavorite, nil
}

// 从缓存中获取封面地址和视频播放地址，缓存中没有的话就从数据库中查询，并缓存视频的播放地址，封面地址
func tryToGetVideoInfo(ctx *gin.Context, videoId *uint64) (*string, *string, *string, error) {
	key := fmt.Sprintf("%s%d", viper.GetString("redis.KeyVideoInfoHashPrefix"), *videoId)
	// 判断视频键是否存在
	n, err := utils.RDB3.Exists(ctx, key).Result()
	if err != nil {
		return nil, nil, nil, err
	}
	// 视频键不存在，写入缓存直接返回
	if n == 0 {
		videoBasic, err := mysql.QueryVideoInfoByVideoId(videoId)
		if err != nil {
			return nil, nil, nil, err
		}
		// 缓存
		err = Myredis.RedisAddVideoInfo(ctx, *videoBasic)
		if err != nil {
			return nil, nil, nil, err
		}
		return &(*videoBasic).CoverUrl, &(*videoBasic).PlayUrl, &(*videoBasic).Title, nil
	}
	// 视频键存在，从缓存中读取封面地址
	coverUrl, err := utils.RDB3.HGet(ctx, key, "cover_url").Result()
	if err != nil {
		return nil, nil, nil, err
	}
	// 未使用该字段，从数据库中获取并写入缓存
	if coverUrl == "" {
		videoBasic, err := mysql.QueryVideoInfoByVideoId(videoId)
		if err != nil {
			return nil, nil, nil, err
		}
		// 缓存
		err = Myredis.RedisAddVideoInfo(ctx, *videoBasic)
		if err != nil {
			return nil, nil, nil, err
		}
		return &(*videoBasic).CoverUrl, &(*videoBasic).PlayUrl, &(*videoBasic).Title, nil
	}
	// 获取播放地址
	playUrl, err := utils.RDB3.HGet(ctx, key, "play_url").Result()
	if err != nil {
		return nil, nil, nil, err
	}
	// 未使用该字段，从数据库中获取并写入缓存
	if playUrl == "" {
		videoBasic, err := mysql.QueryVideoInfoByVideoId(videoId)
		if err != nil {
			return nil, nil, nil, err
		}
		// 缓存
		err = Myredis.RedisAddVideoInfo(ctx, *videoBasic)
		if err != nil {
			return nil, nil, nil, err
		}
		return &(*videoBasic).CoverUrl, &(*videoBasic).PlayUrl, &(*videoBasic).Title, nil
	}
	// 获取视频标题
	title, err := utils.RDB3.HGet(ctx, key, "title").Result()
	if title == "" {
		videoBasic, err := mysql.QueryVideoInfoByVideoId(videoId)
		if err != nil {
			return nil, nil, nil, err
		}
		// 缓存
		err = Myredis.RedisAddVideoInfo(ctx, *videoBasic)
		if err != nil {
			return nil, nil, nil, err
		}
		return &(*videoBasic).CoverUrl, &(*videoBasic).PlayUrl, &(*videoBasic).Title, nil
	}
	return &coverUrl, &playUrl, &title, nil
}

/**
 * @Author jiang
 * @Description 获取用户名
 * @Date 13:00 2023/2/13
 **/
func GetUsername(c *gin.Context, identity string) (string, error) {
	// 先查询缓存
	key := viper.GetString("redis.KeyUserHashPrefix") + identity

	if utils.RDB1.Exists(c, key).Val() == 0 {
		// 不存在，则在数据库中获取，同时将数据存入缓存
		username, err := mysql.FindUserName(identity)
		if err != nil {
			return "", err
		}

		err = Myredis.RedisAddUserInfoHash(c, key, map[string]interface{}{
			"identity": identity,
			"username": username,
		})
		if err != nil {
			return "", err
		}
		return username, err
	}

	// 存在，则获取，并返回
	userName, err := utils.RDB1.HGet(c, key, "username").Result()
	if err != nil {
		return "", err
	}
	return userName, nil
}

/**
 * @Author jiang
 * @Description 获取关注数
 * @Date 13:00 2023/2/13
 **/
func GetFollowCount(c *gin.Context, identity string) (int64, error) {
	// 先查询缓存
	key := viper.GetString("redis.KeyFollowListPrefix") + identity

	if utils.RDB10.Exists(c, key).Val() == 0 {
		// 不存在，则在数据库中获取，同时将数据存入缓存
		user_id, _ := strconv.Atoi(identity)
		followerList, err := mysql.FindUserFollowByIdentity(uint64(user_id))
		if err != nil {
			logger.SugarLogger.Error("FindUserFollowByIdentity Error：", err.Error())
			return 0, err
		}

		// 将数据存入缓存
		// 头部插入数据到key当中
		// LPUSH KEY_NAME VALUE1.. VALUEN
		pipeline := utils.RDB10.Pipeline()
		for _, follower := range followerList {
			pipeline.LPush(c, key, follower.FollowerIdentity)
		}
		pipeline.Expire(c, key, time.Duration(viper.GetInt("redis.RedisExpireTime"))*time.Hour)
		_, err = pipeline.Exec(c)
		return int64(len(followerList)), err
	}

	// 存在，则获取，并返回
	return utils.RDB10.LLen(c, key).Result()
}

/**
 * @Author jiang
 * @Description 获取粉丝数
 * @Date 13:00 2023/2/13
 **/
func GetFollowerCount(c *gin.Context, identity string) (int64, error) {
	// 先查询缓存
	key := viper.GetString("redis.KeyFollowerSortSetPrefix") + identity

	if utils.RDB11.Exists(c, key).Val() == 0 {
		// 不存在，则查询数据库，将数据存入缓存，并返回
		followerList, err := mysql.FindFollower(identity)
		if err != nil {
			return 0, nil
		}

		//  ZADD KEY_NAME SCORE1 VALUE1.. SCOREN VALUEN
		pipeline := utils.RDB11.Pipeline()
		for _, follower := range followerList {
			pipeline.ZAdd(c, key, redis.Z{
				Member: follower.UserIdentity,
				Score:  1,
			})
		}
		pipeline.Expire(c, key, time.Duration(viper.GetInt("redis.RedisExpireTime"))*time.Hour)
		_, err = pipeline.Exec(c)
		return int64(len(followerList)), err
	}

	// 存在，则获取，并返回
	return utils.RDB11.ZCount(c, key, "1", "1").Result()
}

/**
 * @Author jiang
 * @Description 判断是否关注该用户
 * @Date 13:00 2023/2/13
 **/
func IsFollow(c *gin.Context, identity string, follower string) (bool, error) {
	// 先查询缓存
	key := viper.GetString("redis.KeyFollowerSortSetPrefix") + identity

	if utils.RDB11.Exists(c, key).Val() == 0 {
		// 不存在，则查询数据库
		cnt, err := mysql.IsFollow(follower, identity)
		if err != nil {
			return false, err
		}

		if cnt <= 0 {
			return false, nil
		}

		return true, nil
	}

	// 存在，则获取，并返回
	cnt := utils.RDB11.ZScore(c, key, follower).Val()
	if cnt != 1 {
		return false, nil
	}
	return true, nil
}

// 将结果存入redis缓存
func SaveRedisResp(msgid string, code int, msg string) {
	info := map[string]interface{}{
		"status_code": code,
		"status_msg":  msg,
	}

	var ctx = context.Background()
	pipeline := utils.RDB0.Pipeline()
	pipeline.HSet(ctx, msgid, info)
	pipeline.Expire(ctx, msgid, time.Second*70)
	pipeline.Exec(ctx)
}

// 从缓存里取指定用户喜欢的视频id，如果没有就缓存进去
func getFavoriteVideoList(ctx *gin.Context, userId *uint64) (*[]uint64, error) {
	key := fmt.Sprintf("%s%d", viper.GetString("redis.KeyUserFavoriteListPrefix"), *userId)
	n, err := utils.RDB6.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var videoList []uint64
	if n == 0 {
		// 缓存中没有，需要在RDB6中缓存用户的信息
		// 先从数据库中拿到数据
		favoriteList, err := mysql.QueryFavoriteHistoryByUserId(userId)
		if err != nil {
			return nil, err
		}
		// 缓存的同时构造返回值
		for i := range *favoriteList {
			if (*favoriteList)[i].Status == 0 {
				continue
			}
			videoList = append(videoList, (*favoriteList)[i].VideoIdentity)
			err := Myredis.RedisAddListRDB6(ctx, key, fmt.Sprintf("%d", (*favoriteList)[i].VideoIdentity))
			if err != nil {
				return nil, err
			}
		}
		return &videoList, nil
	}
	res, err := utils.RDB6.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	for i := range res {
		id, err := strconv.Atoi(res[i])
		if err != nil {
			return nil, err
		}
		videoList = append(videoList, uint64(id))
	}
	return &videoList, nil
}

// 根据视频id找作者id，先从视频信息的缓存中找，没有就把视频信息加到缓存中去
func getAuthorIdByVideoId(ctx *gin.Context, videoId *uint64) (*uint64, error) {
	key := fmt.Sprintf("%s%d", viper.GetString("redis.KeyVideoInfoHashPrefix"), *videoId)
	n, err := utils.RDB3.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if n == 0 {
		videoBasic, err := mysql.QueryVideoInfoByVideoId(videoId)
		if err != nil {
			return nil, err
		}
		// 缓存
		err = Myredis.RedisAddVideoInfo(ctx, *videoBasic)
		if err != nil {
			return nil, err
		}
		return &videoBasic.UserIdentity, nil
	}
	sId, err := utils.RDB3.HGet(ctx, key, "author_id").Result()
	if err != nil {
		return nil, err
	}
	iId, err := strconv.Atoi(sId)
	if err != nil {
		return nil, err
	}
	authorId := uint64(iId)
	return &authorId, nil
}
