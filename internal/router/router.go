package router

import (
	"net/http"

	"github.com/Chaim-cmd/golang-blog/internal/config"
	"github.com/gin-gonic/gin"
)

func NewRouter(cfg *config.Config) *gin.Engine {
	gin.SetMode(cfg.Server.Mode) //debug 日志详细，release 精简，上线切release

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	return r
}
