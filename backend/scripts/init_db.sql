-- 星穹通行证数据库初始化脚本

-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS astro_pass CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE astro_pass;

-- 注意：表结构会通过GORM自动迁移创建
-- 这里只创建一些初始数据

-- 插入默认角色
INSERT IGNORE INTO roles (id, name, display_name, description, created_at, updated_at) VALUES
(1, 'admin', '管理员', '系统管理员，拥有所有权限', NOW(), NOW()),
(2, 'user', '普通用户', '普通用户角色', NOW(), NOW()),
(3, 'teacher', '教师', '教师角色', NOW(), NOW()),
(4, 'student', '学生', '学生角色', NOW(), NOW());

-- 插入默认权限
INSERT IGNORE INTO permissions (id, name, display_name, resource, action, description, created_at, updated_at) VALUES
(1, 'user:read', '查看用户', 'user', 'read', '查看用户信息', NOW(), NOW()),
(2, 'user:write', '编辑用户', 'user', 'write', '编辑用户信息', NOW(), NOW()),
(3, 'user:delete', '删除用户', 'user', 'delete', '删除用户', NOW(), NOW()),
(4, 'role:read', '查看角色', 'role', 'read', '查看角色信息', NOW(), NOW()),
(5, 'role:write', '编辑角色', 'role', 'write', '编辑角色信息', NOW(), NOW()),
(6, 'permission:read', '查看权限', 'permission', 'read', '查看权限信息', NOW(), NOW()),
(7, 'permission:write', '编辑权限', 'permission', 'write', '编辑权限信息', NOW(), NOW()),
(8, 'audit:read', '查看审计日志', 'audit', 'read', '查看审计日志', NOW(), NOW());

-- 为管理员角色分配所有权限
INSERT IGNORE INTO role_permissions (role_id, permission_id) VALUES
(1, 1), (1, 2), (1, 3), (1, 4), (1, 5), (1, 6), (1, 7), (1, 8);


