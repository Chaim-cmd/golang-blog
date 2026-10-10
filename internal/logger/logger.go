package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log 全局 logger 全局不够优雅,先跑起来
// 优雅姿势看 RecoveryMiddleware：logger 当参数传进去，而不是全局抓
var Log *zap.Logger

// Init 按 level 初始化生产级 logger
func Init(level string) error {
	var lv zapcore.Level
	// UnmarshalText :把 "info" 之类字符串解析成 zap 的级别枚举,拼错返回err
	if err := lv.UnmarshalText([]byte(level)); err != nil {
		return err
	}
	cfg := zap.NewProductionConfig() // JSON encode, 日志平台原生格式
	cfg.Level = zap.NewAtomicLevelAt(lv)
	// 模拟是时间戳换成 IS08601:人类可读,平台可识别
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	var err error
	Log, err = cfg.Build()
	return err
}
