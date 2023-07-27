package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"simple_tiktok/conf"
	"simple_tiktok/internal/common"
	"simple_tiktok/internal/user"
	"simple_tiktok/internal/video"
	"simple_tiktok/model"
	"simple_tiktok/store"
	"simple_tiktok/utils"
	"strconv"
	"strings"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

type videoService struct {
}

// GetVideoInfo implements video.IVideoService.
func (*videoService) GetVideoInfo(req *video.VideoInfoReq) (map[string]common.Video, error) {
	// 查询redis
	exists, result, err := store.GetHashData(req.VideoIds, store.VideoHashPrefix)
	if err != nil {
		return nil, err
	}

	respList := make([]common.Video, len(req.VideoIds))
	videoIds := []uint{}
	authorIds := []uint{}
	for i, videoId := range req.VideoIds {
		if !exists[i] {
			videoIds = append(videoIds, videoId)
			continue
		}

		hashId, err := utils.EncodeID(videoId)
		if err != nil {
			return nil, err
		}
		commentCount, _ := strconv.ParseInt(result[i]["comment_count"], 10, 64)
		favoriteCount, _ := strconv.ParseInt(result[i]["favorite_count"], 10, 64)
		isFavorite := false
		if req.TokenInfo.Id > 0 {
			var cnt int64
			if err := store.GetDB().Model(&model.FavouriteVideo{}).
				Where("video_id = ? AND user_id = ? AND status = 1", videoId, req.TokenInfo.Id).Count(&cnt).
				Error; err != nil {
				return nil, err
			}
			if cnt > 0 {
				isFavorite = true
			}
		}
		authorId, _ := strconv.ParseInt(result[i]["author_id"], 10, 64)
		authorIds = append(authorIds, uint(authorId))
		authorHashId, err := utils.EncodeID(uint(authorId))
		if err != nil {
			return nil, err
		}
		respList = append(respList, common.Video{
			Id:            hashId,
			PlayUrl:       result[i]["play_url"],
			CoverUrl:      result[i]["cover_url"],
			Title:         result[i]["title"],
			FavoriteCount: favoriteCount,
			CommentCount:  commentCount,
			IsFavorite:    isFavorite,
			Author: common.User{
				Id: authorHashId,
			},
		})
	}

	// 查询video mysql
	if len(videoIds) > 0 {
		var videoList []model.VideoBasic
		if err := store.GetDB().Model(model.VideoBasic{}).Where("id IN (?)", videoIds).Scan(&videoList).Error; err != nil {
			return nil, err
		}

		redisHashData := make(map[uint]map[string]interface{}, len(videoIds))
		for _, video := range videoList {
			hashId, err := utils.EncodeID(video.ID)
			if err != nil {
				return nil, err
			}
			playUrl := conf.GetGlobalConf().COSAddr + video.PlayUrl
			coverUrl := conf.GetGlobalConf().UploadAddr + video.CoverUrl
			title := video.Title
			var commentCount, favoriteCount int64
			if err := store.GetDB().Model(&model.CommentVideo{}).Where("video_id = ?", video.ID).
				Scan(&commentCount).Error; err != nil {
				return nil, err
			}
			if err := store.GetDB().Model(&model.FavouriteVideo{}).
				Where("video_id = ? AND status = 1", video.ID).Count(&favoriteCount).
				Error; err != nil {
				return nil, err
			}
			isFavorite := false
			if req.TokenInfo.Id > 0 {
				var cnt int64
				if err := store.GetDB().Model(&model.FavouriteVideo{}).
					Where("video_id = ? AND user_id = ? AND status = 1", video.ID, req.TokenInfo.Id).Count(&cnt).
					Error; err != nil {
					return nil, err
				}
				if cnt > 0 {
					isFavorite = true
				}
			}
			authorIds = append(authorIds, video.AuthorId)
			authorHashId, err := utils.EncodeID(video.AuthorId)
			if err != nil {
				return nil, err
			}
			redisHashData[video.ID] = store.NewRedisVideoHash(title, playUrl, coverUrl, video.AuthorId, commentCount, favoriteCount)
			respList = append(respList, common.Video{
				Id:            hashId,
				PlayUrl:       playUrl,
				CoverUrl:      coverUrl,
				Title:         title,
				FavoriteCount: favoriteCount,
				CommentCount:  commentCount,
				IsFavorite:    isFavorite,
				Author: common.User{
					Id: authorHashId,
				},
			})
		}
		store.CreateHashData(redisHashData, store.VideoHashPrefix)
	}

	// 获取用户信息
	userReq := &user.UserInfoReq{
		HashIds: []string{},
		UserIds: authorIds,
		TokenInfo: common.TokenInfoReq{
			Id:       req.TokenInfo.Id,
			Username: req.TokenInfo.Username,
		},
	}
	userResp, err := UserService.GetUserInfo(userReq)
	if err != nil {
		return nil, err
	}

	resp := make(map[string]common.Video)
	for _, video := range respList {
		if userInfo, ok := userResp[video.Author.Id]; ok {
			video.Author.Name = userInfo.Name
			video.Author.FavoriteCount = userInfo.FavoriteCount
			video.Author.WorkCount = userInfo.WorkCount
			resp[video.Id] = video
		}
	}
	return resp, nil
}

