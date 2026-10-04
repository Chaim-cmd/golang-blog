package main

import (
	"log"

	"github.com/Chaim-cmd/golang-blog/internal/config"
	"github.com/Chaim-cmd/golang-blog/internal/router"
)

func main() {
	//装配配置
	cfg, err := config.Load("internal/config/config.yaml")
	if err != nil {
		log.Fatalf("启动失败：%v", err)
	}
	r := router.NewRouter(cfg)

	log.Printf("[%s] listening on %s", cfg.App.Name, cfg.Server.Port)

	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("服务退出：%v", err)
	}
}
