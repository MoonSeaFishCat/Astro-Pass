-- Casbin策略初始化脚本
-- 注意：这些策略会通过Casbin API自动创建，这里仅供参考

-- Casbin使用casbin_rule表存储策略
-- 表结构由gorm-adapter自动创建

-- 示例：为admin角色添加所有权限
-- 这些策略应该通过API调用创建，而不是直接插入数据库
-- 示例代码：
-- permissionService.AssignPermissionToRole("admin", "user", "read")
-- permissionService.AssignPermissionToRole("admin", "user", "write")
-- permissionService.AssignPermissionToRole("admin", "user", "delete")


