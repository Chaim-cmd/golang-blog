package router

import (
	"net/http"

	"github.com/Chaim-cmd/golang-blog/internal/config"
	"github.com/Chaim-cmd/golang-blog/internal/handler"
	"github.com/Chaim-cmd/golang-blog/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(cfg *config.Config, db *gorm.DB) *gin.Engine {
	gin.SetMode(cfg.Server.Mode) //debug 日志详细，release 精简，上线切release

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	svc := service.NewUserService(db)
	h := handler.NewUserHandler(svc)

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	api := r.Group("/api/v1")
	{
		api.POST("/register", h.Register)
		api.POST("/login", h.Login)
	}
	return r
}
