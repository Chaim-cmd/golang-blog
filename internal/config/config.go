package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config 顶层配置 ： 以后加mysql,redis 就往里面加挂子结构,mapstructure标签是viper解析yaml用的
type Config struct {
	App struct {
		Name string `mapstructure:"name"`
	} `mapstructure:"app"` //修复拼写 mapstructure

	Server struct {
		Port string `mapstructure:"port"`
		Mode string `mapstructure:"mode"`
	} `mapstructure:"server"` //修复拼写 mapstructure
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"`
	Host            string `mapstructure:"Host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	DBName          string `mapstructure:"dbname"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}
type JWTConfig struct {
	Secret      string `mapstructure:"secret"`
	ExpireHours int    `mapstructure:"expire_hours"`
}

func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	//默认值兜底
	v.SetDefault("app.name", "golang-blog")
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.mode", "debug")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置失败：%w", err)
	}
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("读取配置失败：%w", err)
	}
	return &cfg, nil
}
