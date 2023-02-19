package service

import (
	"fmt"
	"simple_tiktok/dao/mysql"
	Myredis "simple_tiktok/dao/redis"
	"simple_tiktok/logger"
	"simple_tiktok/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

/**
 * 获取视频列表的service层
 * @Author pjh
 * @Summary
 * @Tags
 **/

func GetVideoListByUserId(ctx *gin.Context, authorId *uint64, loginUserId *uint64) (*[]VideoInfo, error) {
	// 1. 获取作者用户名
	authorName, err := getAuthorName(ctx, authorId)
	if err != nil {
		return nil, err
	}

	// 获取关注数
	followCount, err := GetFollowCount(ctx, strconv.Itoa(int(*authorId)))
	if err != nil {
		logger.SugarLogger.Error("GetFollowCount Error：", err.Error())
		return nil, err
	}

	// 获取粉丝数
	followerCount, err := GetFollowerCount(ctx, strconv.Itoa(int(*authorId)))
	if err != nil {
		logger.SugarLogger.Error("GetFollowerCount Error：", err.Error())
		return nil, err
	}

	// 判断是否关注该用户
	isFollow, err := IsFollow(ctx, strconv.Itoa(int(*authorId)), strconv.Itoa(int(*loginUserId)))
	if err != nil {
		logger.SugarLogger.Error("IsFollow Error：", err.Error())
		return nil, err
	}

	// 获取点赞数量，作品数和喜欢数
	totalFavourited, workCount, FavouriteCount, err := GetTotalFavouritedANDWorkCountANDFavoriteCount(*authorId)
	if err != nil {
		logger.SugarLogger.Error("GetTotalFavouritedANDWorkCountANDFavoriteCount Error：", err.Error())
		return nil, err
	}

	// 2. 创建作者对象
	author := Author{
		FollowCount:     followCount,
		FollowerCount:   followerCount,
		IsFollow:        isFollow,
		Name:            *authorName,
		Id:              *authorId,
		Avatar:          viper.GetString("defaultAvatarUrl"),
		BackGroundImage: viper.GetString("defaultBackGroudImage"),
		Signature:       viper.GetString("defaultSignature"),
		TotalFavorited:  totalFavourited,
		WorkCount:       workCount,
		FavoriteCount:   FavouriteCount,
	}

	// 3. 尝试从缓存中获取用户发表的视频id列表
	key := fmt.Sprintf("%s%d", viper.GetString("redis.KeyPublishListPrefix"), *authorId)
	n, err := utils.RDB4.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	// 若缓存中存在信息则使用 A 计划
	// A4. 创建返回的视频列表参数
	if n != 0 {
		videoIds, err := utils.RDB4.LRange(ctx, key, 0, -1).Result()
		if err != nil {
			return nil, err
		}
		videoList := &[]VideoInfo{}
		for i := range videoIds {
			sVideoId := videoIds[i]
			iVideoId, err := strconv.Atoi(sVideoId)
			if err != nil {
				return nil, err
			}
			videoId := uint64(iVideoId)

			// A5. 获取赞数
			favoriteCount, err := getVideoFavoriteCount(ctx, videoId)
			if err != nil {
				return nil, err
			}

			// A6. 获取评论数
			commentCount, err := getVideoCommentCount(ctx, videoId)
			if err != nil {
				return nil, err
			}

			// A7. 判断使用者是否喜欢该视频
			bIsFavorite, err := judgeLoginUserLoveVideo(ctx, videoId, *loginUserId)
			isFavorite := *bIsFavorite
			if err != nil {
				return nil, err
			}

			// A8. 获取视频封面url和视频url
			coverUrl, playUrl, title, err := tryToGetVideoInfo(ctx, &videoId)
			if err != nil {
				return nil, err
			}

			// A9. 创建单个视频对象
			*videoList = append(*videoList, VideoInfo{
				Id:            videoId,
				Author:        author,
				PlayUrl:       viper.GetString("cos.addr") + *playUrl,
				CoverUrl:      viper.GetString("uploadAddr") + *coverUrl,
				FavoriteCount: *favoriteCount,
				CommentCount:  *commentCount,
				IsFavorite:    isFavorite,
				Title:         *title,
			})
		}
		return videoList, nil

	}

	// 以下为 B 计划，需要对 用户发表的视频id列表 和 视频信息 进行缓存
	// 从数据库中获取视频信息
	videoListFromDao, err := mysql.QueryVideoList(authorId)
	if err != nil {
		return nil, err
	}

	// B4. 创建返回的视频列表参数
	videoList := &[]VideoInfo{}
	for i := range *videoListFromDao {
		videoId := (*videoListFromDao)[i].Identity

		// B5. 获取赞数
		favoriteCount, err := getVideoFavoriteCount(ctx, videoId)
		if err != nil {
			return nil, err
		}

		// B6. 获取评论数
		commentCount, err := getVideoCommentCount(ctx, videoId)
		if err != nil {
			return nil, err
		}

		// B7. 判断使用者是否喜欢该视频
		bIsFavorite, err := judgeLoginUserLoveVideo(ctx, videoId, *loginUserId)
		isFavorite := *bIsFavorite
		if err != nil {
			return nil, err
		}

		// B8. 对视频信息进行缓存
		// 缓存用户发布的视频列表
		err = Myredis.RedisAddPublishList(ctx, fmt.Sprintf("%d", *authorId), fmt.Sprintf("%d", videoId))
		if err != nil {
			return nil, err
		}
		// 缓存视频信息
		err = Myredis.RedisAddVideoInfo(ctx, (*videoListFromDao)[i])
		if err != nil {
			return nil, err
		}
		// B9. 创建单个视频对象
		*videoList = append(*videoList, VideoInfo{
			Id:            videoId,
			Author:        author,
			PlayUrl:       viper.GetString("cos.addr") + (*videoListFromDao)[i].PlayUrl,
			CoverUrl:      viper.GetString("uploadAddr") + (*videoListFromDao)[i].CoverUrl,
			FavoriteCount: *favoriteCount,
			CommentCount:  *commentCount,
			IsFavorite:    isFavorite,
		})
	}

	return videoList, nil
}
