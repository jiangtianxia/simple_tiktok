package service

import (
	"context"
	"fmt"
	"simple_tiktok/internal/common"
	"simple_tiktok/internal/favorite"
	"simple_tiktok/internal/video"
	"simple_tiktok/model"
	"simple_tiktok/store"
	"strconv"
	"time"
)

type favoriteService struct{}

// FavoriteAction implements favorite.IFavoriteService.
func (*favoriteService) FavoriteAction(req *favorite.FavoriteActionReq) error {
	if req.ActionType == 1 {
		return FavoriteVideo(req)
	}
	return NFavoriteVideo(req)
}

// 点赞赞视频
func FavoriteVideo(req *favorite.FavoriteActionReq) error {
	var cnt int64

	if err := store.GetDB().Model(&model.FavouriteVideo{}).
		Where("video_id = ? AND user_id = ?", req.VideoId, req.TokenInfo.Id).Count(&cnt).Error; err != nil {
		return err
	}

	tx := store.GetDB().Begin()
	if cnt == 0 {
		if err := tx.Create(&model.FavouriteVideo{
			VideoId: req.VideoId,
			UserId:  req.TokenInfo.Id,
			Status:  1,
		}).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		if err := tx.Model(&model.FavouriteVideo{}).
			Where("video_id = ? AND user_id = ?", req.VideoId, req.TokenInfo.Id).
			Update("status", 1).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 修改点赞集合
	key := fmt.Sprintf("%s%d", store.FavoriteSetPrefix, req.TokenInfo.Id)
	exists, err := store.GetRDB().Exists(context.Background(), key).Result()
	if err != nil {
		tx.Rollback()
		return err
	}
	if exists == 1 {
		_, err := store.GetRDB().SAdd(context.Background(), key, req.VideoId).Result()
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// 修改视频哈希
	videoKey := fmt.Sprintf("%s%d", store.VideoHashPrefix, req.VideoId)
	exists, err = store.GetRDB().Exists(context.Background(), videoKey).Result()
	if err != nil {
		tx.Rollback()
		return err
	}
	if exists == 1 {
		_, err := store.GetRDB().HIncrBy(context.Background(), videoKey, "favorite_count", 1).Result()
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// 修改用户哈希
	userKey := fmt.Sprintf("%s%d", store.UserHashPrefix, req.TokenInfo.Id)
	exists, err = store.GetRDB().Exists(context.Background(), userKey).Result()
	if err != nil {
		tx.Rollback()
		return err
	}
	if exists == 1 {
		_, err := store.GetRDB().HIncrBy(context.Background(), userKey, "favorite_count", 1).Result()
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

// 取消点赞
func NFavoriteVideo(req *favorite.FavoriteActionReq) error {
	var cnt int64

	if err := store.GetDB().Model(&model.FavouriteVideo{}).
		Where("video_id = ? AND user_id = ?", req.VideoId, req.TokenInfo.Id).Count(&cnt).Error; err != nil {
		return err
	}
	if cnt == 0 {
		return nil
	}

	tx := store.GetDB().Begin()
	if err := tx.Model(&model.FavouriteVideo{}).
		Where("video_id = ? AND user_id = ?", req.VideoId, req.TokenInfo.Id).
		Update("status", 0).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 修改点赞集合
	key := fmt.Sprintf("%s%d", store.FavoriteSetPrefix, req.TokenInfo.Id)
	exists, err := store.GetRDB().Exists(context.Background(), key).Result()
	if err != nil {
		tx.Rollback()
		return err
	}
	if exists == 1 {
		_, err := store.GetRDB().SRem(context.Background(), key, req.VideoId).Result()
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// 修改视频哈希
	videoKey := fmt.Sprintf("%s%d", store.VideoHashPrefix, req.VideoId)
	exists, err = store.GetRDB().Exists(context.Background(), videoKey).Result()
	if err != nil {
		tx.Rollback()
		return err
	}
	if exists == 1 {
		_, err := store.GetRDB().HIncrBy(context.Background(), videoKey, "favorite_count", -1).Result()
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// 修改用户哈希
	userKey := fmt.Sprintf("%s%d", store.UserHashPrefix, req.TokenInfo.Id)
	exists, err = store.GetRDB().Exists(context.Background(), userKey).Result()
	if err != nil {
		tx.Rollback()
		return err
	}
	if exists == 1 {
		_, err := store.GetRDB().HIncrBy(context.Background(), userKey, "favorite_count", -1).Result()
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

// FavoriteList implements favorite.IFavoriteService.
func (*favoriteService) FavoriteList(req *favorite.FavoriteListReq) (map[string]common.Video, error) {
	// 查询redis
	key := fmt.Sprintf("%s%d", store.FavoriteSetPrefix, req.UserId)
	exists, err := store.GetRDB().Exists(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}
	if exists == 1 {
		setMembers, err := store.GetRDB().SMembers(context.Background(), key).Result()
		if err != nil {
			return nil, err
		}

		memberId := []uint{}
		for _, member := range setMembers {
			member, _ := strconv.ParseInt(member, 10, 64)
			memberId = append(memberId, uint(member))
		}

		videoReq := &video.VideoInfoReq{
			HashIds:  []string{},
			VideoIds: memberId,
			TokenInfo: common.TokenInfoReq{
				Id:       req.UserId,
				Username: "test",
			},
		}

		return VideoService.GetVideoInfo(videoReq)
	}

	// 查询数据库
	var list []model.FavouriteVideo
	if err := store.GetDB().Model(&model.FavouriteVideo{}).Where("user_id = ? AND status = 1", req.UserId).
		Scan(&list).Error; err != nil {
		return nil, err
	}

	// 查询视频信息
	pipeline := store.GetRDB().Pipeline()
	ids := []uint{}
	for _, item := range list {
		pipeline.SAdd(context.Background(), key, item.VideoId)
		ids = append(ids, item.VideoId)
	}
	pipeline.Expire(context.Background(), key, 7*24*time.Hour)
	_, err = pipeline.Exec(context.Background())
	if err != nil {
		return nil, err
	}

	videoReq := &video.VideoInfoReq{
		HashIds:  []string{},
		VideoIds: ids,
		TokenInfo: common.TokenInfoReq{
			Id:       req.UserId,
			Username: "test",
		},
	}
	return VideoService.GetVideoInfo(videoReq)
}
