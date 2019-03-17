package controller

import (
	"io/ioutil"
	"log"
	"net/http"
)

type HomeController struct{}

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
