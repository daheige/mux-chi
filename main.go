package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mux-chi/app/routes"

	"mux-chi/app/middleware"

	"github.com/daheige/thinkgo/common"
	"github.com/go-chi/chi"
	chiWare "github.com/go-chi/chi/middleware"
)

var port int
var log_dir string
var config_dir string
var wait time.Duration //平滑重启的等待时间1s or 1m

func init() {
	flag.IntVar(&port, "port", 1338, "app listen port")
	flag.StringVar(&log_dir, "log_dir", "./logs", "log dir")
	flag.StringVar(&config_dir, "config_dir", "./", "config dir")
	flag.DurationVar(&wait, "graceful-timeout", 3*time.Second, "the server gracefully reload. eg: 15s or 1m")
	flag.Parse()

	//日志文件设置
	common.SetLogDir(log_dir)
	common.SetLogTime(false) //自定义日志格式

	//性能报告监控和健康检测，性能监控的端口port+1000,只能在内网访问
	go func() {
		defer common.CheckPanic()

		pprof_port := port + 1000
		log.Println("server pprof run on: ", pprof_port)

		httpMux := http.NewServeMux() //创建一个http ServeMux实例
		httpMux.HandleFunc("/debug/pprof/", pprof.Index)
		httpMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		httpMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		httpMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		httpMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
		httpMux.HandleFunc("/check", HealthCheckHandler)
		if err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", pprof_port), httpMux); err != nil {
			log.Println(err)
		}
	}()

}

func main() {
	router := chi.NewRouter()

	// A good base middleware stack
	//router.Use(chiWare.RequestID)
	//router.Use(chiWare.RealIP)
	router.Use(chiWare.Logger)
	//router.Use(chiWare.Recoverer)

	//请求中间件，记录日志和异常捕获处理
	reqWare := &middleware.RequestWare{}
	router.Use(reqWare.LogAccess, reqWare.Recover)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	// request timeout
	router.Use(chiWare.Timeout(60 * time.Second))

	//加载路由
	routes.RouterHandler(router)

	//路由找不到404处理
	router.NotFound(middleware.NotFoundHandler)

	//路由walk,打印所有的路由信息，开发环境可以打开，生产环境可以注释掉
	//walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
	//	route = strings.Replace(route, "/*/", "/", -1)
	//	fmt.Printf("%s %s\n", method, route)
	//	return nil
	//}
	//
	//if err := chi.Walk(router, walkFunc); err != nil {
	//	fmt.Printf("Logging err: %s\n", err.Error())
	//}

	server := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("0.0.0.0:%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	//在独立携程中运行
	log.Println("server run on: ", port)
	go func() {
		defer common.CheckPanic()

		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	//mux平滑重启
	ch := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// recivie signal to exit main goroutine
	//window signal
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2, os.Interrupt, syscall.SIGHUP)

	// Block until we receive our signal.
	<-ch

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// if your application should wait for other services
	// to finalize based on context cancellation.
	go server.Shutdown(ctx) //在独立的携程中关闭服务器
	<-ctx.Done()

	log.Println("shutting down")
	os.Exit(0)
}

//健康检测
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	w.Write([]byte(`{"alive": true}`))
}
