package service

import (
	"github.com/gin-gonic/gin"
	"simple_tiktok/dao/mysql"
	"simple_tiktok/logger"
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

	//根据传入的时间，获得传入时间前n个视频，可以通过config.videoCount来控制
	videos, err := mysql.FindVideosByLatestTime(latestTime)
	if err != nil {
		logger.SugarLogger.Error(err)
		return nil, time.Time{}, err
	}
	// 定义要返回的数据
	videoInfos := make([]VideoInfo, len(videos))
	for i, video := range videos {
		user, _ := mysql.FindUserByIdentity(video.UserIdentity)
		videoInfos[i].Identity = video.Identity
		videoInfos[i].Author = Author{
			Id:            user.Identity,
			Name:          user.Username,
			FollowCount:   0,
			FollowerCount: 0,
			IsFollow:      false,
		}
		videoInfos[i].PlayUrl = video.PlayUrl
		videoInfos[i].CoverUrl = video.CoverUrl
		videoInfos[i].CommentCount = 0
		videoInfos[i].IsFavorite = false
		videoInfos[i].Title = video.Title
	}

	return videoInfos, videos[len(videos)-1].PublishTime, nil
}
