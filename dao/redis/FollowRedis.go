package redis

import (
	"context"
	"simple_tiktok/utils"

	"github.com/spf13/viper"
)

/**
 * @Author jiang
 * @Description 删除用户关注的用户的列表
 * @Date 22:00 2023/2/11
 **/
func DeleteFollowList(userIdentity string) error {
	// 1、获取key
	key := viper.GetString("redis.KeyFollowListPrefix") + userIdentity

	// 2、删除
	return utils.RDB10.Del(context.Background(), key).Err()
}

/**
 * @Author jiang
 * @Description 删除粉丝sorted set
 * @Date 22:00 2023/2/11
 **/
func DeleteFollowerSortSet(userIdentity string) error {
	// 1、获取key
	key := viper.GetString("redis.KeyFollowerSortSetPrefix") + userIdentity

	// 2、删除
	return utils.RDB11.Del(context.Background(), key).Err()
}
