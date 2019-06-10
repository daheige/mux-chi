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
	HTTP_SUCCESS_CODE = 200
	HTTP_ERROR_CODE   = 500
)

type H map[string]interface{}

//api 响应的结果
func Json(w http.ResponseWriter, code int, message string, data interface{}) {
	json_data, _ := json.Marshal(H{
		"code":     code,
		"message":  message,
		"data":     data,
		"req_time": time.Now().Unix(),
	})

	w.WriteHeader(HTTP_SUCCESS_CODE)
	w.Write(json_data)
}

//请求成功返回结果
//data,message
func ApiSuccess(w http.ResponseWriter, message string, data interface{}) {
	if message == "" {
		message = "ok"
	}

	Json(w, HTTP_SUCCESS_CODE, message, data)
}

//错误处理code,message
func ApiError(w http.ResponseWriter, message string) {
	Json(w, HTTP_ERROR_CODE, message, nil)
}

//指定http code,message返回
func HttpCode(w http.ResponseWriter, code int, message string) {
	if code <= 0 {
		code = HTTP_ERROR_CODE
	}

	json_data, _ := json.Marshal(map[string]interface{}{
		"code":    code,
		"message": message,
	})

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(code)
	w.Write(json_data)
}
