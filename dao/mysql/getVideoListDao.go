package mysql

import (
	//"gorm.io/gorm"
	"simple_tiktok/models"
)

func QueryVideoList(userId *uint64) (*[]models.VideoBasic, error) {
	return nil, nil
}

func QueryAuthorName(userId *uint64) (*string, error) {
	return nil, nil
}

func QueryVideoFavoriteCount (videoId *uint64) (*int, error) {
	return nil, nil
}

func QueryCommentFavoriteCount (videoId *uint64) (*int, error) {
	return nil, nil
}

func QueryIsFavorite(videoId *uint64, userId *uint64) (*bool, error) {
	return nil, nil
}
