package config

import (
	"errors"
	"log"
	"path/filepath"
	"time"

	"github.com/daheige/thinkgo/gredigo"
	"github.com/daheige/thinkgo/yamlconf"

	"github.com/gomodule/redigo/redis"
)

var (
	conf *yamlconf.ConfigEngine
)

// InitConf init config.
func InitConf(path string) {
	dir, err := filepath.Abs(path)
	if err != nil {
		log.Fatalln("config_dir path error: ", err)
	}

	conf = yamlconf.NewConf()
	err = conf.LoadConf(filepath.Join(dir, "app.yaml"))
	if err != nil {
		log.Fatalln("load config error: ", err)
	}

	InitAppConfig()
}

// InitAppConfig init app config.
func InitAppConfig() {
	conf.GetStruct("AppServer", &AppConf)
	if AppConf.AppEnv == "" {
		AppConf.AppEnv = "production"
	}

	if AppConf.HttpPort == 0 {
		AppConf.HttpPort = 1338
	}

	if AppConf.PProfPort == 0 {
		AppConf.PProfPort = AppConf.HttpPort + 1000
	}

	if AppConf.GracefulWait == 0 {
		AppConf.GracefulWait = 5 * time.Second
	} else {
		AppConf.GracefulWait *= time.Second
	}

	log.Println("app: ", AppConf)
}

// InitRedis初始化redis
func InitRedis() {
	redisConf := &gredigo.RedisConf{}
	conf.GetStruct("RedisCommon", redisConf)

	// log.Println(redisConf)
	redisConf.SetRedisPool("default")
}

// 从连接池中获取redis client
func GetRedisObj(name string) (redis.Conn, error) {
	conn := gredigo.GetRedisClient(name)
	if conn == nil {
		return nil, errors.New("get redis client error")
	}

	if conn.Err() != nil {
		return nil, conn.Err()
	}

	return conn, nil
}
