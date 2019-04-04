package middleware

import (
	"github.com/daheige/thinkgo/mysql"
	"log"
	"mux-chi/app/extensions/Logger"
	"net/http"
	"time"

	"mux-chi/app/utils"

	"github.com/daheige/thinkgo/common"
)

type RequestWare struct{}

func (this *RequestWare) LogAccess(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()

		//获取请求id
		requestId := r.Header.Get("X-Request-Id")
		log.Println("request before")
		logId := common.RndUuidMd5() //日志id

		if requestId == "" {
			requestId = logId

			//如果采用了nginx x-request-id功能，可以注释下面一行
			w.Header().Set("x-request-id", requestId)
		}

		//log.Println("log_id: ", logId)
		//将requestId 写入当前上下文中
		r = utils.ContextSet(r, "log_id", logId)
		r = utils.ContextSet(r, "request_id", requestId)
		// log.Println(utils.ContextGet(r, "log_id"))

		//log.Println("request uri: ", r.RequestURI)
		Logger.Info(r, "exec begin", map[string]interface{}{
			"App": "hg-mux",
		})

		dbConf := mysql.DbConf{
			Ip:        "127.0.0.1",
			Port:      3306,
			User:      "root",
			Password:  "1234",
			Database:  "test",
			ParseTime: true,
			SqlCmd:    true,
		}

		err := dbConf.ShortConnect()
		if err != nil{
			log.Println("db connection error: ",err.Error())
		}else{
			r = utils.ContextSet(r, "db", dbConf.Db())

			//用完短连接建议关闭，因为资源不释放的话，当大量的请求过来的时候，mysql连接过多，就会报错
			//查看db连接信息：show full processlist;或show processlist;
			//就可以看到当前数据库连接信息状态
			defer dbConf.Db().Close()
		}


		h.ServeHTTP(w, r)

		//请求结束后，记录日志
		log.Println("request after")
		Logger.Info(r, "exec end", map[string]interface{}{
			"App":       "hg-mux",
			"exec_time": time.Now().Sub(t).Seconds(),
		})
	})
}

//当请求发生了异常或致命错误，需要捕捉r,w执行上下文的错误
//该Recover设计灵感来源于golang gin框架的WriteHeaderNow()设计
func (this *RequestWare) Recover(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				bytes := common.CatchStack()
				if len(bytes) > 0 {
					Logger.Error(r, "exec recover error", map[string]interface{}{
						"trace_error": string(bytes),
					})
				}

				//当http请求发生了recover或异常就直接终止
				utils.HttpCode(w, http.StatusInternalServerError, "server error!")
				return
			}
		}()

		if r.RequestURI == "/favicon.ico" {
			w.Write([]byte("ok"))
			return
		}

		h.ServeHTTP(w, r)
	})
}

//404处理函数
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 - page not found"))
}
