package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"syscall"
	"time"

	"os/signal"

	"github.com/Chaim-cmd/golang-blog/internal/config"
	"github.com/Chaim-cmd/golang-blog/internal/database"
	"github.com/Chaim-cmd/golang-blog/internal/logger"
	"github.com/Chaim-cmd/golang-blog/internal/router"
	"github.com/Chaim-cmd/golang-blog/internal/service"
	"go.uber.org/zap"
)

func main() {
	//装配配置
	cfg, err := config.Load("internal/config/config.yaml")
	if err != nil {
		log.Fatalf("启动失败：%v", err)
	}
	db, err := database.InitMySQL(cfg.Database)
	if err != nil {
		log.Fatalf("数据库初始化失败：%v", err)

	}

	if err := logger.Init(cfg.Log.Level); err != nil {
		log.Fatalf("init logger : %v", err)
	}
	defer logger.Log.Sync() // 退出前把缓冲刷盘

	noticeSvc := service.NewNotices(db)

	logger.Log.Info("server starting",
		zap.String("mode", cfg.Server.Mode),
		zap.String("port", cfg.Server.Port))
	r := router.NewRouter(cfg, db, noticeSvc)

	log.Printf("[%s] listening on %s", cfg.App.Name, cfg.Server.Port)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	srv := http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("服务退出：%v", err)
		}
	}()

	//阻塞
	<-ctx.Done()
	logger.Log.Info("shutting down...")

	//收尾
	shutDownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutDownCtx); err != nil {
		logger.Log.Error("http shutdown", zap.Error(err))
	}

	//关闭通知
	noticeSvc.Shutdown()
	logger.Log.Info("bye")
}
