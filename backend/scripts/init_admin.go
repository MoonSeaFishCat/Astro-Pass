package main

import (
	"fmt"
	"log"

	"astro-pass/internal/config"
	"astro-pass/internal/database"
	"astro-pass/internal/models"
	"astro-pass/internal/services"
	"astro-pass/internal/utils"
)

func main() {
	fmt.Println("=== 星穹通行证 - 初始化管理员账户 ===")

	// 加载配置
	config.Load()

	// 初始化数据库
	if err := database.Init(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 检查是否已存在管理员
	var existingAdmin models.User
	if err := database.DB.Where("username = ?", "admin").First(&existingAdmin).Error; err == nil {
		fmt.Println("⚠️  管理员账户已存在！")
		fmt.Printf("用户名: %s\n", existingAdmin.Username)
		fmt.Printf("邮箱: %s\n", existingAdmin.Email)
		fmt.Println("\n如需重置密码，请手动修改数据库或删除现有管理员账户。")
		return
	}

	// 创建管理员账户
	password := "Admin@123456"
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		log.Fatalf("密码加密失败: %v", err)
	}

	admin := &models.User{
		UUID:         utils.GenerateUUID(),
		Username:     "admin",
		Email:        "admin@astro-pass.local",
		PasswordHash: passwordHash,
		Nickname:     "系统管理员",
		Status:       "active",
		EmailVerified: true,
	}

	if err := database.DB.Create(admin).Error; err != nil {
		log.Fatalf("创建管理员失败: %v", err)
	}

	fmt.Println("✅ 管理员账户创建成功！")
	fmt.Println("\n账户信息：")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("用户名: %s\n", admin.Username)
	fmt.Printf("邮箱: %s\n", admin.Email)
	fmt.Printf("密码: %s\n", password)
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// 创建admin角色（如果不存在）
	var adminRole models.Role
	if err := database.DB.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		adminRole = models.Role{
			Name:        "admin",
			DisplayName: "系统管理员",
			Description: "拥有所有权限的系统管理员角色",
		}
		if err := database.DB.Create(&adminRole).Error; err != nil {
			log.Printf("创建admin角色失败: %v", err)
		} else {
			fmt.Println("✅ admin角色创建成功")
		}
	}

	// 为管理员分配admin角色
	if err := database.DB.Model(&admin).Association("Roles").Append(&adminRole); err != nil {
		log.Printf("分配角色失败: %v", err)
	} else {
		fmt.Println("✅ 已为管理员分配admin角色")
	}

	// 创建基础权限
	permissions := []models.Permission{
		{Name: "user:read", DisplayName: "查看用户", Resource: "user", Action: "read", Description: "查看用户信息"},
		{Name: "user:write", DisplayName: "管理用户", Resource: "user", Action: "write", Description: "创建和编辑用户"},
		{Name: "user:delete", DisplayName: "删除用户", Resource: "user", Action: "delete", Description: "删除用户"},
		{Name: "role:read", DisplayName: "查看角色", Resource: "role", Action: "read", Description: "查看角色信息"},
		{Name: "role:write", DisplayName: "管理角色", Resource: "role", Action: "write", Description: "创建和编辑角色"},
		{Name: "permission:read", DisplayName: "查看权限", Resource: "permission", Action: "read", Description: "查看权限信息"},
		{Name: "permission:write", DisplayName: "管理权限", Resource: "permission", Action: "write", Description: "创建和编辑权限"},
		{Name: "audit:read", DisplayName: "查看审计日志", Resource: "audit", Action: "read", Description: "查看审计日志"},
		{Name: "backup:manage", DisplayName: "备份管理", Resource: "backup", Action: "manage", Description: "管理数据库备份"},
		{Name: "config:manage", DisplayName: "配置管理", Resource: "config", Action: "manage", Description: "管理系统配置"},
	}

	fmt.Println("\n创建基础权限...")
	for _, perm := range permissions {
		var existing models.Permission
		if err := database.DB.Where("name = ?", perm.Name).First(&existing).Error; err != nil {
			if err := database.DB.Create(&perm).Error; err != nil {
				log.Printf("创建权限 %s 失败: %v", perm.Name, err)
			} else {
				fmt.Printf("✅ 创建权限: %s\n", perm.DisplayName)
			}
		}
	}

	// 初始化权限服务并为admin角色分配所有权限
	permissionService, err := services.NewPermissionService()
	if err != nil {
		log.Printf("初始化权限服务失败: %v", err)
	} else {
		fmt.Println("\n为admin角色分配权限...")
		for _, perm := range permissions {
			if err := permissionService.AssignPermissionToRole("admin", perm.Resource, perm.Action); err != nil {
				log.Printf("分配权限 %s 失败: %v", perm.Name, err)
			}
		}
		fmt.Println("✅ 权限分配完成")
	}

	// 初始化系统配置
	configService := services.NewSystemConfigService()
	if err := configService.InitDefaultConfigs(); err != nil {
		log.Printf("初始化系统配置失败: %v", err)
	} else {
		fmt.Println("✅ 系统配置初始化完成")
	}

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🎉 初始化完成！")
	fmt.Println("\n⚠️  重要提示：")
	fmt.Println("1. 请立即登录并修改默认密码")
	fmt.Println("2. 生产环境请删除或禁用此默认账户")
	fmt.Println("3. 建议创建新的管理员账户后删除此账户")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
