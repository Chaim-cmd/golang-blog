package database

import (
	"fmt"
	"time"

	"github.com/Chaim-cmd/golang-blog/internal/config"
	"github.com/Chaim-cmd/golang-blog/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 初始化 MYSQL 连接
func InitMySQL(cfg config.DatabaseConfig) (*gorm.DB, error) {
	//DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("连接 Mysql 失败:%w", err)
	}

	//拿到底层的sqlDB 对象
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("获取底层数据库实例失败: %w", err)
	}
	//连接池三件套
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	//启动时先 Ping 探活
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("Mysql ping 失败：%w", err)
	}

	//自动迁移 AutoMigrate
	if err := db.AutoMigrate(&model.User{}); err != nil {
		return nil, fmt.Errorf("迁移 users 表失败：%w", err)
	}

	return db, nil
}
