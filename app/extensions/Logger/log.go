package Logger

import (
	"encoding/json"
	"log"
	"mux-chi/app/utils"
	"net/http"
	"strings"
	"time"

	"github.com/daheige/thinkgo/common"
)

const logtmFmtWithMS = "2006-01-02 15:04:05.999"

//日志基本信息
type LogInfo struct {
	Tag        string      `json:"tag"` //uri路径
	Message    interface{} `json:"message"`
	RequestUri string      `json:"request_uri"` //请求的uri
	LogId      string      `json:"log_id"`      //上下文请求的日志id
	LocalTime  string      `json:"local_time"`
	Context    interface{} `json:"context"` //当前请求上下文
	Host       string      `json:"host"`
	Ua         string      `json:"ua"`     //请求的ua
	Platform   string      `json:"plat"`   //请求平台
	Method     string      `json:"method"` //请求的方法 post,get,put,delete等
}

func writeLog(req *http.Request, levelName string, message interface{}, context interface{}) {
	defer func() {
		if err := recover(); err != nil {
			//请求上下文中的log_id
			reqLogId := req.Context().Value("log_id")
			logId := ""
			if reqLogId == nil {
				logId = common.StrJoin("_", "log", common.RndUuidMd5())
			} else {
				logId = reqLogId.(string)
			}

			//记录堆栈信息
			bytes := common.CatchStack()
			ua := req.Header.Get("User-Agent")
			logInfo := &LogInfo{
				Tag:        strings.Replace(req.URL.Path, "/", "_", -1),
				Message:    "write log error",
				RequestUri: req.RequestURI,
				LogId:      logId,
				LocalTime:  time.Now().Format(logtmFmtWithMS),
				Context: map[string]interface{}{
					"trace_error": string(bytes),
				},
				Host:     req.RemoteAddr,
				Ua:       ua,
				Platform: utils.GetDeviceByUa(ua),
				Method:   req.Method,
			}

			json_data, _ := json.Marshal(logInfo)
			common.ErrorLog(string(json_data))
		}
	}()

	tag := strings.Replace(req.RequestURI, "/", "_", -1)
	ua := req.Header.Get("User-Agent")

	//日志信息
	log_id := req.Context().Value("log_id")

	// log.Println("log_id: ", log_id)

	if _, ok := context.(map[string]interface{}); !ok {
		context = nil
	}

	logInfo := &LogInfo{
		Tag:        tag,
		Message:    message,
		RequestUri: req.RequestURI,
		LogId:      log_id.(string),
		LocalTime:  time.Now().Format(logtmFmtWithMS),
		Context:    context,
		Host:       req.RemoteAddr,
		Ua:         ua,
		Platform:   utils.GetDeviceByUa(ua), //当前设备匹配
		Method:     req.Method,
	}

	json_data, err := json.Marshal(logInfo)
	if err != nil {
		log.Println(err)
		return
	}

	switch levelName {
	case "info":
		common.InfoLog(string(json_data))
	case "debug":
		common.DebugLog(string(json_data))
	case "notice":
		common.NoticeLog(string(json_data))
	case "warn":
		common.WarnLog(string(json_data))
	case "error":
		common.ErrorLog(string(json_data))
	}
}

func Info(req *http.Request, message interface{}, context map[string]interface{}) {
	writeLog(req, "info", message, context)
}

func Debug(req *http.Request, message interface{}, context map[string]interface{}) {
	writeLog(req, "debug", message, context)
}

func Notice(req *http.Request, message interface{}, context map[string]interface{}) {
	writeLog(req, "notice", message, context)
}

func Warn(req *http.Request, message interface{}, context map[string]interface{}) {
	writeLog(req, "warn", message, context)
}

func Error(req *http.Request, message interface{}, context map[string]interface{}) {
	writeLog(req, "error", message, context)
}
