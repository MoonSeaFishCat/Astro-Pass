package main

import (
	"log"

	"astro-pass/internal/config"
	"astro-pass/internal/database"
	"astro-pass/internal/routes"
	"astro-pass/internal/services"
	"astro-pass/internal/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	config.Load()

	// 验证配置
	if err := config.Cfg.Validate(); err != nil {
		log.Fatalf("配置验证失败: %v", err)
	}

	// 初始化日志
	utils.InitLogger()
	utils.Info("星穹通行证服务启动中...")

	// 初始化数据库
	if err := database.Init(); err != nil {
		utils.Error("数据库初始化失败: %v", err)
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 初始化Redis（可选）
	if err := services.InitRedis(); err != nil {
		utils.Warn("Redis初始化失败（将继续运行，但缓存功能不可用）: %v", err)
	} else {
		utils.Info("Redis初始化完成")
	}

	// 初始化RSA密钥对（用于JWT签名）
	if err := utils.InitRSAKeys(); err != nil {
		utils.Error("RSA密钥初始化失败: %v", err)
		log.Fatalf("RSA密钥初始化失败: %v", err)
	}
	utils.Info("RSA密钥初始化完成")

	// 初始化权限服务（数据库初始化后）
	// 注意：这里只是预初始化，实际使用时会延迟初始化
	utils.Info("数据库初始化完成")

	// 初始化系统配置
	configService := services.NewSystemConfigService()
	if err := configService.InitDefaultConfigs(); err != nil {
		utils.Warn("初始化系统配置失败: %v", err)
	} else {
		utils.Info("系统配置初始化完成")
	}

	// 启动定时任务调度器
	scheduler := utils.NewScheduler()
	scheduler.Start()
	defer scheduler.Stop()
	utils.Info("定时任务调度器已启动")

	// 设置Gin模式
	if config.Cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 设置路由
	router := routes.SetupRoutes()

	// 启动服务器
	addr := config.Cfg.Server.Host + ":" + config.Cfg.Server.Port
	utils.Info("星穹通行证服务启动在 %s", addr)
	log.Printf("星穹通行证服务启动在 %s", addr)
	if err := router.Run(addr); err != nil {
		utils.Error("服务器启动失败: %v", err)
		log.Fatalf("服务器启动失败: %v", err)
	}
}

