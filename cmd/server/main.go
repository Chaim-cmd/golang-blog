package main

import (
	"log"

	"github.com/Chaim-cmd/golang-blog/internal/config"
	"github.com/Chaim-cmd/golang-blog/internal/database"
	"github.com/Chaim-cmd/golang-blog/internal/logger"
	"github.com/Chaim-cmd/golang-blog/internal/router"
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

	logger.Log.Info("server starting",
		zap.String("mode", cfg.Server.Mode),
		zap.String("port", cfg.Server.Port))
	r := router.NewRouter(cfg, db)

	log.Printf("[%s] listening on %s", cfg.App.Name, cfg.Server.Port)

	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("服务退出：%v", err)
	}
}
