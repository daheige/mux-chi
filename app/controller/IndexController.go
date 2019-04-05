package controller

import (
	"log"
	"mux-chi/app/model"
	"net/http"

	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-xorm/xorm"
	"github.com/jinzhu/gorm"

	"mux-chi/app/config"
	"mux-chi/app/utils"

	"github.com/gomodule/redigo/redis"
)

type IndexController struct{}

func (this *IndexController) Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello hg-mux"))
}

func (this *IndexController) Test(w http.ResponseWriter, r *http.Request) {
	// log.Println("log_id: ", r.Context().Value("log_id"))
	conn, err := config.GetRedisObj("default")
	log.Println("err: ", err)
	if err != nil {
		utils.ApiError(w, "redis connection error")
		return
	}

	defer conn.Close()

	v, err := redis.String(conn.Do("get", "myname"))
	log.Println(v, err)
	utils.ApiSuccess(w, "ok: "+v, nil)
}

//from a route like /info/{userID}
func (this *IndexController) Info(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	log.Println("uid: ", userID)

	utils.ApiSuccess(w, "hello world", nil)
}

//模拟发生panic操作
func (this *IndexController) MockPanic(w http.ResponseWriter, r *http.Request) {
	panic(111)
}

func (this *IndexController) User(w http.ResponseWriter, r *http.Request) {
	dbObj := r.Context().Value("db")
	db := dbObj.(*gorm.DB)

	log.Println(dbObj)

	user := &model.User{}
	err := db.Where("name = ?", "hello").First(user).Error
	if err != nil {
		log.Println("get data error: ", err.Error())
		w.Write([]byte("get data error"))
		return
	}

	//defer db.Close() //当db在中间件退出的时候已经关闭了，也不会影响服务运行

	log.Println("user: ", user)
	log.Println(user.Id)

	w.Write([]byte("ok"))
}

func (this *IndexController) GetUser(w http.ResponseWriter, r *http.Request) {
	dbObj := r.Context().Value("db2")
	if dbObj == nil{
			w.Write([]byte("db2 not exist"))
			return
	}

	db := dbObj.(*xorm.Engine)

	//由于xorm是使用的时候，才真正建立连接，这里使用的是短连接方式，因此这里需要关闭，也可以在中间件中进行关闭
	//defer db.Close()

	log.Println("db: ", db)
	user := &model.User{}
	has, err := db.Where("name = ?", "hello").Get(user)
	if err != nil || !has {
		if err != nil {
			log.Println("get data error: ", err.Error())
		}

		w.Write([]byte("get data error"))
		return
	}

	log.Println("user: ", user)
	log.Println(user.Id)
	jsonData, _ := json.Marshal(user)
	w.Write(jsonData)
}
