package service

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/spf13/viper"
	"simple_tiktok/utils"
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
		// 获取视频作者信息
		// TODO
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

// 从缓存里取，如果没有就缓存进去
func getFavoriteVideoList(userId *uint64) (*[]uint64, error){
	key := fmt.Sprintf("%s%d", viper.GetString("redis.KeyUserFavoriteListPrefix"))
	n, err := utils.RDB6.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if n == 0 {
		//TODO
		return nil, nil
	}
	// TODO
	return nil, nil
}
