package dao

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"user_system/internal/model"
	"user_system/utils"
)

// GetUserByname 根据姓名获取用户
func GetUserByName(name string) (*model.User, error) {
	//初始化一个User模型
	user := &model.User{}
	//按照名字查找符合的第一条记录
	err := utils.GetDB().Model(user).Where("name=?", name).First(user).Error
	if err != nil {
		//没找到记录
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			return nil, nil
		}
		log.Errorf("GetUserByName fail:%v", err)
	}
	return user, nil
}

// CreatUser 创建一个用户
func CreateUser(user *model.User) error {
	if err := utils.GetDB().Model(&user).Create(user).Error; err != nil {
		log.Errorf("CreateUser fail: %v", err)
		return fmt.Errorf("CreateUser fail: %v", err)
	}
	log.Infof("insert success")
	return nil
}

// UpdateUserInfo 更新昵称，根据Username去查，昵称已经封装好在User
func UpdateUserInfo(userName string, user *model.User) int64 {
	return utils.GetDB().Model(&model.User{}).Where("name = ?", userName).Updates(user).RowsAffected
}
