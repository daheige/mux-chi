package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/daheige/mux-chi/app/extensions/logger"

	"github.com/daheige/mux-chi/app/utils"

	"github.com/daheige/thinkgo/grecover"
	"github.com/daheige/thinkgo/gutils"
)

// RequestWare request middleware.
type RequestWare struct{}

// LogAccess access log.
func (reqWare *RequestWare) LogAccess(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()

		// 获取请求id
		requestId := r.Header.Get("X-Request-Id")
		// log.Println("request before")
		logId := gutils.RndUuidMd5() // 日志id

		if requestId == "" {
			requestId = logId

			// 如果采用了nginx x-request-id功能，可以注释下面一行
			w.Header().Set("x-request-id", requestId)
		}

		// log.Println("log_id: ", logId)
		// 将requestId 写入当前上下文中
		r = utils.ContextSet(r, "log_id", logId)
		r = utils.ContextSet(r, "request_id", requestId)
		// log.Println(utils.ContextGet(r, "log_id"))

		ctx := r.Context()
		// log.Println("request uri: ", r.RequestURI)
		logger.Info(ctx, "exec begin", map[string]interface{}{
			"App": "hg-mux",
		})

		h.ServeHTTP(w, r)

		// 请求结束后，记录日志
		// log.Println("request after")
		logger.Info(ctx, "exec end", map[string]interface{}{
			"App":       "hg-mux",
			"exec_time": time.Now().Sub(t).Seconds(),
		})
	})
}

// Recover当请求发生了异常或致命错误，需要捕捉r,w执行上下文的错误
// 该Recover设计灵感来源于golang gin框架的WriteHeaderNow()设计
func (reqWare *RequestWare) Recover(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				bytes := grecover.CatchStack()

				ctx := r.Context()
				if ctx == nil {
					ctx = context.Background()
				}

				logger.Error(ctx, "exec recover error", map[string]interface{}{
					"trace_error": string(bytes),
				})

				// 当http请求发生了recover或异常就直接终止
				// w.Header().Set("X-Content-Type-Options", "nosniff")

				utils.HttpCode(w, http.StatusInternalServerError, "server error!")
				return
			}
		}()

		h.ServeHTTP(w, r)
	})
}

// NotFoundHandler 404处理函数
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	utils.HttpCode(w, http.StatusNotFound, "404 - page not found")
}
