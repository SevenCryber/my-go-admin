package config

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"time"
)

// Logger 全局变量声明
var Logger *zap.Logger

// 确保日志目录和当天的日志文件存在，并返回日志文件路径
func ensureLogFilePath() (string, error) {
	logDir := "log"                          // 日志主目录
	today := time.Now().Format("2006-01-02") // 格式化为 YYYY-MM-DD 形式的字符串
	dirPath := filepath.Join(logDir, today)  // 拼接成完整的路径

	// 如果目录不存在，则创建它
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create log directory: %w", err)
	}

	// 创建日志文件名，例如：2024-12-28.log
	fileName := filepath.Join(dirPath, fmt.Sprintf("%s.log", today))

	// 检查文件是否存在，如果不存在则创建
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		file, err := os.Create(fileName)
		if err != nil {
			return "", fmt.Errorf("failed to create log file: %w", err)
		}
		defer file.Close()
	}

	return fileName, nil
}

// 初始化 Zap Logger 并设置全局变量
func InitLogger() error {
	// 获取日志文件路径
	logFilePath, err := ensureLogFilePath()
	if err != nil {
		return err
	}

	// 打开或创建日志文件
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	// 定义编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder          // 设置时间格式
	encoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder // 小写且带颜色的日志级别

	// 创建核心组件
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // 或者使用 zapcore.NewConsoleEncoder(encoderConfig) 对于更友好的控制台输出
		zapcore.AddSync(file),
		zap.InfoLevel, // 设置最低日志级别
	)

	// 构建 logger
	Logger = zap.New(core, zap.AddCaller(), zap.Development()) // 开发模式下显示调用者信息
	zap.ReplaceGlobals(Logger)                                 // 替换全局默认logger

	return nil
}
