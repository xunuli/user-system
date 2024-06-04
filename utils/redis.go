package utils

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"sync"
	"user_system/config"
)

//连接redis

var (
	redisConn *redis.Client
	redisOnce sync.Once
)

// 初始化redis连接
func initRedis() {
	redisConfig := config.GetGlobalConf().RedisConfig //获取redis配置
	log.Infof("redisConfig======%+v", redisConfig)    //打印日志信息
	addr := fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port)
	//初始化一个客户端
	redisConn = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: redisConfig.PassWord,
		DB:       redisConfig.Rdb,
		PoolSize: redisConfig.PoolSize,
	})
	if redisConn == nil {
		panic("failed to call redis.NewClient")
	}
	//设置值，测试
	res, err := redisConn.Set(context.Background(), "abc", 100, 60).Result()
	log.Infof("res ====== %v, err ======== %v", res, err)
	pong, err := redisConn.Ping(context.Background()).Result()
	fmt.Println(pong)
	if err != nil {
		panic("failed to ping redis, err:" + err.Error())
	}
}

func GetRediscli() *redis.Client {
	redisOnce.Do(initRedis)
	return redisConn
}
