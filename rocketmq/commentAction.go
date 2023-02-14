package rocket

import (
	"fmt"
	myRedis "simple_tiktok/dao/redis"
	"simple_tiktok/logger"
	"simple_tiktok/models"
	"strconv"

	"github.com/spf13/viper"
)

func CommentAction(commentinfo models.CommentVideo) {
	var newCathe = map[string]interface{}{
		"id":           commentinfo.Identity,
		"video_id":     commentinfo.VideoIdentity,
		"author_id":    commentinfo.UserIdentity,
		"text":         commentinfo.Text,
		"comment_time": commentinfo.CommentTime,
	}
	hashKey := viper.GetString("redis.KeyCommentHashPrefix") + strconv.Itoa(int(commentinfo.Identity))
	err := myRedis.RedisAddCommentInfo(hashKey, newCathe)
	if err != nil {
		logger.SugarLogger.Error(err)
	}
	fmt.Println("消息队列处理成功")
}
