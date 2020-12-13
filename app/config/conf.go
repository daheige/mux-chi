package config

import (
	"time"
)

// AppServerConf config.
type AppServerConf struct {
	AppEnv       string
	AppDebug     bool
	AppName      string
	HttpPort     int
	PProfPort    int
	LogDir       string
	GracefulWait time.Duration
}

// AppConf app config.
var AppConf AppServerConf
