package routes

import (
	"mux-chi/app/controller"
	"net/http"

	"github.com/go-chi/chi"
)

func RouterHandler(router *chi.Mux) {
	//测试get/post接收数据
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	//调用控制器上的方法
	homeCtrl := &controller.HomeController{}
	router.HandleFunc("/test", homeCtrl.Test)
	router.HandleFunc("/post", homeCtrl.Post)
	router.Post("/home/post",homeCtrl.Post) //post路由

	indexCtrl := &controller.IndexController{}
	router.HandleFunc("/home", indexCtrl.Home)
	router.HandleFunc("/index", indexCtrl.Home)
	router.HandleFunc("/index/test", indexCtrl.Test)
	router.HandleFunc("/info", indexCtrl.Info)

	//路由参数
	router.Get("/info/{userID}",indexCtrl.Info)

	//模拟panic操作
	//http://localhost:1338/mock-panic
	router.HandleFunc("/mock-panic", indexCtrl.MockPanic)
}
