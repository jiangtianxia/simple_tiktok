package service

import (
	"errors"
	"simple_tiktok/internal/common"
	"simple_tiktok/internal/user"
	"simple_tiktok/model"
	"simple_tiktok/store"
	"simple_tiktok/utils"
	"strconv"
)

type userService struct {
}

// UserRegister implements user.IUserService.
func (*userService) UserRegister(req *user.NormalizeUserReq) (*user.UserRegisterResp, error) {
	// 判断用户名是否存在
	var cnt int64
	if err := store.GetDB().Model(&model.UserBasic{}).Where("username = ?", req.Username).Count(&cnt).Error; err != nil {
		return nil, err
	}
	if cnt != 0 {
		return nil, errors.New("用户名已存在")
	}

	// 用户注册
	password := utils.MakePassword(req.Password)
	u := model.UserBasic{
		Username: req.Username,
		Password: password,
	}

	tx := store.GetDB().Begin()
	if err := tx.Create(&u).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 生成token
	token, err := utils.GenerateToken(u.ID, u.Username)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	hashId, err := utils.EncodeID(u.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return &user.UserRegisterResp{
		HashId: hashId,
		Token:  token,
	}, nil
}

// UserLogin implements user.IUserService.
func (*userService) UserLogin(req *user.NormalizeUserReq) (*user.UserLoginResp, error) {
	// 判断用户名是否存在
	var u *model.UserBasic
	if err := store.GetDB().Model(&model.UserBasic{}).Where("username = ?", req.Username).Scan(&u).Error; err != nil {
		return nil, err
	}
	if u == nil {
		return nil, errors.New("用户名不已存在")
	}

	// 验证密码
	if flag := utils.ValidPassword(req.Password, u.Password); !flag {
		return nil, errors.New("密码错误")
	}

	// 生成token
	token, err := utils.GenerateToken(u.ID, u.Username)
	if err != nil {
		return nil, err
	}

	hashId, err := utils.EncodeID(u.ID)
	if err != nil {
		return nil, err
	}

	return &user.UserLoginResp{
		HashId: hashId,
		Token:  token,
	}, nil
}

// GetUserInfo implements user.IUserService.
func (*userService) GetUserInfo(req *user.UserInfoReq) (map[string]common.User, error) {
	// 查询redis
	exists, result, err := store.GetHashData(req.UserIds, store.UserHashPrefix)
	if err != nil {
		return nil, err
	}
	resp := map[string]common.User{}
	userIds := []uint{}
	for i, userId := range req.UserIds {
		if !exists[i] {
			userIds = append(userIds, userId)
			continue
		}

		hashId, err := utils.EncodeID(userId)
		if err != nil {
			return nil, err
		}
		workCount, _ := strconv.ParseInt(result[i]["work_count"], 10, 64)
		favoriteCount, _ := strconv.ParseInt(result[i]["favorite_count"], 10, 64)
		resp[hashId] = common.User{
			Id:            hashId,
			Name:          result[i]["username"],
			WorkCount:     workCount,
			FavoriteCount: favoriteCount,
		}
	}

	// 不存在的查询mysql, 并创建redis数据
	if len(userIds) != 0 {
		var userList []model.UserBasic
		if err := store.GetDB().Model(model.UserBasic{}).Where("id IN (?)", userIds).Scan(&userList).Error; err != nil {
			return nil, err
		}

		redisHashData := make(map[uint]map[string]interface{}, len(userIds))
		for _, user := range userList {
			hashId, err := utils.EncodeID(user.ID)
			if err != nil {
				return nil, err
			}
			var favoriteCount int64
			if err := store.GetDB().Model(&model.FavouriteVideo{}).Where("user_id = ? AND status = 1", user.ID).Count(&favoriteCount).Error; err != nil {
				return nil, err
			}
			var workCount int64
			if err := store.GetDB().Model(&model.VideoBasic{}).Where("author_id = ?", user.ID).Count(&workCount).Error; err != nil {
				return nil, err
			}
			redisHashData[user.ID] = store.NewRedisUserHash(user.Username, workCount, favoriteCount)
			resp[hashId] = common.User{
				Id:            hashId,
				Name:          user.Username,
				WorkCount:     workCount,
				FavoriteCount: favoriteCount,
			}
		}
		store.CreateHashData(redisHashData, store.UserHashPrefix)
	}
	return resp, nil
}
