package routes

import (
	"net/http"

	"mux-chi/app/controller"

	"github.com/daheige/thinkgo/monitor"
	"github.com/go-chi/chi"
)

// RunRouter router list.
func RunRouter(router *chi.Mux) {
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	// 健康检测 http://localhost:1338/check
	router.Get("/check", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"code":200,"active":true"}`))
	})

	// 调用控制器上的方法
	homeCtrl := &controller.HomeController{}
	router.HandleFunc("/test", homeCtrl.Test)
	router.HandleFunc("/post", homeCtrl.Post)
	router.Post("/home/post", homeCtrl.Post) // post路由

	router.Get("/get-info", homeCtrl.Info) // 模拟参数校验

	indexCtrl := &controller.IndexController{}
	router.HandleFunc("/home", indexCtrl.Home)
	router.HandleFunc("/index/test", indexCtrl.Test)
	router.HandleFunc("/info", indexCtrl.Info)

	// 路由参数
	router.Get("/info/{userID}", indexCtrl.Info)

	// 对单个接口做性能监控打点
	router.Get("/index", monitor.MonitorHandlerFunc(indexCtrl.Home))

	// 正则路由
	// http://localhost:1338/api/user/123
	router.Get("/api/{category}/{id:[0-9]+}", indexCtrl.Category)

	// 路由前缀 /road开始，子路由设置
	router.Route("/v1", func(router chi.Router) { // chi.Router作为参数
		router.Get("/left", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("left road"))
			return
		})

		router.Post("/right", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("right road"))
		})
	})

	// 设置article_id到上下文
	// 上下文中间件处理
	// http://localhost:1338/api/articles/123fe
	router.Route("/api/articles/{articleID}", func(r chi.Router) {
		r.Use(indexCtrl.ArticleIdCtx)
		r.Get("/", indexCtrl.GetArticleId)       // GET /api/articles/123
		r.Put("/", indexCtrl.UpdateArticleId)    // PUT /api/articles/articles/123
		r.Delete("/", indexCtrl.DeleteArticleId) // DELETE /api/articles/123
	})

	// 模拟panic操作
	// http://localhost:1338/mock-panic
	router.HandleFunc("/mock-panic", indexCtrl.MockPanic)

	// 测试路由交叉的情况
	router.HandleFunc("/api/v1/hello", homeCtrl.Hello)
}
