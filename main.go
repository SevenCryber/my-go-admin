package main

import (
	"fmt"
	"github.com/SevenCryber/my-go-admin/config"
	"github.com/SevenCryber/my-go-admin/config/task"
	"github.com/SevenCryber/my-go-admin/initialize/dal"
	"github.com/SevenCryber/my-go-admin/router"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"strconv"
)

func main() {

	// 初始化配置文件
	config.InitConfig()
	task.InitTaskConfig()
	// 初始化全局 logger
	if err := config.InitLogger(); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		return
	}
	defer func() {
		_ = config.Logger.Sync() // 确保所有日志都被刷新到磁盘
	}()

	// dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := config.Data.Mysql.Username + ":" + config.Data.Mysql.Password + "@tcp(" + config.Data.Mysql.Host + ":" + strconv.Itoa(config.Data.Mysql.Port) + ")/" + config.Data.Mysql.Database + "?charset=" + config.Data.Mysql.Charset + "&parseTime=True&loc=Local"

	// 初始化数据访问层
	dal.InitDal(&dal.Config{
		GomrConfig: &dal.GomrConfig{
			Dialector: mysql.Open(dsn),
			Opts: &gorm.Config{
				SkipDefaultTransaction: true, // 跳过默认事务
				NamingStrategy: schema.NamingStrategy{
					SingularTable: true,
				},
				Logger: logger.New(log.Default(), logger.Config{
					LogLevel: logger.Silent, // 不打印日志
					//LogLevel: logger.Error, // 打印错误日志
				}),
			},
			MaxOpenConns: config.Data.Mysql.MaxOpenConns,
			MaxIdleConns: config.Data.Mysql.MaxIdleConns,
		},
		// RedisConfig: &dal.RedisConfig{
		// 	Host:     config.Data.Redis.Host,
		// 	Port:     config.Data.Redis.Port,
		// 	Database: config.Data.Redis.Database,
		// 	Password: config.Data.Redis.Password,
		// },
	})

	// 设置模式
	gin.SetMode(config.Data.App.Server.Mode)

	// 初始化gin
	server := gin.New()

	// 注册恢复中间件
	server.Use(gin.Recovery())

	// 注册路由
	router.ApiRegister(server)
	config.Logger.Info("Application started")
	server.Run(":" + strconv.Itoa(config.Data.App.Server.Port))
}
