package service

import (
	"simple_tiktok/dao/mysql"
	"github.com/spf13/viper"
	"github.com/gin-gonic/gin"
	"simple_tiktok/utils"
	"fmt"
	"strconv"
	Myredis "simple_tiktok/dao/redis"
)

/** 
 * 获取视频列表的service层，决定是采用缓存策略还是数据库查询
 * @Author pjh
 * @Summary 
 * @Tags 
 **/

var ctx *gin.Context

func GetVideoListByUserId(c *gin.Context, authorId *uint64, loginUserId *uint64) (*[]Video, error){
	ctx = c

	// 1. 获取作者用户名
	authorName,err := getAuthorName(authorId)
	if err != nil {
		return nil, err
	}

	// 2. 创建作者对象
	author := Author{
		FollowCount: 0,
		FollowerCount: 0,
		IsFollow: false,
		Name: *authorName,
		Id: *authorId,
	}

	// 3. 尝试从缓存中获取用户发表的视频id列表
	key := fmt.Sprintf("%s%d", viper.GetString("redis.KeyUserPublishSetPrefix"), *authorId)
	n, err := utils.RDB4.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	// 若缓存中存在信息则使用 A 计划
	// A4. 创建返回的视频列表参数
	if n != 0 {
		videoIds, err := utils.RDB4.SMembers(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		videoList := &[]Video{}
		for i := range videoIds {
			sVideoId := videoIds[i]
			iVideoId, err := strconv.Atoi(sVideoId)
			if err != nil {
				return nil, err
			}
			videoId := uint64(iVideoId)

			// A5. 获取赞数
			favoriteCount, err := getVideoFavoriteCount(videoId)
			if err != nil {
				return nil, err
			}

			// A6. 获取评论数
			commentCount, err := getVideoCommentCount(videoId)
			if err != nil {
				return nil, err
			}

			// A7. 判断使用者是否喜欢该视频
			bIsFavorite, err := judgeLoginUserLoveVideo(videoId, *loginUserId)
			isFavorite := *bIsFavorite
			if err != nil {
				return nil, err
			}

			// A8. 获取视频封面url和视频url
			coverUrl, playUrl, title, err := tryToGetVideoInfo(&videoId)
			if err != nil {
				return nil, err
			}

			// A9. 创建单个视频对象
			*videoList = append(*videoList, Video{
				Id: videoId,
				Author: author,
				PlayUrl: *playUrl,
				CoverUrl: *coverUrl,
				FavoriteCount: *favoriteCount,
				CommentCount: *commentCount,
				IsFavorite: isFavorite,
				Title: *title,
			})
		}
		return videoList,nil
		
	}

	// 以下为 B 计划，需要对 用户发表的视频id列表 和 视频信息 进行缓存
	// 从数据库中获取视频信息
	videoListFromDao, err := mysql.QueryVideoList(authorId)
	if err != nil {
		return nil, err
	}

	// B4. 创建返回的视频列表参数
	videoList := &[]Video{}
	for i := range *videoListFromDao {
		videoId := (*videoListFromDao)[i].Identity 

		// B5. 获取赞数
		favoriteCount, err := getVideoFavoriteCount(videoId)
		if err != nil {
			return nil, err
		}

		// B6. 获取评论数
		commentCount, err := getVideoCommentCount(videoId)
		if err != nil {
			return nil, err
		}

		// B7. 判断使用者是否喜欢该视频
		bIsFavorite, err := judgeLoginUserLoveVideo(videoId, *loginUserId)
		isFavorite := *bIsFavorite
		if err != nil {
			return nil, err
		}

		// B8. 对视频信息进行缓存
		// 缓存用户发布的视频列表
		key = fmt.Sprintf("%s%d", viper.GetString("redis.KeyUserPublishSetPrefix"), *authorId)
		err = Myredis.RedisAddSetRDB4(key, fmt.Sprintf("%d", videoId))
		if err != nil {
			return nil, err
		}
		// 缓存视频信息
		err = Myredis.RedisAddVideoInfo((*videoListFromDao)[i])
		if err != nil {
			return nil, err
		}
		// B9. 创建单个视频对象
		*videoList = append(*videoList, Video{
			Id: videoId,
			Author: author,
			PlayUrl: (*videoListFromDao)[i].PlayUrl,
			CoverUrl: (*videoListFromDao)[i].CoverUrl,
			FavoriteCount: *favoriteCount,
			CommentCount: *commentCount,
			IsFavorite: isFavorite,
		})
	}

	return videoList, nil
}
func getAuthorName(authorId *uint64) (*string, error){
	authorName := new(string)
	key := fmt.Sprintf("%s%d",viper.GetString("redis.KeyUserHashPrefix"), *authorId)
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
			"identity":       *authorId,
			"username":       *authorName,
			"follow_count":   0,
			"follower_count": 0,
			"is_follow":      false,
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
			"identity":       *authorId,
			"username":       *authorName,
			"follow_count":   0,
			"follower_count": 0,
			"is_follow":      false,
		})
		if err != nil {
			return nil, err
		}
	}
	return authorName, nil
}

