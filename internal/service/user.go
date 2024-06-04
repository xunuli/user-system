package service

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"user_system/internal/cache"
	"user_system/internal/dao"
	"user_system/internal/model"
	"user_system/pkg/constant"
	"user_system/utils"
)

//具体的用户服务

// Register 用户注册
func Register(req *RegisterRequest) error {
	//校验请求的格式是否正确
	if req.UserName == "" || req.Password == "" || req.Age <= 0 || !utils.Contains([]string{constant.GenderMale, constant.GenderFeMale}, req.Gender) {
		log.Errorf("register param invalid")
		return fmt.Errorf("register param invalid")
	}
	//查看是否用户已经存在
	existedUser, err := dao.GetUserByName(req.UserName)
	if err != nil {
		log.Errorf("Register|%v", err)
		return fmt.Errorf("Register|%v", err)
	}
	//用户已经存在
	if existedUser != nil {
		log.Errorf("用户已经注册，user_name==%s", req.UserName)
		return fmt.Errorf("用户已经注册，不能重复注册")
	}
	//创建一个用户的记录
	user := &model.User{
		Name:     req.UserName,
		Age:      req.Age,
		Gender:   req.Gender,
		Password: req.Password,
		Nickname: req.NickName,
		CreateModel: model.CreateModel{
			Creator: req.UserName,
		},
		ModifyModel: model.ModifyModel{
			Modifier: req.UserName,
		},
	}

	log.Infof("user ==== %+v", user)
	if err := dao.CreateUser(user); err != nil {
		log.Errorf("Register|%v", err)
		return fmt.Errorf("Register|%v", err)
	}
	return nil
}

// 用户登录
func Login(ctx context.Context, req *LoginRequest) (string, error) {
	//获取uuid
	uuid := ctx.Value(constant.ReqUuid)
	log.Debugf(" %s| Login access from:%s,@,%s", uuid, req.UserName, req.PassWord)
	//获取用户信息
	user, err := getUserInfo(req.UserName)
	if err != nil {
		log.Errorf("Login|%v", err)
		return "", fmt.Errorf("Login|%v", err)
	}

	//用户存在，密码不正确
	if req.PassWord != user.Password {
		log.Errorf("Login|password err: req.password=%s|user.password=%s", req.PassWord, user.Password)
		return "", fmt.Errorf("password is not correct")
	}
	//生成session，并保存到缓存中
	session := utils.GenerateSession(user.Name)
	err = cache.SetSessionInfo(user, session)

	if err != nil {
		log.Errorf(" Login|Failed to SetSessionInfo, uuid=%s|user_name=%s|session=%s|err=%v", uuid, user.Name, session, err)
		return "", fmt.Errorf("login|SetSessionInfo fail:%v", err)
	}

	log.Infof("Login successfully, %s@%s with redis_session session_%s", req.UserName, req.UserName, session)
	return session, nil
}

// 退出登录
func Logout(ctx context.Context, req *LogoutRequest) error {
	uuid := ctx.Value(constant.ReqUuid)
	session := ctx.Value(constant.SessionKey).(string)
	log.Infof("%s|Logout access from,user_name=%s|session=%s", uuid, req.UserName, session)
	//要退出登录，必须要是在登录状态
	_, err := cache.GetSessionInfo(session)
	if err != nil {
		log.Errorf("%s|Failed to get with session=%s|err =%v", uuid, session, err)
		return fmt.Errorf("Logout|GetSessionInfo err:%v", err)
	}
	//删除session
	err = cache.DelSessionInfo(session)
	if err != nil {
		log.Errorf("%s|Failed to delSessionInfo :%s", uuid, session)
		return fmt.Errorf("del session err:%v", err)
	}
	log.Infof("%s|Success to delSessionInfo :%s", uuid, session)
	return nil
}

