package router

import (
	user "go-app/business"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	// 註冊
	r.POST("/user", user.CreateUser)
	// 登入
	r.POST("/login", user.Login)
	// 更新會員資料√
	r.PUT("/user/:id", user.UpdateUser)
	// 刪除會員
	r.DELETE("/user/:id", user.DeleteUser)
	// 取得會員資料列表
	r.GET("/users", user.GetUsers)
	// 確認會員是否在線
	r.GET("/user/:username", user.CheckUser)

	return r
}
