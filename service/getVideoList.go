package service

import (
	"simple_tiktok/dao/mysql"
	"github.com/spf13/viper"
	"simple_tiktok/utils"
	"simple_tiktok/models"
	"fmt"
	"context"
	"strconv"
	Myredis "simple_tiktok/dao/redis"
)

/** 
 * 获取视频列表的service层，决定是采用缓存策略还是数据库查询
 * @Author pjh
 * @Summary 
 * @Tags 
 **/

func GetVideoListByUserId(authorId *uint64, loginUserId *uint64) (*[]Video, error){
	ctx := context.Background()
	// 1. 获取作者用户名
	authorName := new(string)
	key := fmt.Sprintf("%s%d",viper.GetString("redis.KeyUserHashPrefix"), *authorId)
	fmt.Println(key)
	n, err := utils.RDB1.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if n == 0 {
		authorName, err = mysql.QueryAuthorName(authorId)
		if err != nil {
			return nil, err
		}
		err := Myredis.RedisSetHashRDB1(key, "username", *authorName)
		if err != nil {
			return nil, err
		}
	}
	*authorName, _ = utils.RDB1.HGet(ctx, key, "username").Result()
	// 没有username字段
	if *authorName == "" {
		authorName, err = mysql.QueryAuthorName(authorId)
		if err != nil {
			return nil, err
		}
		err := Myredis.RedisSetHashRDB1(key, "username", *authorName)
		if err != nil {
			return nil, err
		}
	}
	// 2. 创建作者对象
	author := Author{
		FollowCount: 0,
		FollowerCount: 0,
		IsFollow: false,
		Name: *authorName,
		Id: *authorId,
	}
	// 3. 获取视频信息
	videoListFromDao := &[]models.VideoBasic{}
	// 感觉缓存起来有点麻烦，先直接数据库查了，看看能不能实现基础功能再缓存优化
	// 要缓存的话key: userId, value: list[videoId, ...]，缓存位置在RDB4
	// 然后在得到的list中遍历是否有videoId的缓存
	// 但是这样的话下面的函数是不用写的，需要写的是查询某个用户的所有视频id，但你查了id，不如直接返回所有字段
	// 或者在返回的所有字段里，处理各个字段，达到上面的缓存效果
	videoListFromDao, err = mysql.QueryVideoList(authorId)
	if err != nil {
		return nil, err
	}
	// 4. 创建返回的视频列表参数
	videoList := &[]Video{}
	for i := range *videoListFromDao {
		videoId := (*videoListFromDao)[i].Identity 
		// 5. 获取赞数
		favoriteCount := new(int64)
		key := fmt.Sprintf("%s%d",viper.GetString("redis.KeyVideoFavoriteCountStringPrefix"), videoId)
		// 先从RDB0中查看键值对是否存在
		n, err := utils.RDB0.Exists(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		// 如果redis中没有key，调用mysql的方法获得状态
		if n == 0 {
			favoriteCount, err = mysql.QueryVideoFavoriteCount(&videoId)
			if err != nil {
				return nil, err
			}
			err = Myredis.RedisAddStringRDB0(key, fmt.Sprintf("%d", *favoriteCount))
			if err != nil {
				return nil, err
			}
		// redis中有key，获取string转换成int64
		} else {
			sFavoriteCount, err := utils.RDB0.Get(ctx, key).Result()
			if err != nil {
				return nil, err
			}
			iFavoriteCount, err := strconv.Atoi(sFavoriteCount)
			if err != nil {
				return nil, err
			}
			*favoriteCount = int64(iFavoriteCount)
		}
		// 6. 获取评论数
		commentCount := new(int64)
		key = fmt.Sprintf("%s%d",viper.GetString("redis.KeyVideoCommentCountStringPrefix"), videoId)
		// 先从RDB0中查看键值对是否存在
		n, err = utils.RDB0.Exists(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		// 如果redis中没有key，调用mysql的方法获得状态
		if n == 0 {
			commentCount, err = mysql.QueryCommentCount(&videoId)
			if err != nil {
				return nil, err
			}
			err = Myredis.RedisAddStringRDB0(key, fmt.Sprintf("%d", *commentCount))
			if err != nil {
				return nil, err
			}
		// redis中有key，获取string转换成int64
		} else {
			sCommentCount, err := utils.RDB0.Get(ctx, key).Result()
			if err != nil {
				return nil, err
			}
			iCommentCount, err := strconv.Atoi(sCommentCount)
			if err != nil {
				return nil, err
			}
			*commentCount = int64(iCommentCount)
		}
		// 7. 判断使用者是否喜欢该视频
		var isFavorite bool
		key = fmt.Sprintf("%s%d-%d",viper.GetString("redis.KeyUserLoveVideoStringPrefix"), videoId, *loginUserId)
		// 先从RDB0中查看键值对是否存在
		n, err = utils.RDB0.Exists(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		// 如果redis中没有key，调用mysql的方法获得状态
		if n == 0 {
			isFavorite, err = mysql.QueryIsFavorite(&videoId, loginUserId)
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
		// redis中有key，获取string转换成bool
		} else {
			isFavorite = true
			if utils.RDB0.Get(ctx, key).Val() == "0" {
				isFavorite = false
			}
		}
		// 8. 创建单个视频对象
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
