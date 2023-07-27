package service

import (
	"errors"
	"fmt"
	"simple_tiktok/internal/comment"
	"simple_tiktok/internal/common"
	"simple_tiktok/internal/user"
	"simple_tiktok/model"
	"simple_tiktok/store"
	"simple_tiktok/utils"
	"time"
)

type commentService struct {
}

// CommentAction implements comment.ICommentService.
func (*commentService) CommentAction(req *comment.CommentActionReq) (*comment.Comment, error) {
	if req.ActionType == 2 {
		// 删除评论, 解析commentId后删除
		commentId, err := utils.DecodeID(req.CommentId)
		if err != nil {
			return nil, err
		}

		// 获取评论内容
		var c model.CommentVideo
		if err := store.GetDB().Model(&model.CommentVideo{}).
			Where("id = ? AND video_id = ?", commentId, req.VideoId).Scan(&c).Error; err != nil {
			return nil, err
		}

		resp := &comment.Comment{
			Id:         req.CommentId,
			User:       common.User{},
			Content:    c.Text,
			CreateDate: c.CommentTime,
		}
		userReq := &user.UserInfoReq{
			HashIds:   []string{},
			UserIds:   []uint{c.UserId},
			TokenInfo: req.TokenInfo,
		}
		userResp, err := UserService.GetUserInfo(userReq)
		if err != nil {
			return nil, err
		}
		for _, userInfo := range userResp {
			resp.User = userInfo
		}

		// 删除评论
		if err := store.GetDB().Delete(&c).Error; err != nil {
			return nil, err
		}
		return resp, nil
	}

	// 判断视频是否存在
	var cnt int64
	if err := store.GetDB().Model(model.VideoBasic{}).Where("id = ?", req.VideoId).Count(&cnt).Error; err != nil {
		return nil, err
	}
	if cnt == 0 {
		return nil, errors.New("视频错误")
	}

	// 初始化评论信息
	currentTime := time.Now()
	month := currentTime.Format("01")
	day := currentTime.Format("02")
	mmdd := fmt.Sprintf("%s-%s", month, day)
	c := &model.CommentVideo{
		VideoId:     req.VideoId,
		UserId:      req.TokenInfo.Id,
		Text:        req.CommentText,
		CommentTime: mmdd,
	}

	// 获取用户信息
	resp := &comment.Comment{
		Id:         "",
		User:       common.User{},
		Content:    c.Text,
		CreateDate: c.CommentTime,
	}
	userReq := &user.UserInfoReq{
		HashIds:   []string{},
		UserIds:   []uint{c.UserId},
		TokenInfo: req.TokenInfo,
	}
	userResp, err := UserService.GetUserInfo(userReq)
	if err != nil {
		return nil, err
	}
	for _, userInfo := range userResp {
		resp.User = userInfo
	}

	// 创建评论
	if err = store.GetDB().Create(c).Error; err != nil {
		return nil, err
	}
	hashId, _ := utils.EncodeID(c.ID)
	resp.Id = hashId
	return resp, nil
}

// CommentList implements comment.ICommentService.
func (*commentService) CommentList(req *comment.CommentListReq) (*comment.CommentListResp, error) {
	var total int64
	var commentList []model.CommentVideo
	offset := (req.Page - 1) * req.PageSize
	if err := store.GetDB().Model(&model.CommentVideo{}).Where("video_id = ?", req.VideoId).Count(&total).
		Order("comment_time DESC").
		Limit(int(req.PageSize)).Offset(int(offset)).Scan(&commentList).Error; err != nil {
		return nil, err
	}

	// 获取评论详细信息
	var commentListResp []comment.Comment
	userIds := map[string]string{}
	ids := []uint{}
	for _, c := range commentList {
		hashId, _ := utils.EncodeID(c.ID)
		tmp := comment.Comment{
			Id:         hashId,
			User:       common.User{},
			Content:    c.Text,
			CreateDate: c.CommentTime,
		}
		commentListResp = append(commentListResp, tmp)
		uid, _ := utils.EncodeID(c.UserId)
		userIds[hashId] = uid
		ids = append(ids, c.UserId)
	}

	// 获取用户信息
	userReq := &user.UserInfoReq{
		HashIds:   []string{},
		UserIds:   ids,
		TokenInfo: req.TokenInfo,
	}
	userResp, err := UserService.GetUserInfo(userReq)
	if err != nil {
		return nil, err
	}

	for i, v := range commentListResp {
		uid := userIds[v.Id]
		commentListResp[i].User = userResp[uid]
	}

	totalPages := total / req.PageSize
	if total%req.PageSize != 0 {
		totalPages++
	}
	return &comment.CommentListResp{
		PaginateResp: common.PaginateResp{
			Total:     total,
			Page:      req.Page,
			PageSize:  req.PageSize,
			TotalPage: totalPages,
		},
		CommentList: commentListResp,
	}, nil
}
