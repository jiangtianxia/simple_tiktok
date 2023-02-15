package mysql

import (
	"simple_tiktok/logger"
	"simple_tiktok/models"
	"simple_tiktok/utils"
)

// 新增聊天信息
func CreateUserMessage(userMessage models.UserMessage) error {
	return utils.DB.Create(&userMessage).Error
}

// 根据接收者identity查询聊天信息
func QueryMessageByToUserIdentity(fromUserIdentity, toUserIdentity uint64) ([]models.UserMessage, error) {
	var messageList []models.UserMessage
	if err := utils.DB.Table("user_message").Where("from_user_identity = ? And to_user_identity = ?", fromUserIdentity, toUserIdentity).Find(&messageList).Error; err != nil {
		logger.SugarLogger.Error(err)
		return nil, err
	}
	return messageList, nil
}

// 根据消息Identity查询聊天信息
func QueryMessageByIdentity(identity uint64) (models.UserMessage, error) {
	var message models.UserMessage
	if err := utils.DB.Table("user_message").Where("identity = ?", identity).First(&message).Error; err != nil {
		logger.SugarLogger.Error(err)
		return message, err
	}
	return message, nil
}

// 查询最新一条消息
func QueryNewMessage(fromUserIdentity, toUserIdentity string) (models.UserMessage, error) {
	var messageInfo models.UserMessage
	if err := utils.DB.Table("user_message").Where("from_user_identity = ? And to_user_identity = ?", fromUserIdentity, toUserIdentity).Order("create_time desc").First(&messageInfo).Error; err != nil {
		logger.SugarLogger.Error(err)
		return messageInfo, err
	}
	return messageInfo, nil
}
