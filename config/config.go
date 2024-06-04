package config

import (
	rlog "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"sync"
	"time"
)

// 定义全局变量，全局配置，初始化一次
var (
	config GlobalConfig //全局业务的配置文件
	once   sync.Once    //用于初始化，只执行一次
)

// LogConf 日志配置
type LogConf struct {
	LogPattern string `yaml:"log_pattern" mapstructure:"log_pattern"` //日志输出标准，终端输出/文件输出
	LogPath    string `yaml:"log_path" mapstructure:"log_path"`       //日志路径
	SaveDays   uint   `yaml:"save_days" mapstructure:"save_days"`     //日志保存天数
	Level      string `yaml:"level" mapstructure:"level"`             //日志级别
}

// DBConf数据库配置
type DbConf struct {
	Host        string `yaml:"host" mapstructure:"host"`                   //db主机地址
	Port        string `yaml:"port" mapstructure:"port"`                   //db端口
	User        string `yaml:"user" mapstructure:"user"`                   //用户名
	PassWord    string `yaml:"password" mapstructure:"password"`           //密码
	DbName      string `yaml:"dbname" mapstructure:"dbname"`               //db名
	MaxIdleConn int    `yaml:"max_idle_conn" mapstructure:"max_idle_conn"` //最大空闲连接
	MaxOpenConn int    `yaml:"max_open_conn" mapstructure:"max_open_conn"` //最大打开的连接数
	MaxIdleTime int64  `yaml:"max_idle_time" mapstructure:"max_idle_time"` //连接最大空闲时间
}

// RedisConf配置
type RedisConf struct {
	Host     string `yaml:"host" mapstructure:"host"`         //redis主机地址
	Port     string `yaml:"port" mapstructure:"port"`         //redis端口
	Rdb      int    `yaml:"rdb" mapstructure:"rdb"`           //redis db名
	PassWord string `yaml:"password" mapstructure:"password"` //密码
	PoolSize int    `yaml:"poolsize" mapstructure:"poolsize"`
}

// cache配置
type Cache struct {
	SessionExpired int `yaml:"session_expired" mapstructure:"session_expired"` //会话过期时间
	UserExpired    int `yaml:"user_expired" mapstructure:"user_expired"`       //用户过期时间
}

// AppConf 服务配置
type AppConf struct {
	AppName string `yaml:"app_name" mapstructure:"app_name"` //业务名
	Version string `yaml:"version" mapstructure:"version"`   //版本号
	Port    string `yaml:"port" mapstructure:"port"`         //业务端口
	RunMode string `yaml:"run_mode" mapstructure:"run_mode"` //运行模式
}

// GlobalConfig 业务配置结构体
type GlobalConfig struct {
	LogConfig   LogConf   `yaml:"log" mapstructure:"log"`                   //日志配置
	DbConfig    DbConf    `yaml:"db" mapstructure:"db"`                     //db配置
	RedisConfig RedisConf `yaml:"redis" mapstructure:"redis"`               //Redis配置
	Cache       Cache     `yaml:"cache" mapstructure:"cache"`               //缓存
	AppConfig   AppConf   `yaml:"app" mapstructure:"app"`                   //服务配置
	CoresOrigin []string  `yaml:"cores_origin" mapstructure:"cores_origin"` //跨域源列表
}

// GetGlobalConf获取全局配置文件，并初始化服务
func GetGlobalConf() *GlobalConfig {
	once.Do(readConf)
	return &config
}

func readConf() {
	//viper是一个go语言的配置管理工具
	viper.SetConfigName("app")     //为配置文件设置名字
	viper.SetConfigType("yml")     //配置文件的格式
	viper.AddConfigPath("./conf")  //配置文件的路径
	viper.AddConfigPath(".")       //配置文件的路径
	viper.AddConfigPath("../conf") //配置文件的路径
	err := viper.ReadInConfig()    //发现并加载配置文件
	if err != nil {
		panic("read config file err:" + err.Error())
	}
	err = viper.Unmarshal(&config) //将配置文件反序列到config结构体
	if err != nil {
		panic("config file unmarshal err:" + err.Error())
	}
	//logrus日志管理工具
	log.Infof("config === %+v", config)
}

// 初始化日志
func InitConfig() {
	globalConf := GetGlobalConf()
	// 设置日志级别
	level, err := log.ParseLevel(globalConf.LogConfig.Level)
	if err != nil {
		panic("log level parse err:" + err.Error())
	}
	//设置日志格式为json格式
	log.SetFormatter(&logFormatter{
		log.TextFormatter{
			DisableColors:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		},
	})
	log.SetReportCaller(true) // 打印文件位置，行号
	log.SetLevel(level)
	switch globalConf.LogConfig.LogPattern {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "stderr":
		log.SetOutput(os.Stderr)
	case "file":
		logger, err := rlog.New(
			globalConf.LogConfig.LogPath+".%Y%m%d",
			rlog.WithLinkName(globalConf.LogConfig.LogPath),
			rlog.WithRotationCount(globalConf.LogConfig.SaveDays),
			//rlog.WithMaxAge(time.Minute*3),
			rlog.WithRotationTime(time.Hour*24),
		)
		if err != nil {
			panic("log conf err: " + err.Error())
		}
		log.SetOutput(logger)
	default:
		panic("log conf err, check log_pattern in logsvr.yaml")
	}

}
