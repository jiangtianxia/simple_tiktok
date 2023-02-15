package service

import (
	"github.com/gin-gonic/gin"
)

func GetFavoriteListByUserId(ctx *gin.Context, userId *uint64, loginUserId *uint64) (*[]VideoInfo, error) {
	// 获取用户喜欢的视频列表
	favoriteVideo, err := getFavoriteVideoList(ctx, userId)
	if err != nil {
		return nil, err
	}
	videoList := &[]VideoInfo{}
	// 遍历视频，调用之前写的函数获取视频的信息
	for i := range *favoriteVideo {
		// 获取视频信息
		videoId := (*favoriteVideo)[i]
		// 封面，播放地址，标题
		coverUrl, playUrl, title, err := tryToGetVideoInfo(ctx, &videoId)
		if err != nil {
			return nil, err
		}
		// 获赞数
		favoriteCount, err := getVideoFavoriteCount(ctx, videoId)
		if err != nil {
			return nil, err
		}
		// 评论数
		commentCount, err := getVideoCommentCount(ctx, videoId)
		if err != nil {
			return nil, err
		}
		// 用户是否喜欢
		isFavorite, err := judgeLoginUserLoveVideo(ctx, videoId, *loginUserId)
		if err != nil {
			return nil, err
		}
		// 获取视频作者信息
		authorId, err := getAuthorIdByVideoId(ctx, &videoId)
		authorName, err := getAuthorName(ctx, authorId)
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
		*videoList = append(*videoList, VideoInfo{
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
	return videoList, nil
}
