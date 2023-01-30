package service

import (
	"simple_tiktok/dao/mysql"
)

func GetVideoListByUserId(userName *string) (*[]mysql.Video, error){
	// 获取作者信息
	var author mysql.Author
	author.FollowCount = 0
	author.FollowerCount = 0
	author.IsFollow = false
	author.Name = *userName
	authorId, err := mysql.QueryAuthorInfo(userName)
	if err != nil {
		return nil, err
	}
	author.Id = *authorId
	// 获取视频信息
	videoList, err := mysql.QueryVideoList(userName)
	if err != nil {
		return nil, err
	}
	// 给视频添加作者信息
	for i := range *videoList {
		(*videoList)[i].Author = author
	}

	return videoList, nil
}
