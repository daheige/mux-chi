package controller

import (
	"io/ioutil"
	"log"
	"mux-chi/app/utils"
	"net/http"
)

type HomeController struct {
	BaseController
}

func (this *HomeController) Test(w http.ResponseWriter, r *http.Request) {
	log.Println(r.FormValue("id"))
	log.Println(r.Form) //所有的form get数据 //www.hgmux.com/home/test?id=1&name=daheige
	//map[name:[daheige] id:[1]] 类型是 map[string][]string
	w.Write([]byte("ok"))
}

func (this *HomeController) Post(w http.ResponseWriter, r *http.Request) {
	log.Println(r.PostFormValue("name")) //会自动调用r.ParseForm()解析header,body

	log.Println(r.PostForm) //所有的form post数据
	log.Println(r.Body)
	body, err := ioutil.ReadAll(r.Body) //读取body内容 {"ids":{"a":1}} <nil>
	log.Println(string(body), err)
	w.Write([]byte("ok"))
}

type InfoReq struct {
	Uid   int    `json:"uid" validate:"required,min=1"`
	Limit int    `json:"limit" validate:"required,min=1,max=20"`
	Name  string `json:"name" validate:"omitempty,max=10"`
}

// Info 测试参数校验
//http://localhost:1338/get-info?limit=12&uid=1&name=abcdeabcde
func (this *HomeController) Info(w http.ResponseWriter, r *http.Request) {
	//接收参数
	req := &InfoReq{
		Uid:   this.GetInt(r.FormValue("uid")),
		Limit: this.GetInt(r.FormValue("limit")),
		Name:  r.FormValue("name"),
	}

	//校验参数
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

func (this *HomeController) Hello(w http.ResponseWriter, r *http.Request) {
	utils.ApiSuccess(w, "success", utils.H{
		"id":   1,
		"name": "mux-chi",
	})
}
