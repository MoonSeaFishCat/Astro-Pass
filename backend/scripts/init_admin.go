package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"astro-pass/internal/config"
	"astro-pass/internal/database"
	"astro-pass/internal/models"
	"astro-pass/internal/utils"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	config.Load()

	// 验证配置
	if err := config.Cfg.Validate(); err != nil {
		log.Fatalf("配置验证失败: %v", err)
	}

	// 初始化数据库
	if err := database.Init(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 检查是否已存在管理员用户
	var adminUser models.User
	err := database.DB.Where("username = ?", "admin").First(&adminUser).Error
	if err == nil {
		fmt.Println("✓ 管理员账户已存在，跳过创建")
		os.Exit(0)
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Fatalf("检查管理员账户失败: %v", err)
	}

	// 创建默认管理员账户
	fmt.Println("正在创建默认管理员账户...")

	// 生成密码哈希
	passwordHash, err := utils.HashPassword("Admin@123456")
	if err != nil {
		log.Fatalf("密码加密失败: %v", err)
	}

	// 创建用户
	adminUser = models.User{
		UUID:         utils.GenerateUUID(),
		Username:     "admin",
		Email:        "admin@astro-pass.local",
		PasswordHash: passwordHash,
		Nickname:     "系统管理员",
		Status:       "active",
		EmailVerified: true, // 默认管理员邮箱已验证
	}

	if err := database.DB.Create(&adminUser).Error; err != nil {
		log.Fatalf("创建管理员账户失败: %v", err)
	}

	fmt.Printf("✓ 管理员账户创建成功 (ID: %d)\n", adminUser.ID)

	// 查找 admin 角色
	var adminRole models.Role
	if err := database.DB.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		// 如果角色不存在，创建它
		adminRole = models.Role{
			Name:        "admin",
			DisplayName: "管理员",
			Description: "系统管理员，拥有所有权限",
		}
		if err := database.DB.Create(&adminRole).Error; err != nil {
			log.Fatalf("创建管理员角色失败: %v", err)
		}
		fmt.Println("✓ 管理员角色创建成功")
	}

	// 为用户分配管理员角色
	if err := database.DB.Model(&adminUser).Association("Roles").Append(&adminRole); err != nil {
		log.Fatalf("分配管理员角色失败: %v", err)
	}

	fmt.Println("✓ 管理员角色分配成功")

	// 输出登录信息
	separator := strings.Repeat("=", 50)
	fmt.Println("\n" + separator)
	fmt.Println("🎉 默认管理员账户创建成功！")
	fmt.Println(separator)
	fmt.Println("用户名: admin")
	fmt.Println("邮箱: admin@astro-pass.local")
	fmt.Println("密码: Admin@123456")
	fmt.Println("\n⚠️  重要提示：")
	fmt.Println("   1. 首次登录后请立即修改密码")
	fmt.Println("   2. 建议启用 MFA 多因素认证")
	fmt.Println("   3. 生产环境请删除或修改此账户")
	fmt.Println(separator)
}

