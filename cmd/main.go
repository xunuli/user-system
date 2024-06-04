package main

import (
	"user_system/config"
	"user_system/internal/router"
)

// 初始化配置
func init() {
	config.InitConfig()
}

// 启动服务
func main() {
	router.InitRouterAndServe()
}
