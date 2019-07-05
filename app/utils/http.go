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
	HTTP_SUCCESS_CODE = 200 //http ok 200
	HTTP_ERROR_CODE   = 500 //http error 500
	API_SUCCESS_CODE  = 0   //api success
)

//map短类型声明
type H map[string]interface{}

//空数组[]兼容其他语言php,js,python等
type EmptyArray []struct{}

//直接返回json data数据
func Json(w http.ResponseWriter, data interface{}) {
	writeJson(w, HTTP_SUCCESS_CODE, data)
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

//请求成功返回结果
//data,message
func ApiSuccess(w http.ResponseWriter, message string, data interface{}) {
	if message == "" {
		message = "ok"
	}

	writeJson(w, HTTP_SUCCESS_CODE, H{
		"code":     API_SUCCESS_CODE,
		"message":  message,
		"data":     data,
		"req_time": time.Now().Unix(),
	})
}

//错误处理code,message
func ApiError(w http.ResponseWriter, code int, message string) {
	writeJson(w, HTTP_SUCCESS_CODE, H{
		"code":     code,
		"message":  message,
		"req_time": time.Now().Unix(),
	})
}

//指定http code,message返回
func HttpCode(w http.ResponseWriter, httpCode int, message string) {
	if httpCode <= 0 {
		httpCode = HTTP_ERROR_CODE
	}

	writeJson(w, httpCode, H{
		"message": message,
	})
}
