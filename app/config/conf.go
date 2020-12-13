package config

// AppServerConf config.
type AppServerConf struct {
	AppEnv       string
	AppDebug     bool
	AppName      string
	HttpPort     int
	PProfPort    int
	LogDir       string
	GracefulWait int
}

// AppConf app config.
var AppConf AppServerConf
