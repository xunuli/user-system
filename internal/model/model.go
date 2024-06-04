package model

import "time"

// CreateModel 内嵌model
type CreateModel struct {
	Creator    string    `gorm:"type:varchar(100);not null;default ''"`
	CreateTime time.Time `gorm:"autoCreateTime"` //在创建时自动生成时间
}

// ModifyModel 内嵌model
type ModifyModel struct {
	Modifier   string    `gorm:"type:varchar(100); not null; default ''"`
	ModifyTime time.Time `gorm:"autoUpdateTime"` //在更新记录时自动生成时间
}

// 数据库表
type User struct {
	ID       int    `gorm:"cloumn:id"`
	Name     string `gorm:"cloumn:name"`     //姓名
	Age      int    `gorm:"cloumn:age"`      //年龄
	Gender   string `gorm:"cloumn:gender"`   //性别
	Password string `gorm:"cloumn:password"` //密码
	Nickname string `gorm:"cloumn:nickname"` //昵称
	CreateModel
	ModifyModel
}

// TableName 表名
func (t *User) TableName() string {
	return "t_user"
}
