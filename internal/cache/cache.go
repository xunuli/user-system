package cache

import (
	"context"
	"encoding/json"
	"time"
	"user_system/config"
	"user_system/internal/model"
	"user_system/pkg/constant"
	"user_system/utils"
)

// 从缓存中获取信息
func GetUserInfoFromCache(username string) (*model.User, error) {
	//用户信息的key为前缀+用户名
	redisKey := constant.UserInfoPrefix + username
	//根据key获取对应的value，用户信息
	val, err := utils.GetRediscli().Get(context.Background(), redisKey).Result()
	if err != nil {
		return nil, err
	}
	user := &model.User{}
	//将json(字符串)格式的用户信息转化为结构体形式的信息
	err = json.Unmarshal([]byte(val), user)
	return user, err
}

// 设置用户缓存信息
func SetUserCacheInfo(user *model.User) error {
	//更具用户的名称得到对应的Key
	redisKey := constant.UserInfoPrefix + user.Name
	//将struct类型的用户信息序列化为json形式的信息
	val, err := json.Marshal(user)
	if err != nil {
		return err
	}
	//设置对应的过期时间
	expired := time.Second * time.Duration(config.GetGlobalConf().Cache.UserExpired)
	//设置对应的key，value和ex
	_, err = utils.GetRediscli().Set(context.Background(), redisKey, val, expired*time.Second).Result()
	return err
}

// 更新用户缓存信息
func UpdateCachedUserInfo(user *model.User) error {
	//直接设置用户信息（直接更新）
	err := SetUserCacheInfo(user)
	if err != nil {
		//更新失败，则尝试删除缓存
		redisKey := constant.SessionKeyPrefix + user.Name
		utils.GetRediscli().Del(context.Background(), redisKey).Result()
	}
	return err
}

// 获取缓存的Session信息
func GetSessionInfo(session string) (*model.User, error) {
	//生成对应的Key
	redisKey := constant.UserInfoPrefix + session
	//依据key获取对应的session信息
	val, err := utils.GetRediscli().Get(context.Background(), redisKey).Result()
	if err != nil {
		return nil, err
	}
	user := &model.User{}
	//反序列化为结构体
	err = json.Unmarshal([]byte(val), &user)
	return user, err
}

// 设置session缓存信息
func SetSessionInfo(user *model.User, session string) error {
	//生成对应的Key
	redisKey := constant.SessionKeyPrefix + session
	//将用户信息序列化为对应的value
	val, err := json.Marshal(&user)
	if err != nil {
		return err
	}
	expired := time.Second * time.Duration(config.GetGlobalConf().Cache.SessionExpired)
	_, err = utils.GetRediscli().Set(context.Background(), redisKey, val, expired*time.Second).Result()
	return err
}

// 删除对应的session信息
func DelSessionInfo(session string) error {
	redisKey := constant.UserInfoPrefix + session
	_, err := utils.GetRediscli().Del(context.Background(), redisKey).Result()
	return err
}
