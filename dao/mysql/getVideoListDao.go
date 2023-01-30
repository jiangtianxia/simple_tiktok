package mysql

import (
	//"gorm.io/gorm"
)

// 视频参数
type Video struct {
	Id int `json:"id"`
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
	Id int `json:"id"`
	Name string `json:"name"`
	FollowCount int `json:"follow_count"` // default
	FollowerCount int `json:"follower_count"` // default
	IsFollow bool `json:"is_follow"` // defalut
}

func QueryVideoList(userName *string) (*[]Video, error) {
	return nil, nil
}

func QueryAuthorInfo(userName *string) (*int, error) {
	return nil, nil
}
