package config

import (
	"errors"

	"github.com/daheige/thinkgo/common"

	"github.com/gomodule/redigo/redis"
)

var AppEnv string
var conf *common.ConfigEngine

func InitConf(path string) {
	conf = common.NewConf()
	conf.LoadConf(path + "/app.yaml")
}

func InitRedis() {
	//初始化redis
	redisConf := &common.RedisConf{}
	conf.GetStruct("RedisCommon", redisConf)

	// log.Println(redisConf)
	redisConf.SetRedisPool("default")
}

//从连接池中获取redis client
func GetRedisObj(name string) (redis.Conn, error) {
	conn := common.GetRedisClient(name)
	if conn == nil {
		return nil, errors.New("get redis client error")
	}

	return conn, nil
}
