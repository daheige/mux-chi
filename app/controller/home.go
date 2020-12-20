package controller

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/daheige/mux-chi/app/utils"
)

// HomeController home ctrl.
type HomeController struct {
	BaseController
}

// Test test.
func (h *HomeController) Test(w http.ResponseWriter, r *http.Request) {
	log.Println(r.FormValue("id"))
	log.Println(r.Form) // 所有的form get数据 //www.hgmux.com/home/test?id=1&name=daheige
	// map[name:[daheige] id:[1]] 类型是 map[string][]string
	w.Write([]byte("ok"))
}

// Post post.
func (h *HomeController) Post(w http.ResponseWriter, r *http.Request) {
	log.Println(r.PostFormValue("name")) // 会自动调用r.ParseForm()解析header,body

	log.Println(r.PostForm) // 所有的form post数据
	log.Println(r.Body)
	body, err := ioutil.ReadAll(r.Body) // 读取body内容 {"ids":{"a":1}} <nil>
	log.Println(string(body), err)
	w.Write([]byte("ok"))
}

// InfoReq info req.
type InfoReq struct {
	Uid   int    `json:"uid" validate:"required,min=1"`
	Limit int    `json:"limit" validate:"required,min=1,max=20"`
	Name  string `json:"name" validate:"omitempty,max=10"`
}

// Info 测试参数校验
// http://localhost:1338/get-info?limit=12&uid=1&name=abcdeabcde
func (h *HomeController) Info(w http.ResponseWriter, r *http.Request) {
	// 接收参数
	req := &InfoReq{
		Uid:   h.GetInt(r.FormValue("uid")),
		Limit: h.GetInt(r.FormValue("limit")),
		Name:  r.FormValue("name"),
	}

	// 校验参数
	err := validate.Struct(req)
	if err != nil {
		utils.ApiError(w, 1001, "param error")
		return
	}

	if req.Limit == 0 {
		utils.ApiSuccess(w, "ok", utils.EmptyArray{})
		return
	}

	utils.ApiSuccess(w, "success", utils.H{
		"uid":   req.Uid,
		"limit": req.Limit,
	})
}

// Hello hello func.
func (h *HomeController) Hello(w http.ResponseWriter, r *http.Request) {
	utils.ApiSuccess(w, "success", utils.H{
		"id":   1,
		"name": "mux-chi",
	})
}