// 获取视频赞数的函数
//TODO
func getVideoFavoriteCount(videoId uint64) (*int64, error) {
	favoriteCount := new(int64)
	key := fmt.Sprintf("%s%d",viper.GetString("redis.KeyVideoFavoriteCountStringPrefix"), videoId)
	// 先从RDB0中查看键值对是否存在
	n, err := utils.RDB0.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	// redis中有key，获取string转换成int64
	if n != 0 {
		sFavoriteCount, err := utils.RDB0.Get(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		iFavoriteCount, err := strconv.Atoi(sFavoriteCount)
		if err != nil {
			return nil, err
		}
		*favoriteCount = int64(iFavoriteCount)
		return favoriteCount, nil
	}
	// 如果redis中没有key，调用mysql的函数获得状态
	favoriteCount, err = mysql.QueryVideoFavoriteCount(&videoId)
	if err != nil {
		return nil, err
	}
	err = Myredis.RedisAddStringRDB0(key, fmt.Sprintf("%d", *favoriteCount))
	if err != nil {
		return nil, err
	}
	return favoriteCount, nil
}

// 获取视频评论数的函数
func getVideoCommentCount(videoId uint64) (*int64, error) {
	commentCount := new(int64)
	key := fmt.Sprintf("%s%d",viper.GetString("KeyCommentListPrefix"), videoId)
	// 先从RDB0中查看键值对是否存在
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
		err = Myredis.RedisAddListRBD8(key, fmt.Sprintf("%d", (*commentList)[i].Identity))
		if err != nil {
			return nil, err
		}
		// 缓存评论信息
		commentKey := fmt.Sprintf("%s%d", viper.GetString("KeyCommentInfoHashPrefix"), (*commentList)[i].Identity)
		err = Myredis.RedisSetHashRDB7(commentKey, map[string]interface{}{
			"video_identity": (*commentList)[i].VideoIdentity,
			"user_identity": (*commentList)[i].UserIdentity,
			"text": (*commentList)[i].Text,
			"comment_time": (*commentList)[i].CommentTime,
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
// TODO
func judgeLoginUserLoveVideo(videoId uint64, loginUserId uint64) (*bool, error) {
	var isFavorite bool
	key := fmt.Sprintf("%s%d-%d",viper.GetString("redis.KeyUserLoveVideoStringPrefix"), videoId, loginUserId)
	// 先从RDB0中查看键值对是否存在
	n, err := utils.RDB0.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	// redis中有key，获取string转换成bool
	if n != 0 {
		isFavorite = true
		if utils.RDB0.Get(ctx, key).Val() == "0" {
			isFavorite = false
		}
		return &isFavorite, nil
	}
	// 如果redis中没有key，调用mysql的方法获得状态
	isFavorite, err = mysql.QueryIsFavorite(&videoId, &loginUserId)
	if err != nil {
		return nil, err
	}
	// 在redis中存储，0表示不喜欢，1表示喜欢
	if isFavorite {
		err = Myredis.RedisAddStringRDB0(key, "1")
	} else {
		err = Myredis.RedisAddStringRDB0(key, "0")
	}
	if err != nil {
		return nil, err
	}
	return &isFavorite, nil
}

// 从缓存中获取封面地址和视频播放地址，缓存中没有的话就从数据库中查询，并缓存视频的播放地址，封面地址
func tryToGetVideoInfo(videoId *uint64) (*string, *string, *string, error) {
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
		err = Myredis.RedisAddVideoInfo(*videoBasic)
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
		err = Myredis.RedisAddVideoInfo(*videoBasic)
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
		err = Myredis.RedisAddVideoInfo(*videoBasic)
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
		err = Myredis.RedisAddVideoInfo(*videoBasic)
		if err != nil {
			return nil, nil, nil, err
		}
		return &(*videoBasic).CoverUrl, &(*videoBasic).PlayUrl, &(*videoBasic).Title, nil
	}
	return &coverUrl, &playUrl, &title, nil
}

// 视频参数
type Video struct {
	Id uint64 `json:"id"`
	Author Author `json:"author"`
	PlayUrl string `json:"play_url"`
	CoverUrl string `json:"cover_url"`
	FavoriteCount int64 `json:"favorite_count"`
	CommentCount int64 `json:"comment_count"`
	IsFavorite bool `json:"is_favorite"`
	Title string `json:"title"`
}

// 作者参数
type Author struct {
	Id uint64 `json:"id"`
	Name string `json:"name"`
	FollowCount int64 `json:"follow_count"` // default
	FollowerCount int64 `json:"follower_count"` // default
	IsFollow bool `json:"is_follow"` // defalut
}