// GetVideoPublishList implements video.IVideoService.
func (*videoService) GetVideoPublishList(req *video.VideoPublishListReq) (*common.VideoListResp, error) {
	var total int64
	var videoList []model.VideoBasic
	offset := (req.Page - 1) * req.PageSize
	query := store.GetDB().Model(&model.VideoBasic{}).Where("author_id = ?", req.UserId)
	if req.Where != "" {
		w := strings.Split(req.Where, ",")
		for _, v := range w {
			k := strings.Split(v, ":")
			query = query.Where(fmt.Sprintf("%s = ?", k[0]), k[1])
		}
	}
	t := ""
	sort := "DESC"
	if req.Sort == 1 {
		sort = "ASC"
	}
	if req.Order != "" {
		t = req.Order + sort
	} else {
		t = "publish_time " + sort
	}
	if err := query.Order(t).Count(&total).Limit(int(req.PageSize)).Offset(int(offset)).Scan(&videoList).Error; err != nil {
		return nil, err
	}

	if len(videoList) == 0 {
		return nil, nil
	}

	totalPages := total / req.PageSize
	if total%req.PageSize != 0 {
		totalPages++
	}
	resp := &common.VideoListResp{
		PaginateResp: common.PaginateResp{
			Total:     total,
			Page:      req.Page,
			PageSize:  req.PageSize,
			TotalPage: totalPages,
		},
		VideoList: []common.Video{},
	}

	ids := []uint{}
	for _, video := range videoList {
		ids = append(ids, video.ID)
	}
	videoReq := &video.VideoInfoReq{
		HashIds:   []string{},
		VideoIds:  ids,
		TokenInfo: req.TokenInfo,
	}
	videoResp, err := VideoService.GetVideoInfo(videoReq)
	if err != nil {
		return nil, err
	}

	for _, video := range videoList {
		hashId, err := utils.EncodeID(video.ID)
		if err != nil {
			return nil, err
		}

		videoInfo := videoResp[hashId]
		resp.VideoList = append(resp.VideoList, videoInfo)
	}
	return resp, nil
}

// VideoFeed implements video.IVideoService.
func (*videoService) VideoFeed(req *video.VideoFeedReq) (*video.VideoFeedResp, error) {
	// 获取视频列表
	var videoList []model.VideoBasic
	if err := store.GetDB().Model(&model.VideoBasic{}).Where("publish_time < ?", req.LatestTime).
		Limit(30).Scan(&videoList).Error; err != nil {
		return nil, err
	}
	resp := &video.VideoFeedResp{
		VideoList: []common.Video{},
		NextTime:  time.Now().Unix(),
	}
	var Ids []uint
	for _, video := range videoList {
		resp.NextTime = utils.Min[int64](resp.NextTime, video.PublishTime)
		Ids = append(Ids, video.ID)
	}

	// 获取视频信息
	videoReq := &video.VideoInfoReq{
		HashIds:   []string{},
		VideoIds:  Ids,
		TokenInfo: req.TokenInfo,
	}
	videoResp, err := VideoService.GetVideoInfo(videoReq)
	if err != nil {
		return nil, err
	}
	for _, video := range videoResp {
		resp.VideoList = append(resp.VideoList, video)
	}
	return resp, nil
}

// VideoPublishAction implements video.IVideoService.
func (*videoService) VideoPublishAction(req *video.VideoPublishActionReq) error {
	// 拼接文件名
	suffix := path.Ext(req.FileHead.Filename)
	identity, err := utils.GetID()
	if err != nil {
		return err
	}
	filename := strconv.Itoa(int(time.Now().Unix())) + strconv.Itoa(int(identity))
	path := filename + suffix
	hashId, _ := utils.EncodeID(req.TokenInfo.Id)
	dir := conf.GetGlobalConf().UploadBase + hashId
	if err = os.MkdirAll(dir, 0777); err != nil {
		return err
	}
	coverPath := dir + "/" + "cover" + filename + ".png"

	// 操作Mysql数据库
	tx := store.GetDB().Begin()
	coverUrl := coverPath[1:]
	videoInfo := &model.VideoBasic{
		AuthorId:    req.TokenInfo.Id,
		PlayUrl:     "/" + path,
		CoverUrl:    coverUrl,
		Title:       req.Title,
		PublishTime: time.Now().Unix(),
	}

	if err := tx.Create(videoInfo).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 保存到COS
	if err := saveCOS(req.File, path, coverPath); err != nil {
		tx.Rollback()
		return err
	}

	// 增加redis
	userInfoKey := fmt.Sprintf("%s%d", store.UserHashPrefix, req.TokenInfo.Id)
	cnt, err := store.GetRDB().Exists(context.Background(), userInfoKey).Result()
	if err != nil {
		tx.Rollback()
		return err
	}
	if cnt != 0 {
		_, err := store.GetRDB().HIncrBy(context.Background(), userInfoKey, "work_count", 1).Result()
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func saveCOS(file *multipart.File, path string, coverPath string) error {
	// 保存到COS
	_, err := store.GetCos().Object.Put(context.Background(), path, *file, nil)
	if err != nil {
		return err
	}
	fd, err := os.Create(coverPath)
	if err != nil {
		return err
	}

	// 读取COS的封面信息，保存到本地
	opt := &cos.GetSnapshotOptions{
		Time:   1,
		Format: "png",
	}
	resp, err := store.GetCos().CI.GetSnapshot(context.Background(), path, opt)
	if err != nil {
		return err
	}
	_, err = io.Copy(fd, resp.Body)
	if err != nil {
		return err
	}
	fd.Close()
	return nil
}
