/*
 * @Description:
 * @Author: liuxin
 * @Date: 2023-01-28 09:33:15
 * @LastEditTime: 2023-01-29 13:45:24
 * @LastEditors:
 */
package service

import (
	"errors"
	"simple_tiktok/dao/mysql"
	"simple_tiktok/logger"
	"simple_tiktok/models"
	"simple_tiktok/utils"
)

type UserRegister struct {
	identity uint64
	username string
	password string

	token string
	data  *RegisterResponse
}

// 注册得到token和id
func PostUserRegister(username, password string) (*RegisterResponse, error) {
	return NewPostUserRegister(username, password).Solve()
}

func NewPostUserRegister(username, password string) *UserRegister {
	return &UserRegister{username: username, password: password}
}

func (req *UserRegister) Solve() (*RegisterResponse, error) {
	//验证
	if req.username == "" {
		return nil, errors.New("用户名为空")
	}
	if req.password == "" {
		return nil, errors.New("密码为空")
	}

	//更新数据
	if err := req.update(); err != nil {
		return nil, err
	}

	//response
	if err := req.registerResponse(); err != nil {
		return nil, err
	}
	return req.data, nil
}

func (req *UserRegister) update() error {
	getid, err := utils.GetID()
	if err != nil {
		return err
	}
	ur := models.UserBasic{Identity: getid, Username: utils.MakePassword(req.username), Password: utils.MakePassword(req.password)}

	//判断用户名
	ret1 := mysql.NewUserRegisterDAO()
	if ret1.IsExist(req.username) {
		return errors.New("用户名已存在")
	}

	//更新数据
	ret2 := mysql.NewUserBasicDAO()
	err = ret2.AddUserBasic(&ur)
	if err != nil {
		return err
	}

	//给token
	token, err3 := utils.GenerateToken(ur.Identity, ur.Username) //按格式填，这里还缺点东西，填userRegister
	if err3 != nil {
		logger.SugarLogger.Error("Generate Token Error:" + err.Error())
		return err3
	}
	req.token = token
	req.identity = ur.Identity
	return nil
}

func (req *UserRegister) registerResponse() error {
	req.data = &RegisterResponse{
		Identity: req.identity,
		Token:    req.token,
	}
	return nil
}

type RegisterResponse struct {
	Identity uint64 `json:"identity"`
	Token    string `json:"token"`
}
