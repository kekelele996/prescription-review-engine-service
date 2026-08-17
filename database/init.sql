-- 处方智能审核引擎数据库初始化脚本
-- 业务表结构由后端 GORM AutoMigrate 自动创建（backend/internal/database/database.go）。
-- 本脚本仅确保扩展与基础环境就绪，幂等可重复执行。
CREATE EXTENSION IF NOT EXISTS pgcrypto;
