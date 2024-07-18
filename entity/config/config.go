package config

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var once sync.Once

type Config struct {
	Server struct {
		// eg :8080
		Listen string
	}
	Redis struct {
		// 支持单机和集群版
		Addrs    []string
		Username string
		Password string `json:"-"`
		Db       int
		// 连接超时单位秒
		DialTimeout int
		// 读超时（单位秒）
		ReadTimeout int
		// 写超时（单位秒）
		WriteTimeout int
		// 最大连接数
		PoolSize int
		// 从池子中获取连接超时时间（单位秒）
		PoolTimeout int
		// 最小空闲连接数
		MinIdleConns int
		// 空闲回收时间（单位秒）
		ConnMaxIdleTime int
		// 连接最大存活时间 （单位秒）
		ConnMaxLifetime int
	}
	Log struct {
		Filename string
		// M
		MaxSize int
		// Day
		MaxAge          int
		Level           string
		ReportCaller    bool
		OutputToConsole bool
	}
}

const (
	defaultConfigPath = "./conf"
	configName        = "config"
	configType        = "yaml"
)

var (
	DebugEnable      bool
	ServerConfigPath = defaultConfigPath
	Global           = new(Config)
	RedisClient      redis.UniversalClient
)

// Init 初始化业务全局配置
// config file dir default is ./conf and can be set by flag -configPath
func Init() {
	once.Do(func() {
		// 执行文件 -h 可以查看说明
		if ServerConfigPath == defaultConfigPath {
			flag.StringVar(&ServerConfigPath, "configPath", defaultConfigPath, "server config path")
			if !flag.Parsed() {
				flag.Parse()
			}
		}
		fullConfigPath, err := filepath.Abs(ServerConfigPath)
		if err != nil {
			log.Fatalf("find configPath abs errs: %s", err)
		}
		log.Infof("read config file from %s", fullConfigPath)
		viper := viper.New()
		viper.AddConfigPath(fullConfigPath)
		viper.SetConfigName(configName)
		viper.SetConfigType(configType)
		// 环境变量不能有点，viper 对大小写的都能识别
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.AutomaticEnv()
		if err = viper.ReadInConfig(); err != nil {
			log.Fatalf("read config file errors:%s", err)
		}
		if err = viper.Unmarshal(Global); err != nil {
			log.Fatalf("unmarshal config file errors:%s", err)
		}
		if marshal, err := json.Marshal(Global); err != nil {
			log.Fatalf("json unmarshal errors:%s", err)
		} else {
			log.Infof("config is:%s\n", marshal)
		}
		initLog()
		initRedis()
	})
}

// 初始化log https://github.com/sirupsen/logrus
func initLog() {
	logConf := Global.Log
	if logConf.Level == "debug" {
		DebugEnable = true
	}
	level, err := log.ParseLevel(logConf.Level)
	if err != nil {
		panic(err)
	}
	log.SetLevel(level)

	logger := &lumberjack.Logger{
		Filename:   logConf.Filename,
		MaxSize:    logConf.MaxSize,
		MaxAge:     logConf.MaxAge,
		MaxBackups: 0,
		LocalTime:  true,
	}
	log.SetReportCaller(logConf.ReportCaller)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true, TimestampFormat: time.DateTime + ".000"})
	if !logConf.OutputToConsole {
		log.SetOutput(logger)
	}
}

func initRedis() {
	redisConf := Global.Redis
	RedisClient = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:           redisConf.Addrs,
		Password:        redisConf.Password, // no password set
		DB:              redisConf.Db,       // use default DB
		DialTimeout:     time.Duration(redisConf.DialTimeout) * time.Second,
		ReadTimeout:     time.Duration(redisConf.ReadTimeout) * time.Second,
		WriteTimeout:    time.Duration(redisConf.WriteTimeout) * time.Second,
		PoolSize:        redisConf.PoolSize,
		PoolTimeout:     time.Duration(redisConf.PoolTimeout) * time.Second,
		MinIdleConns:    redisConf.MinIdleConns,
		ConnMaxIdleTime: time.Duration(redisConf.ConnMaxIdleTime) * time.Second,
		ConnMaxLifetime: time.Duration(redisConf.ConnMaxLifetime) * time.Second,
	})
	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("connect redis err %s", err)
	}
}

func Shutdown() {
	if RedisClient != nil {
		if err := RedisClient.Close(); err != nil {
			log.Errorf("close redis error %s", err)
			return
		}
	}
}
