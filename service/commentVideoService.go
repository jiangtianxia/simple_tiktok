package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"simple_tiktok/dao/mysql"
	myredis "simple_tiktok/dao/redis"
	"simple_tiktok/logger"
	"simple_tiktok/utils"
	"strconv"
)

/**
 * @Author
 * @Description 视频流接口
 * @Date 21:00 2023/2/11
 **/
func CommentList(c *gin.Context, user_id uint64, video_id uint64) ([]CommentInfo, error) {

	listKey := viper.GetString("redis.KeyCommentListPrefix") + strconv.FormatUint(video_id, 10)
	//先判断存储视频评论id的是否存在
	if utils.RDB8.Exists(c, listKey).Val() == 0 {
		// 查询数据库 获得评论
		comments, err := mysql.QueryVideoCommentInfo(&video_id)
		if err != nil {
			logger.SugarLogger.Error("FindComment Error：", err.Error())
			return nil, err
		}
		for i := range *comments {
			// 缓存评论id
			err = myredis.RedisAddListRBD8(c, listKey, fmt.Sprintf("%d", (*comments)[i].Identity))
		}
	}

	// 使用缓存 查询出发布时间小于latestTime的30条记录  记录中包含视频的identity
	identityList := utils.RDB8.LRange(c, listKey, 0, -1).Val()

	commentInfos := make([]CommentInfo, len(identityList))
	for i, identity := range identityList {
		// 查询缓存 查询RDB7
		hashKey := viper.GetString("redis.KeyCommentInfoHashPrefix") + identity
		if utils.RDB7.Exists(c, hashKey).Val() == 0 {
			id, _ := strconv.Atoi(identity)
			commentID := uint64(id)
			comment, err := mysql.QueryCommentInfoByID(commentID)
			if err != nil {
				logger.SugarLogger.Error(err)
				return nil, err
			}

			err = myredis.RedisSetHashRDB7(c, hashKey, map[string]interface{}{
				"identity":       comment.Identity,
				"video_identity": comment.VideoIdentity,
				"user_identity":  comment.UserIdentity,
				"text":           comment.Text,
				"comment_time":   comment.CommentTime,
			})
			if err != nil {
				logger.SugarLogger.Error(err)
				return nil, err
			}
			// fmt.Println("数据库")
		}

		// 使用缓存获取信息
		comment := utils.RDB7.HGetAll(c, hashKey).Val()

		commentInfos[i].Id, _ = strconv.ParseUint(comment["identity"], 10, 64)
		userId, _ := strconv.ParseUint(comment["user_identity"], 10, 64)
		user, _ := mysql.FindUserByIdentity(userId)

		// 获取关注数
		followCount, err := GetFollowCount(c, strconv.Itoa(int(user.Identity)))
		if err != nil {
			logger.SugarLogger.Error("GetFollowCount Error：", err.Error())
			return nil, err
		}

		// 获取粉丝数
		followerCount, err := GetFollowerCount(c, strconv.Itoa(int(user.Identity)))
		if err != nil {
			logger.SugarLogger.Error("GetFollowerCount Error：", err.Error())
			return nil, err
		}

		// 判断是否关注该用户
		flag := false
		if user_id != 0 {
			if user.Identity == user_id {
				flag = true
			} else {
				flag, err = IsFollow(c, strconv.Itoa(int(user.Identity)), strconv.Itoa(int(user_id)))
				// fmt.Println(flag)
				if err != nil {
					logger.SugarLogger.Error("IsFollow Error：", err.Error())
					return nil, err
				}
			}
		}

		commentInfos[i].User = Author{
			Id:            user.Identity,
			Name:          user.Username,
			FollowCount:   followCount,
			FollowerCount: followerCount,
			IsFollow:      flag,
		}

		commentInfos[i].Content = comment["text"]
		commentInfos[i].CreateDate = comment["comment_time"]
	}

	return commentInfos, nil
}
