package logger

import (
	"mux-chi/app/utils"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/daheige/thinkgo/logger"
)

func writeLog(req *http.Request, levelName string, message string, opts map[string]interface{}) {
	tag := strings.Replace(req.RequestURI, "/", "_", -1)
	ua := req.Header.Get("User-Agent")

	//日志信息
	logId := req.Context().Value("log_id")

	//函数调用信息
	_, file, line, _ := runtime.Caller(2)

	logInfo := map[string]interface{}{
		"tag":         strings.TrimLeft(tag, "_"),
		"request_uri": req.RequestURI,
		"log_id":      logId,
		"host":        req.Host, //host
		"trace_line":  line,
		"trace_file":  filepath.Base(file),
		"ua":          ua,
		"client_ip":   req.RemoteAddr,          //客户端真实ip地址
		"plat":        utils.GetDeviceByUa(ua), //当前设备匹配
		"method":      req.Method,
	}

	if len(opts) > 0 {
		logInfo["context"] = opts
	}

	switch levelName {
	case "info":
		logger.Info(message, logInfo)
	case "debug":
		logger.Debug(message, logInfo)
	case "warn":
		logger.Warn(message, logInfo)
	case "error":
		logger.Error(message, logInfo)
	case "emergency":
		logger.DPanic(message, logInfo)
	default:
		logger.Info(message, logInfo)
	}
}

func Info(req *http.Request, message string, options map[string]interface{}) {
	writeLog(req, "info", message, options)
}

func Debug(req *http.Request, message string, options map[string]interface{}) {
	writeLog(req, "debug", message, options)
}

func Warn(req *http.Request, message string, options map[string]interface{}) {
	writeLog(req, "warn", message, options)
}

func Error(req *http.Request, message string, options map[string]interface{}) {
	writeLog(req, "error", message, options)
}

func Emergency(req *http.Request, message string, options map[string]interface{}) {
	writeLog(req, "emergency", message, options)
}

