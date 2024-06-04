package router

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	api "user_system/api/http/v1"
	"user_system/config"
	"user_system/pkg/constant"
)
import "github.com/gin-gonic/gin"

//设置路由

// InitRouterAndServer 路由配置、启动服务
func InitRouterAndServe() {
	//设置服务的运行模式
	setAppRunMode()
	//启动引擎
	router := gin.Default()
	//健康检查
	router.GET("ping", api.Ping)
	//用户注册
	router.POST("/user/register", api.Register)
	//用户登录
	router.POST("/user/login", api.Login)
	//用户登出
	router.POST("/user/logout", AuthMiddleWare(), api.Logout)
	//获取用户信息
	router.GET("/user/get_user_info", AuthMiddleWare(), api.GetUserInfo)
	//更新用户信息
	router.POST("/user/update_nick_name", AuthMiddleWare(), api.UpdateNickName)

	//配置静态文件服务的函数，简而言之就是当浏览器访问前面的路径时，会访问实际的后面的路径，（后面路径中是项目根目录下的文件夹）
	router.Static("/static/", "./web/static/")
	router.Static("/upload/images/", "./web/upload/images/")

	//启动服务
	port := config.GetGlobalConf().AppConfig.Port

	if err := router.Run(":" + port); err != nil {
		log.Error("start server err:" + err.Error())
	}

}

// setAppRunMode 设置运行模式
func setAppRunMode() {
	if config.GetGlobalConf().AppConfig.RunMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
}

// 中间件，查询请求中是否存在session信息，存在说明合法，继续向下执行，否则终止执行
func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		if session, err := c.Cookie(constant.SessionKey); err == nil {
			if session != "" {
				//继续向下执行
				c.Next()
				return
			}
		}
		//返回错误
		c.JSON(http.StatusUnauthorized, gin.H{"error": "err"})
		//执行完当前中间件后，出栈返回，后续中间件以及路由处理函数不执行
		c.Abort()
		return
	}
}
