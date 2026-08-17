-- 处方智能审核引擎数据库迁移脚本
-- 说明：业务表由后端 GORM AutoMigrate 自动创建（见 backend/internal/database/database.go），
-- 本脚本为数据库级初始化（扩展、索引建议），幂等可重复执行。

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- 处方号唯一索引（AutoMigrate 已创建，此处为冗余保障）
CREATE UNIQUE INDEX IF NOT EXISTS uk_prescription_no ON prescriptions (prescription_no);
-- 审核状态统计索引
CREATE INDEX IF NOT EXISTS idx_prescriptions_status ON prescriptions (status);
CREATE INDEX IF NOT EXISTS idx_prescriptions_overall_risk ON prescriptions (overall_risk);
CREATE INDEX IF NOT EXISTS idx_review_items_rule_type ON review_items (rule_type);
