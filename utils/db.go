package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"time"
	"user_system/config"
)

//连接数据库mysql

// 定义全局变量数据库
var (
	db     *gorm.DB
	dbOnce sync.Once
)

func openDB() {
	mysqlConf := config.GetGlobalConf().DbConfig //获取数据库配置
	//拼接字符串，获取dsn
	connArgs := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", mysqlConf.User,
		mysqlConf.PassWord, mysqlConf.Host, mysqlConf.Port, mysqlConf.DbName)
	log.Info("mysql addr:" + connArgs) //加载日志

	var err error
	db, err = gorm.Open(mysql.Open(connArgs), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	//_ = db.AutoMigrate(model.User{})
	//生成通用数据库接口
	sqlDB, err := db.DB()
	if err != nil {
		panic("fetch db connection err:" + err.Error())
	}
	sqlDB.SetMaxIdleConns(mysqlConf.MaxIdleConn)                                        //设置最大空闲连接
	sqlDB.SetMaxOpenConns(mysqlConf.MaxOpenConn)                                        //设置最大连接数
	sqlDB.SetConnMaxLifetime(time.Duration(mysqlConf.MaxIdleTime * int64(time.Second))) //设置空间时间
}

// 获取数据库连接，只执行一次
func GetDB() *gorm.DB {
	dbOnce.Do(openDB)
	return db
}
