/**
http返回json格式和http code返回处理
*/
package utils

import (
	"encoding/json"
	"net/http"
	"time"
)

const (
	HttpSuccessCode = 200 // http ok 200
	HttpErrorCode   = 500 // http error 500
	ApiSuccessCode  = 0   // api success
)

// H map短类型声明
type H map[string]interface{}

// EmptyArray 空数组[]兼容其他语言php,js,python等
type EmptyArray []struct{}

// Json 直接返回json data数据
func Json(w http.ResponseWriter, data interface{}) {
	writeJson(w, HttpSuccessCode, data)
}

func writeJson(w http.ResponseWriter, httpCode int, data interface{}) {
	json_data, err := json.Marshal(data)
	if err != nil {
		json_data = []byte(`{"code":500,"message":"server error"}`)
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(httpCode)
	w.Write(json_data)
}

// ApiSuccess 请求成功返回结果
// data,message
func ApiSuccess(w http.ResponseWriter, message string, data interface{}) {
	if message == "" {
		message = "ok"
	}

	writeJson(w, HttpSuccessCode, H{
		"code":     ApiSuccessCode,
		"message":  message,
		"data":     data,
		"req_time": time.Now().Unix(),
	})
}

// ApiError 错误处理code,message
func ApiError(w http.ResponseWriter, code int, message string) {
	writeJson(w, HttpSuccessCode, H{
		"code":     code,
		"message":  message,
		"req_time": time.Now().Unix(),
	})
}

// HttpCode 指定http code,message返回
func HttpCode(w http.ResponseWriter, httpCode int, message string) {
	if httpCode <= 0 {
		httpCode = HttpErrorCode
	}

	writeJson(w, httpCode, H{
		"message": message,
	})
}