// 获取用户的信息
func GetUserInfo(ctx context.Context, req *GetUserInfoRequest) (*GetUserInfoResponse, error) {
	//从上下文中获取对应的uuid，和session
	uuid := ctx.Value(constant.ReqUuid)
	session := ctx.Value(constant.SessionKey).(string)
	log.Infof("%s|GetUserInfo access from,user_name=%s|session=%s", uuid, req.UserName, session)

	if session == "" || req.UserName == "" {
		return nil, fmt.Errorf("GetUserInfo|request params invalid")
	}
	//根据seesion获取用户信息
	user, err := cache.GetSessionInfo(session)
	if err != nil {
		log.Errorf("%s|Failed to get with session=%s|err =%v", uuid, session, err)
		return nil, fmt.Errorf("getUserInfo|GetSessionInfo err:%v", err)
	}
	//
	if user.Name != req.UserName {
		log.Errorf("%s|session info not match with username=%s", uuid, req.UserName)
	}
	log.Infof("%s|Succ to GetUserInfo|user_name=%s|session=%s", uuid, req.UserName, session)
	return &GetUserInfoResponse{
		UserName: user.Name,
		Age:      user.Age,
		Gender:   user.Gender,
		PassWord: user.Password,
		NickName: user.Nickname,
	}, nil
}

// 更新昵称
func UpdateUserNickName(ctx context.Context, req *UpdateNickNameRequest) error {
	uuid := ctx.Value(constant.ReqUuid)
	session := ctx.Value(constant.SessionKey).(string)
	log.Infof("%s|UpdateUserNickName access from,user_name=%s|session=%s", uuid, req.UserName, session)
	log.Infof("UpdateUserNickName|req==%v", req)

	if session == "" || req.UserName == "" {
		return fmt.Errorf("UpdateUserNickName|request params invalid")
	}

	user, err := cache.GetSessionInfo(session)
	if err != nil {
		log.Errorf("%s|Failed to get with session=%s|err =%v", uuid, session, err)
		return fmt.Errorf("UpdataUserNickName|GetSession err:%v", err)
	}

	if user.Name != req.UserName {
		log.Errorf("%s|Failed to get with session=%s|err =%v", uuid, session, err)
	}
	//将要更新的昵称信息封装到结构中
	updateUser := &model.User{
		Nickname: req.NewNickName,
	}

	return updateUserInfo(updateUser, req.UserName, session)
}

// 查询用户信息，先从缓存中查询用户信息，查询不到再到数据库中查询
func getUserInfo(userName string) (*model.User, error) {
	//从缓存中查询，查到了直接返回
	user, err := cache.GetUserInfoFromCache(userName)
	if err == nil && user.Name == userName {
		log.Infof("cahce_user ======= %v", user)
		return user, nil
	}
	//查不到再到数据库中查找
	user, err = dao.GetUserByName(userName)
	if err != nil {
		return user, err
	}

	if user == nil {
		return nil, fmt.Errorf("用户尚未注册")
	}
	log.Infof("user === %+v", user)
	//从数据库中查找，说明数据库中和缓存不一致，需要更新缓存
	err = cache.SetUserCacheInfo(user)
	if err != nil {
		log.Errorf("cache userinfo failed for user:", user.Name, " with err:", err.Error())
	}
	log.Infof("getUserInfo successfully, with key userinfo_%s", user.Name)
	return user, nil
}

// 更新用户信息，直接更新数据库
func updateUserInfo(user *model.User, userName, session string) error {
	//直接更新数据库中的记录
	affectedRows := dao.UpdateUserInfo(userName, user)

	//db更新
	if affectedRows == 1 {
		//数据库中的记录更新成功
		user, err := dao.GetUserByName(userName)
		if err == nil {
			//更新缓存
			cache.UpdateCachedUserInfo(user)
			if session != "" {
				//更新缓存中的session信息
				err = cache.SetSessionInfo(user, session)
				if err != nil {
					log.Errorf("update session failed:", err.Error())
					cache.DelSessionInfo(session)
				}
			}
		} else {
			log.Error("Failed to get dbUserInfo for cache, username=%s with err:", userName, err.Error())
		}
	}
	return nil
}
