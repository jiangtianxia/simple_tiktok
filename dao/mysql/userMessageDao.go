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
func QueryMessageByIdentity(identity uint64) ([]models.UserMessage, error) {
	var messageList []models.UserMessage
	if err := utils.DB.Table("user_message").Where("identity = ?", identity).Find(&messageList).Error; err != nil {
		logger.SugarLogger.Error(err)
		return messageList, err
	}
	return messageList, nil
}
