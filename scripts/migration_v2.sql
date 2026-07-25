-- --------------------------------------------------------------------------------
-- v2 BREAKING CHANGE 迁移脚本
-- --------------------------------------------------------------------------------
-- 用途：在升级到 v2 实体结构（嵌入 xModels.BaseEntity + SnowflakeID + timestamptz）前
--       清空旧表，让 GORM AutoMigrate 在服务启动时按新结构重建。
--
-- 警告：此脚本为破坏性操作，执行后旧数据无法恢复。
--       执行前请自行导出/备份需要保留的业务数据。
--
-- 表名前缀 "bm_" 来自 xOptionDB.WithTablePrefix("bm_")，GORM 默认复数蛇形。
-- DROP 顺序按外键反向依赖，避免约束冲突。
-- --------------------------------------------------------------------------------

BEGIN;

DROP TABLE IF EXISTS bm_sponsor_records;
DROP TABLE IF EXISTS bm_sponsor_channels;
DROP TABLE IF EXISTS bm_systems;
DROP TABLE IF EXISTS bm_system_logs;
DROP TABLE IF EXISTS bm_link_friends;
DROP TABLE IF EXISTS bm_link_colors;
DROP TABLE IF EXISTS bm_link_groups;
DROP TABLE IF EXISTS bm_system_users;

COMMIT;

-- --------------------------------------------------------------------------------
-- 重建工作交由服务启动时的 GORM AutoMigrate 完成，无需在此手写 CREATE TABLE。
-- 启动服务即可按 v2 结构自动建表，并经 prepare.DefaultData 重新注入种子数据。
-- --------------------------------------------------------------------------------
