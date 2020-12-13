package controller

import (
	"context"
	"log"
	"net/http"

	"mux-chi/app/extensions/logger"

	"mux-chi/app/config"
	"mux-chi/app/utils"

	"github.com/go-chi/chi"

	"github.com/gomodule/redigo/redis"
)

// IndexController index ctrl.
type IndexController struct{}

// Home home.
func (ctrl *IndexController) Home(w http.ResponseWriter, r *http.Request) {
	logger.Info(r.Context(), "fefefe", nil)

	w.Write([]byte("hello hg-mux"))
}

// Test test.
func (ctrl *IndexController) Test(w http.ResponseWriter, r *http.Request) {
	// log.Println("log_id: ", r.Context().Value("log_id"))
	conn, err := config.GetRedisObj("default")
	log.Println("err: ", err)
	if err != nil {
		utils.ApiError(w, 1001, "redis connection error")
		return
	}

	defer conn.Close()

	v, err := redis.String(conn.Do("get", "myname"))
	log.Println(v, err)
	utils.ApiSuccess(w, "ok: "+v, nil)
}

// Info from a route like /info/{userID}
func (ctrl *IndexController) Info(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	log.Println("uid: ", userID)

	utils.ApiSuccess(w, "hello world", nil)
}

// MockPanic 模拟发生panic操作
func (ctrl *IndexController) MockPanic(w http.ResponseWriter, r *http.Request) {
	panic(111)
}

// Category category
func (ctrl *IndexController) Category(w http.ResponseWriter, r *http.Request) {
	category := chi.URLParam(r, "category")
	id := chi.URLParam(r, "id")

	utils.Json(w, utils.H{
		"code":    200,
		"message": "ok",
		"data": map[string]interface{}{
			"category": category,
			"id":       id,
		},
	})
}

// ArticleIdCtx 设置http context value 中间件
func (ctrl *IndexController) ArticleIdCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		articleID := chi.URLParam(r, "articleID")

		// 设置上下文article_id
		ctx := context.WithValue(r.Context(), "article_id", articleID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetArticleId get art id.
func (ctrl *IndexController) GetArticleId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	articleId, ok := ctx.Value("article_id").(string)
	if !ok {
		utils.ApiError(w, 422, http.StatusText(422))
		return
	}

	utils.ApiSuccess(w, "ok", utils.H{
		"articleId": articleId,
	})
}

// UpdateArticleId update art id.
func (ctrl *IndexController) UpdateArticleId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	articleId, ok := ctx.Value("article_id").(string)
	if !ok {
		utils.ApiError(w, 422, http.StatusText(422))
		return
	}

	articleId = "456"
	utils.ApiSuccess(w, "ok", utils.H{
		"articleId": articleId,
	})
}

// DeleteArticleId delete.
func (ctrl *IndexController) DeleteArticleId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	articleId, ok := ctx.Value("article_id").(string)
	if !ok {
		utils.ApiError(w, 422, http.StatusText(422))
		return
	}

	log.Println("article_id: ", articleId)

	articleId = ""
	utils.ApiSuccess(w, "delete success", nil)
}
