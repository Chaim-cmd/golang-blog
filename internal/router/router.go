package router

import (
	"net/http"

	"github.com/Chaim-cmd/golang-blog/internal/config"
	"github.com/Chaim-cmd/golang-blog/internal/handler"
	"github.com/Chaim-cmd/golang-blog/internal/logger"
	"github.com/Chaim-cmd/golang-blog/internal/middleware"
	"github.com/Chaim-cmd/golang-blog/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(cfg *config.Config, db *gorm.DB, noticeSvc *service.NoticeService) *gin.Engine {
	gin.SetMode(cfg.Server.Mode) //debug 日志详细，release 精简，上线切release

	r := gin.New()
	r.Use(middleware.RecoveryMiddleware(logger.Log))
	r.Use(gin.Logger())

	svc := service.NewUserService(db)
	artSvc := service.NewArticleService(db)
	LikeSvc := service.NewLike(db, noticeSvc)
	commentSvc := service.NewComments(db)

	jwtSvc := service.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpireHours)
	h := handler.NewUserHandler(svc, jwtSvc)
	h2 := handler.NewArticleHandler(artSvc, commentSvc, LikeSvc)

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	api := r.Group("/api/v1")
	{
		api.POST("/register", h.Register)
		api.POST("/login", h.Login)
		api.GET("/articles", h2.List)
		api.GET("/articles/:id", h2.Detail)
		api.GET("/articles/:id/comments", h2.ListByArticle)

	}
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware([]byte(cfg.JWT.Secret)))
	{
		protected.GET("/me", h.Me)
		protected.POST("/articles", h2.Create)
		protected.PUT("/articles/:id", h2.Update)
		protected.DELETE("/articles/:id", h2.Delete)
		protected.POST("/articles/:id/comments", h2.CreateComments)
		protected.POST("/articles/:id/like", h2.CreateLike)
		protected.DELETE("/articles/:id/like", h2.DeleteLike)
	}
	return r
}
