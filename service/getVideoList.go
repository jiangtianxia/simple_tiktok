package service

import (
	"simple_tiktok/dao/mysql"
)

func GetVideoListByUserId(authorId *uint64, loginUserId *uint64) (*[]Video, error){
	// 获取作者用户名
	authorName, err := mysql.QueryAuthorName(authorId)
	if err != nil {
		return nil, err
	}
	// 创建作者对象
	author := Author{
		FollowCount: 0,
		FollowerCount: 0,
		IsFollow: false,
		Name: *authorName,
		Id: *authorId,
	}
	// 获取视频信息
	videoListFromDao, err := mysql.QueryVideoList(authorId)
	if err != nil {
		return nil, err
	}
	// 创建返回的视频列表参数
	videoList := &[]Video{}
	for i := range *videoListFromDao {
		// 获取赞数
		favoriteCount, err := mysql.QueryVideoFavoriteCount(&(*videoListFromDao)[i].Identity)
		if err != nil {
			return nil, err
		}
		// 获取评论数
		commentCount, err := mysql.QueryCommentFavoriteCount(&(*videoListFromDao)[i].Identity)
		if err != nil {
			return nil, err
		}
		// 判断使用者是否喜欢该视频
		isFavorite, err := mysql.QueryIsFavorite(&(*videoListFromDao)[i].Identity, loginUserId)
		if err != nil {
			return nil, err
		}
		// 创建单个视频对象
		*videoList = append(*videoList, Video{
			Id: (*videoListFromDao)[i].Identity,
			Author: author,
			PlayUrl: (*videoListFromDao)[i].PlayUrl,
			CoverUrl: (*videoListFromDao)[i].CoverUrl,
			FavoriteCount: *favoriteCount,
			CommentCount: *commentCount,
			IsFavorite: *isFavorite,
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
	FavoriteCount int `json:"favorite_count"`
	CommentCount int `json:"comment_count"`
	IsFavorite bool `json:"is_favorite"`
	Title string `json:"title"`
}
// 作者参数
type Author struct {
	Id uint64 `json:"id"`
	Name string `json:"name"`
	FollowCount int `json:"follow_count"` // default
	FollowerCount int `json:"follower_count"` // default
	IsFollow bool `json:"is_follow"` // defalut
}
