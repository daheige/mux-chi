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
	"strings"
	"syscall"
	"time"

	"github.com/daheige/mux-chi/app/config"

	"github.com/daheige/thinkgo/monitor"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/daheige/mux-chi/app/routes"

	"github.com/daheige/mux-chi/app/middleware"

	"github.com/daheige/thinkgo/logger"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/go-chi/chi"

	chiWare "github.com/go-chi/chi/middleware"
)

var configDir string

func init() {
	flag.StringVar(&configDir, "config_dir", "./", "config dir")
	flag.Parse()

	// 初始化配置文件
	config.InitConf(configDir)
	config.InitRedis()

	// 日志文件设置
	logger.SetLogDir(config.AppConf.LogDir)
	logger.SetLogFile("go-web.log")
	logger.MaxSize(200)

	// 由于app/extensions/logger基于thinkgo/logger又包装了一层，所以这里是3
	logger.InitLogger(3)

	// 注册监控指标
	prometheus.MustRegister(monitor.WebRequestTotal)
	prometheus.MustRegister(monitor.WebRequestDuration)
	prometheus.MustRegister(monitor.CpuTemp)
	prometheus.MustRegister(monitor.HdFailures)

	log.Println("pprof port: ", config.AppConf.HttpPort)
	// 性能报告监控和健康检测，性能监控的端口port+1000,只能在内网访问
	go func() {
		defer logger.Recover()

		pprofAddress := fmt.Sprintf("0.0.0.0:%d", config.AppConf.PProfPort)
		log.Println("server pprof run on: ", pprofAddress)

		httpMux := http.NewServeMux() // 创建一个http ServeMux实例
		httpMux.HandleFunc("/debug/pprof/", pprof.Index)
		httpMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		httpMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		httpMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		httpMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
		httpMux.HandleFunc("/check", HealthCheckHandler)

		// metrics监控
		httpMux.Handle("/metrics", promhttp.Handler())

		if err := http.ListenAndServe(pprofAddress, httpMux); err != nil {
			log.Println(err)
		}
	}()

}

func main() {
	router := chi.NewRouter()

	// A good base middleware stack
	// router.Use(chiWare.RealIP) //获取客户端真实ip地址中间件

	// 请求中间件，记录日志和异常捕获处理
	reqWare := &middleware.RequestWare{}

	// 对请求打点监控
	router.Use(reqWare.LogAccess, reqWare.Recover)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	// request timeout 请求超时设置
	router.Use(chiWare.Timeout(10 * time.Second))

	// prometheus性能监控打点
	router.Use(monitor.MonitorHandler)

	// 加载路由
	routes.RouterHandler(router)

	// 路由找不到404处理
	router.NotFound(middleware.NotFoundHandler)

	// 路由walk,打印所有的路由信息，开发环境可以打开，生产环境可以注释掉
	if config.AppConf.AppDebug {
		walkFunc := func(method string, route string, handler http.Handler,
			middlewares ...func(http.Handler) http.Handler) error {
			route = strings.Replace(route, "/*/", "/", -1)
			fmt.Printf("router: %s %s\n", method, route)
			return nil
		}

		if err := chi.Walk(router, walkFunc); err != nil {
			fmt.Printf("Logging err: %s\n", err.Error())
		}
	}

	server := &http.Server{
		Handler: router,
		Addr:    fmt.Sprintf("0.0.0.0:%d", config.AppConf.HttpPort),
		// Good practice to set timeouts to avoid Slowloris attacks.
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 2 << 20, // header max 2MB
	}

	// 在独立携程中运行
	log.Println("server run on: ", config.AppConf.HttpPort)
	go func() {
		defer logger.Recover()

		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// mux平滑重启
	ch := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// recivie signal to exit main goroutine
	// window signal
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGHUP)

	// linux signal
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2, os.Interrupt, syscall.SIGHUP)

	// Block until we receive our signal.
	<-ch

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(config.AppConf.GracefulWait)*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// if your application should wait for other services
	// to finalize based on context cancellation.
	go server.Shutdown(ctx) // 在独立的携程中关闭服务器
	<-ctx.Done()

	log.Println("shutting down")
	os.Exit(0)
}

// 健康检测
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	w.Write([]byte(`{"alive": true}`))
}
