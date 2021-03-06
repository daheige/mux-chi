package logger

import (
	"context"
	"runtime/debug"
	"strings"

	"github.com/daheige/mux-chi/app/utils"

	"github.com/daheige/thinkgo/gutils"
	"github.com/daheige/thinkgo/logger"
)

func writeLog(ctx context.Context, levelName string, message string, options map[string]interface{}) {
	reqUri := getStringByCtx(ctx, "request_uri")
	tag := strings.Replace(reqUri, "/", "_", -1)
	tag = strings.Replace(tag, ".", "_", -1)
	tag = strings.TrimLeft(tag, "_")

	logId := getStringByCtx(ctx, "log_id")
	if logId == "" {
		logId = gutils.RndUuid()
	}

	ua := getStringByCtx(ctx, "user_agent")

	logInfo := map[string]interface{}{
		"tag":            tag,
		"request_uri":    reqUri,
		"log_id":         logId,
		"options":        options,
		"ip":             getStringByCtx(ctx, "client_ip"),
		"ua":             ua,
		"plat":           utils.GetDeviceByUa(ua), // 当前设备匹配
		"request_method": getStringByCtx(ctx, "request_method"),
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

func getStringByCtx(ctx context.Context, key string) string {
	return utils.GetStringByCtx(ctx, key)
}

//Info info log.
func Info(ctx context.Context, message string, context map[string]interface{}) {
	writeLog(ctx, "info", message, context)
}

// Debug debug log.
func Debug(ctx context.Context, message string, context map[string]interface{}) {
	writeLog(ctx, "debug", message, context)
}

// Warn warn log.
func Warn(ctx context.Context, message string, context map[string]interface{}) {
	writeLog(ctx, "warn", message, context)
}

// Error error.
func Error(ctx context.Context, message string, context map[string]interface{}) {
	writeLog(ctx, "error", message, context)
}

// Emergency 致命错误或panic捕获
func Emergency(ctx context.Context, message string, context map[string]interface{}) {
	writeLog(ctx, "emergency", message, context)
}

// Recover 异常捕获处理
func Recover() {
	if err := recover(); err != nil {
		logger.DPanic("exec panic", map[string]interface{}{
			"error":       err,
			"error_trace": string(debug.Stack()),
		})
	}
}
