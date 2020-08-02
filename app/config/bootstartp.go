package config

import (
	"errors"

	"github.com/daheige/thinkgo/gredigo"
	"github.com/daheige/thinkgo/yamlconf"

	"github.com/gomodule/redigo/redis"
)

var AppEnv string
var conf *yamlconf.ConfigEngine

func InitConf(path string) {
	conf = yamlconf.NewConf()
	conf.LoadConf(path + "/app.yaml")
}

func InitRedis() {
	//初始化redis
	redisConf := &gredigo.RedisConf{}
	conf.GetStruct("RedisCommon", redisConf)

	// log.Println(redisConf)
	redisConf.SetRedisPool("default")
}

//从连接池中获取redis client
func GetRedisObj(name string) (redis.Conn, error) {
	conn := gredigo.GetRedisClient(name)
	if conn == nil {
		return nil, errors.New("get redis client error")
	}

	return conn, nil
}
