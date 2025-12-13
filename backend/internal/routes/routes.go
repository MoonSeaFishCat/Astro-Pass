package routes

import (
	"astro-pass/internal/controllers"
	"astro-pass/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	// 全局中间件
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.SecurityHeadersMiddleware())
	router.Use(middleware.MetricsMiddleware())

	// CORS配置
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // 前端地址
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

	// 速率限制（全局限制）
	router.Use(middleware.RateLimitMiddleware(120)) // 每分钟120次

	// 健康检查
	healthController := controllers.NewHealthController()
	router.GET("/health", healthController.Health)
	router.GET("/ready", healthController.Ready)

	// Prometheus指标
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// API路由组
	api := router.Group("/api")
	{
		// 认证相关路由（适度的速率限制）
		authController := controllers.NewAuthController()
		auth := api.Group("/auth")
		auth.Use(middleware.RateLimitMiddleware(30)) // 每分钟30次
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
			auth.POST("/refresh", authController.RefreshToken)
			auth.GET("/profile", middleware.AuthMiddleware(), authController.GetProfile)
		// 忘记密码和重置密码在 UserController 中
		userController := controllers.NewUserController()
		auth.POST("/forgot-password", userController.ForgotPassword)
		auth.POST("/reset-password", userController.ResetPassword)
		}

		// 用户相关路由
		userController := controllers.NewUserController()
		user := api.Group("/user")
		user.Use(middleware.AuthMiddleware())
		{
			user.PUT("/profile", userController.UpdateProfile)
			user.POST("/change-password", userController.ChangePassword)
		}

		// 管理员用户管理路由
		adminUser := api.Group("/admin/users")
		adminUser.Use(middleware.AuthMiddleware())
		adminUser.Use(middleware.PermissionMiddleware("user", "read")) // 需要用户读取权限
		{
			adminUser.GET("", userController.GetAllUsers)
			adminUser.GET("/stats", userController.GetUserStats)
			adminUser.GET("/:id", userController.GetUser)
			adminUser.PUT("/:id", userController.UpdateUser)
			adminUser.DELETE("/:id", userController.DeleteUser)
			adminUser.POST("/:id/roles", userController.AssignRoleToUser)
			adminUser.DELETE("/:id/roles", userController.RemoveRoleFromUser)
		}

		// 权限管理路由
		permissionController := controllers.NewPermissionController()
		permission := api.Group("/permission")
		permission.Use(middleware.AuthMiddleware())
		{
			permission.POST("/assign-role", permissionController.AssignRole)
			permission.GET("/roles", permissionController.GetUserRoles)
			permission.POST("/role", permissionController.CreateRole)
			permission.POST("/permission", permissionController.CreatePermission)
			permission.POST("/role/:role/permission", permissionController.AssignPermissionToRole)
		}

		// 管理员权限管理路由
		adminPermission := api.Group("/admin")
		adminPermission.Use(middleware.AuthMiddleware())
		adminPermission.Use(middleware.PermissionMiddleware("role", "read")) // 需要角色读取权限
		{
			adminPermission.GET("/roles", permissionController.GetAllRoles)
			adminPermission.PUT("/roles/:id", permissionController.UpdateRole)
			adminPermission.DELETE("/roles/:id", permissionController.DeleteRole)
			adminPermission.GET("/permissions", permissionController.GetAllPermissions)
			adminPermission.PUT("/permissions/:id", permissionController.UpdatePermission)
			adminPermission.DELETE("/permissions/:id", permissionController.DeletePermission)
		}

		// 审计日志路由
		auditController := controllers.NewAuditController()
		audit := api.Group("/audit")
		audit.Use(middleware.AuthMiddleware())
		{
			audit.GET("/logs", auditController.GetAuditLogs)
			audit.GET("/log/:id", auditController.GetAuditLog)
		}

		// 会话管理路由
		sessionController := controllers.NewSessionController()
		session := api.Group("/session")
		session.Use(middleware.AuthMiddleware())
		{
			session.GET("/list", sessionController.GetSessions)
			session.DELETE("/:id", sessionController.RevokeSession)
			session.DELETE("/all", sessionController.RevokeAllSessions)
		}

		// MFA相关路由
		mfaController := controllers.NewMFAController()
		mfa := api.Group("/mfa")
		mfa.Use(middleware.AuthMiddleware())
		{
			mfa.POST("/generate", mfaController.GenerateTOTP)
			mfa.POST("/enable", mfaController.EnableMFA)
			mfa.POST("/disable", mfaController.DisableMFA)
			mfa.GET("/recovery-codes", mfaController.GetRecoveryCodes)
		}

		// OAuth2/OIDC路由
		oauth2Controller := controllers.NewOAuth2Controller()
		tokenController := controllers.NewTokenController()
		oauth2 := api.Group("/oauth2")
		{
			oauth2.GET("/authorize", middleware.AuthMiddleware(), oauth2Controller.Authorize)
			oauth2.POST("/token", oauth2Controller.Token)
			oauth2.GET("/userinfo", oauth2Controller.UserInfo)
			oauth2.GET("/jwks", tokenController.GetJWKS)
			oauth2.POST("/revoke", tokenController.RevokeToken)
			oauth2.POST("/introspect", tokenController.IntrospectToken)
		}

		// OAuth2客户端管理路由
		oauth2ClientController := controllers.NewOAuth2ClientController()
		oauth2Clients := api.Group("/oauth2/clients")
		oauth2Clients.Use(middleware.AuthMiddleware())
		{
			oauth2Clients.POST("", oauth2ClientController.CreateClient)
			oauth2Clients.GET("", oauth2ClientController.GetUserClients)
			oauth2Clients.DELETE("/:id", oauth2ClientController.RevokeClient)
		}

		// 授权同意路由
		consentController := controllers.NewConsentController()
		consent := api.Group("/oauth2/consent")
		{
			consent.GET("/info", consentController.GetConsentInfo)
			consent.POST("/approve", middleware.AuthMiddleware(), consentController.ApproveConsent)
			consent.POST("/deny", middleware.AuthMiddleware(), consentController.DenyConsent)
			consent.GET("/list", middleware.AuthMiddleware(), consentController.GetUserConsents)
			consent.DELETE("/:client_id", middleware.AuthMiddleware(), consentController.RevokeConsent)
		}

		// OIDC发现端点
		router.GET("/.well-known/openid-configuration", tokenController.GetOpenIDConfiguration)

		// WebAuthn路由
		webauthnController := controllers.NewWebAuthnController()
		if webauthnController != nil {
			webauthn := api.Group("/webauthn")
			{
				// 注册流程
				webauthn.POST("/register/begin", middleware.AuthMiddleware(), webauthnController.BeginRegistration)
				webauthn.POST("/register/finish", middleware.AuthMiddleware(), webauthnController.FinishRegistration)
				
				// 登录流程
				webauthn.POST("/login/begin", webauthnController.BeginLogin)
				webauthn.POST("/login/finish", webauthnController.FinishLogin)
				
				// 凭证管理（需要认证）
				webauthn.GET("/credentials", middleware.AuthMiddleware(), webauthnController.GetCredentials)
				webauthn.DELETE("/credentials", middleware.AuthMiddleware(), webauthnController.DeleteCredential)
			}
		}

		// 邮箱验证路由
		emailVerificationController := controllers.NewEmailVerificationController()
		emailVerification := api.Group("/email-verification")
		emailVerification.Use(middleware.AuthMiddleware())
		{
			emailVerification.POST("/send", emailVerificationController.SendVerificationEmail)
		}
		api.POST("/email-verification/verify", emailVerificationController.VerifyEmail)

		// 通知路由
		notificationController := controllers.NewNotificationController()
		notification := api.Group("/notifications")
		notification.Use(middleware.AuthMiddleware())
		{
			notification.GET("", notificationController.GetNotifications)
			notification.PUT("/:id/read", notificationController.MarkAsRead)
			notification.PUT("/read-all", notificationController.MarkAllAsRead)
			notification.DELETE("/:id", notificationController.DeleteNotification)
		}

		// 社交登录路由
		socialAuthController := controllers.NewSocialAuthController()
		socialAuth := api.Group("/auth/social")
		{
			socialAuth.GET("/github/url", socialAuthController.GetGitHubAuthURL)
			socialAuth.POST("/github/callback", socialAuthController.HandleGitHubCallback)
		}

		// 备份管理路由（需要管理员权限）
		backupController := controllers.NewBackupController()
		backup := api.Group("/admin/backup")
		backup.Use(middleware.AuthMiddleware())
		backup.Use(middleware.PermissionMiddleware("backup", "manage"))
		{
			backup.POST("", backupController.CreateBackup)
			backup.GET("", backupController.GetBackupList)
			backup.GET("/stats", backupController.GetBackupStats)
			backup.DELETE("/:id", backupController.DeleteBackup)
			backup.POST("/:id/restore", backupController.RestoreBackup)
			backup.GET("/:id/download", backupController.DownloadBackup)
			backup.POST("/clean", backupController.CleanOldBackups)
		}

		// 系统配置路由（需要管理员权限）
		systemConfigController := controllers.NewSystemConfigController()
		systemConfig := api.Group("/admin/config")
		systemConfig.Use(middleware.AuthMiddleware())
		systemConfig.Use(middleware.PermissionMiddleware("config", "manage"))
		{
			systemConfig.GET("", systemConfigController.GetAllConfigs)
			systemConfig.GET("/category/:category", systemConfigController.GetConfigsByCategory)
			systemConfig.PUT("", systemConfigController.UpdateConfig)
			systemConfig.GET("/backup", systemConfigController.GetBackupConfig)
			systemConfig.PUT("/backup", systemConfigController.UpdateBackupConfig)
			systemConfig.GET("/export", systemConfigController.ExportConfigs)
			systemConfig.POST("/import", systemConfigController.ImportConfigs)
		}

		// SSO单点登出路由
		sloController := controllers.NewSLOController()
		sso := api.Group("/sso")
		{
			sso.POST("/logout", middleware.AuthMiddleware(), sloController.InitiateLogout)
			sso.GET("/sessions", middleware.AuthMiddleware(), sloController.GetUserSessions)
			sso.GET("/logout/:request_id/status", middleware.AuthMiddleware(), sloController.GetLogoutStatus)
			sso.POST("/logout/callback", sloController.HandleLogoutCallback)
		}

		// 管理员SSO管理路由
		adminSSO := api.Group("/admin/sso")
		adminSSO.Use(middleware.AuthMiddleware())
		adminSSO.Use(middleware.PermissionMiddleware("sso", "manage"))
		{
			adminSSO.POST("/users/:user_id/revoke-sessions", sloController.AdminRevokeUserSessions)
		}

		// OIDC登出端点
		oidc := api.Group("/oidc")
		{
			oidc.GET("/logout", sloController.GetOIDCLogout)
		}

		// SAML路由
		samlController := controllers.NewSAMLController()
		saml := api.Group("/saml")
		{
			saml.GET("/metadata", samlController.GetMetadata)
			saml.GET("/sso", samlController.HandleSSO)
			saml.POST("/sso", samlController.HandleSSO)
			saml.GET("/login-complete", middleware.AuthMiddleware(), samlController.HandleSAMLLogin)
		}

		// 管理员SAML配置路由
		adminSAML := api.Group("/admin/saml")
		adminSAML.Use(middleware.AuthMiddleware())
		adminSAML.Use(middleware.PermissionMiddleware("saml", "manage"))
		{
			adminSAML.POST("/configs", samlController.CreateSAMLConfig)
			adminSAML.GET("/configs", samlController.GetSAMLConfigs)
			adminSAML.PUT("/configs/:id", samlController.UpdateSAMLConfig)
			adminSAML.DELETE("/configs/:id", samlController.DeleteSAMLConfig)
		}
	}

	return router
}

