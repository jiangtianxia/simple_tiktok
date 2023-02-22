package service

import (
	"encoding/json"
	"simple_tiktok/dao/mysql"
	"simple_tiktok/logger"
	"simple_tiktok/models"
	"simple_tiktok/utils"
	"strconv"

	"github.com/spf13/viper"
)

// 接收参数结构体
type FavouriteReqStruct struct {
	UserId     uint64
	VideoId    string
	ActionType string
}

/**
 * @Author
 * @Description 赞接口
 * @Date 12:00 2023/2/12
 **/
func DealFavourite(msgid string, data []byte) {
	FavouriteInfo := &FavouriteReqStruct{}
	json.Unmarshal(data, FavouriteInfo)

	flag := 1
	if FavouriteInfo.ActionType == "2" {
		flag = 0
	}

	// 判断数据库中是否存在该数据
	cnt, err := mysql.IsExistsFavouriteVideoCount(FavouriteInfo.VideoId, FavouriteInfo.UserId)
	if err != nil {
		logger.SugarLogger.Error("IsExistsFavouriteVideoCount Error：" + err.Error())
		SaveRedisResp(msgid, -1, "操作失败")
	}

	if cnt > 0 {
		//修改数据库
		err := mysql.UpdateFavourite(FavouriteInfo.VideoId, FavouriteInfo.UserId, flag)
		if err != nil {
			logger.SugarLogger.Error("修改点赞失败" + err.Error())
			SaveRedisResp(msgid, -1, "操作失败")
		}
	} else {
		// 创建数据库
		videoId, _ := strconv.Atoi(FavouriteInfo.VideoId)
		info := models.FavouriteVideo{
			VideoIdentity: uint64(videoId),
			UserIdentity:  FavouriteInfo.UserId,
			Status:        flag,
		}
		err := mysql.CreateFavouriteVideo(info)
		if err != nil {
			logger.SugarLogger.Error("创建点赞失败" + err.Error())
			SaveRedisResp(msgid, -1, "操作失败")
		}
	}

	//发送延迟消息，删除缓存
	RetryTopic := viper.GetString("rocketmq.RetryTopic")
	DeleteFavouriteRedisTag := viper.GetString("rocketmq.DeleteFavouriteRedisTag")
	err = utils.SendDelayMsg(RetryTopic, DeleteFavouriteRedisTag, data)
	if err != nil {
		logger.SugarLogger.Error("SendDelayMsg Error：", err.Error())
	}
	SaveRedisResp(msgid, 0, "操作成功")
}
