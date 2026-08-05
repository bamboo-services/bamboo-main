-- --------------------------------------------------------------------------------
-- 配色调整：16 色统一为「竹林淡彩」淡雅三色（主色 / 副色 / 悬停色）
-- --------------------------------------------------------------------------------
-- 背景：原色值为标准纯色（如 #FF0000），过于"正"，未贴合本站「竹林水墨」
--       主题（淡绿宣纸 + 墨色 + 竹绿，低饱和高明度）。本脚本按主题调淡：
--       - primary 主色：色相正确但淡雅（亮度 L≈60-79%，降饱和）
--       - sub 副色：极浅，渐变起点（L≈85-90%）
--       - hover 悬停色：比主色深一档，悬停强调（L≈50-55%）
--       普通色与高级色均配置完整三色。
--
-- 表名单数：bm_link_color。执行后需清理颜色缓存。
-- --------------------------------------------------------------------------------

BEGIN;

-- 普通 12 色：type=0，淡雅三色
UPDATE bm_link_color SET primary_color='#E0655A', sub_color='#F3B4AE', hover_color='#BE4639' WHERE name='红色';
UPDATE bm_link_color SET primary_color='#E8904E', sub_color='#F4CDA6', hover_color='#C86C2F' WHERE name='橙色';
UPDATE bm_link_color SET primary_color='#DCCB52', sub_color='#F1E7A9', hover_color='#B3A02E' WHERE name='黄色';
UPDATE bm_link_color SET primary_color='#7CAE87', sub_color='#C4E0C9', hover_color='#4E8C5A' WHERE name='绿色';
UPDATE bm_link_color SET primary_color='#5DB5AE', sub_color='#BCE2DE', hover_color='#36928A' WHERE name='青色';
UPDATE bm_link_color SET primary_color='#6E9ED6', sub_color='#C2D7EF', hover_color='#3D71AC' WHERE name='蓝色';
UPDATE bm_link_color SET primary_color='#A58AD0', sub_color='#D9CEEB', hover_color='#7E60AC' WHERE name='紫色';
UPDATE bm_link_color SET primary_color='#E8A9B5', sub_color='#F6DBE0', hover_color='#C97F8E' WHERE name='粉色';
UPDATE bm_link_color SET primary_color='#4A4A4A', sub_color='#9A9A9A', hover_color='#2E2E2E' WHERE name='黑色';
UPDATE bm_link_color SET primary_color='#F2EFE4', sub_color='#FAF8F0', hover_color='#D8D3C0' WHERE name='白色';
UPDATE bm_link_color SET primary_color='#A9A9A1', sub_color='#D7D7D0', hover_color='#7F7F77' WHERE name='灰色';
UPDATE bm_link_color SET primary_color='#B0805E', sub_color='#DCC5A9', hover_color='#8A623C' WHERE name='棕色';

-- 高级 4 色：type=1，精细淡雅三色
UPDATE bm_link_color SET primary_color='#D9B45E', sub_color='#F1E1B0', hover_color='#A9852F' WHERE name='金色';
UPDATE bm_link_color SET primary_color='#D8A1AA', sub_color='#F1D3D8', hover_color='#B1727E' WHERE name='玫瑰金';
UPDATE bm_link_color SET primary_color='#C76B82', sub_color='#E9B3C0', hover_color='#9A4458' WHERE name='红宝石';
UPDATE bm_link_color SET primary_color='#7B9ED2', sub_color='#BCCFE9', hover_color='#4B6FA0' WHERE name='蓝宝石';

COMMIT;
