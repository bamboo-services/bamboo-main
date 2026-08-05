-- --------------------------------------------------------------------------------
-- 颜色分级整理：12 基础色（普通 type=0）+ 4 高级色（type=1）+ 删除 20 个相近/无用色
-- --------------------------------------------------------------------------------
-- 背景：bm_link_color 表为旧系统迁移而来，36 色全部只有 primary_color（sub/hover 为空），
--       且遗留 type 列（bigint 默认 0）但 Go 实体未映射。本脚本重新启用 type 列：
--       type=0 普通色（单色渲染）、type=1 高级色（三色渐变渲染）。
--       炫彩为内置虚拟色（ID=1，不落库），不受本脚本影响。
--
-- 注意：表名为单数（bm_link_color / bm_link_friend）。
-- 执行后需清理颜色缓存：redis-cli --scan --pattern 'bm:link:color:*' | xargs -r redis-cli del
-- --------------------------------------------------------------------------------

BEGIN;

-- 0) 确保 type 列存在（旧脚本 migration_link_color_drop_type.sql 曾删列；AutoMigrate 只加列不删列）
ALTER TABLE bm_link_color ADD COLUMN IF NOT EXISTS type bigint NOT NULL DEFAULT 0;

-- 1) 引用检查：先查后删，确认被删颜色下的友链引用（预期仅蛋白石 1 条）
SELECT c.id, c.name, COUNT(f.id) AS link_count
FROM bm_link_color c
LEFT JOIN bm_link_friend f ON f.color_id = c.id
WHERE c.name IN ('银色','铜色','香槟金','桃红','杏色','珊瑚红','鲑鱼红','番茄红','栗色','酒红',
                 '翡翠绿','紫水晶','黄玉','绿松石','海蓝宝石','橄榄石','蛋白石','珍珠','月光石','钻石')
GROUP BY c.id, c.name
HAVING COUNT(f.id) > 0;

-- 2) 处理引用：被删颜色下友链 color_id 置空（恢复默认颜色）
UPDATE bm_link_friend SET color_id = NULL
WHERE color_id IN (SELECT id FROM bm_link_color WHERE name IN ('银色','铜色','香槟金','桃红','杏色','珊瑚红','鲑鱼红','番茄红','栗色','酒红',
                 '翡翠绿','紫水晶','黄玉','绿松石','海蓝宝石','橄榄石','蛋白石','珍珠','月光石','钻石'));

-- 3) 12 基础色：type=0，清空 sub/hover 保持纯单色
UPDATE bm_link_color SET type = 0, sub_color = NULL, hover_color = NULL
WHERE name IN ('红色','橙色','黄色','绿色','青色','蓝色','紫色','粉色','黑色','白色','灰色','棕色');

-- 4) 4 高级色：type=1 + 补齐三色（sub=浅色渐变起点，hover=深色强调）
UPDATE bm_link_color SET type = 1, primary_color='#FFD700', sub_color='#FFF3B0', hover_color='#D4A017' WHERE name='金色';
UPDATE bm_link_color SET type = 1, primary_color='#B76E79', sub_color='#F3C2C6', hover_color='#8E4A54' WHERE name='玫瑰金'; -- 修正色值，避免与粉色 #FFC0CB 重复
UPDATE bm_link_color SET type = 1, primary_color='#E0115F', sub_color='#F79AC0', hover_color='#A80E46' WHERE name='红宝石';
UPDATE bm_link_color SET type = 1, primary_color='#0F52BA', sub_color='#A5C8E8', hover_color='#0A3B8A' WHERE name='蓝宝石';

-- 5) 删除其余 20 色
DELETE FROM bm_link_color WHERE name IN ('银色','铜色','香槟金','桃红','杏色','珊瑚红','鲑鱼红','番茄红','栗色','酒红',
                 '翡翠绿','紫水晶','黄玉','绿松石','海蓝宝石','橄榄石','蛋白石','珍珠','月光石','钻石');

COMMIT;
