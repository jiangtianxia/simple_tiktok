package service

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/spf13/viper"
	"simple_tiktok/utils"
	"strconv"
	"simple_tiktok/dao/mysql"
	Myredis "simple_tiktok/dao/redis"
)

func GetFavoriteListByUserId(c *gin.Context, userId *uint64, loginUserId *uint64) (*[]Video, error) {
	if ctx == nil {
		ctx = c
	}
	// 获取用户喜欢的视频列表
	favoriteVideo, err := getFavoriteVideoList(userId)
	if err != nil {
		return nil, err
	}
	videoList := &[]Video{}
	// 遍历视频，调用之前写的函数获取视频的信息
	for i := range *favoriteVideo {
		// 获取视频信息
		videoId := (*favoriteVideo)[i]
		// 封面，播放地址，标题
		coverUrl, playUrl, title, err := tryToGetVideoInfo(&videoId)
		if err != nil {
			return nil, err
		}
		// 获赞数
		favoriteCount, err := getVideoFavoriteCount(videoId)
		if err != nil {
			return nil, err
		}
		// 评论数
		commentCount, err := getVideoCommentCount(videoId)
		if err != nil {
			return nil, err
		}
		// 用户是否喜欢
		isFavorite, err := judgeLoginUserLoveVideo(videoId, *loginUserId)
		if err != nil {
			return nil, err
		}
		// 获取视频作者信息
		// TODO
		authorId, err := getAuthorIdByVideoId(&videoId)
		authorName, err := getAuthorName(authorId)
		if err != nil {
			return nil, err
		}
		author := Author{
			FollowCount: 0,
			FollowerCount: 0,
			IsFollow: false,
			Name: *authorName,
			Id: *authorId,
		}
		// 添加进视频列表
		*videoList = append(*videoList, Video{
			Id: videoId,
			Author: author,
			PlayUrl: *playUrl,
			CoverUrl: *coverUrl,
			FavoriteCount: *favoriteCount,
			CommentCount: *commentCount,
			IsFavorite: *isFavorite,
			Title: *title,
		})
	}
	
	// 返回给controller层
	return nil, nil
}

// 从缓存里取指定用户喜欢的视频id，如果没有就缓存进去
func getFavoriteVideoList(userId *uint64) (*[]uint64, error){
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
			if (*favoriteList)[i].Status == "0" {
				continue
			}
			videoList = append(videoList, (*favoriteList)[i].VideoIdentity)
			err := Myredis.RedisAddListRDB6(key, fmt.Sprintf("%d", (*favoriteList)[i].VideoIdentity))
			if err != nil {
				return nil, err
			}
		}
		return &videoList, nil
	}
	res, err := utils.RDB6.LRange(ctx, key, 0, -1).Result()
	for i := range res {
		id, err := strconv.Atoi(res[i])
		if err != nil {
			return nil, err
		}
		videoList = append(videoList, uint64(id))
	}
	return &videoList, nil
}

// 先从视频信息的缓存中找，没有就把视频信息加到缓存中去
func getAuthorIdByVideoId(videoId *uint64) (*uint64, error) {
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
		err = Myredis.RedisAddVideoInfo(*videoBasic)
		if err != nil {
			return nil, err
		}
		return &videoBasic.UserIdentity, nil
	}
	sId, err := utils.RDB3.HGet(ctx, key, "user_identity").Result()
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
