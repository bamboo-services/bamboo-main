/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2025 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

-- xf_link_color 友链颜色表
CREATE TABLE xf_link_color
(
    color_uuid UUID NOT NULL PRIMARY KEY,
    color_name VARCHAR(50) NOT NULL,
    color_value VARCHAR(20) NOT NULL,
    color_desc TEXT,
    color_created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    color_updated_at TIMESTAMP DEFAULT NOW() NOT NULL
);

-- 表注释
COMMENT ON TABLE xf_link_color IS '友链颜色表';
COMMENT ON COLUMN xf_link_color.color_uuid IS '颜色唯一标识符';
COMMENT ON COLUMN xf_link_color.color_name IS '颜色名称';
COMMENT ON COLUMN xf_link_color.color_value IS '颜色值（如HEX值：#FFFFFF）';
COMMENT ON COLUMN xf_link_color.color_desc IS '颜色描述';
COMMENT ON COLUMN xf_link_color.color_created_at IS '颜色创建时间';
COMMENT ON COLUMN xf_link_color.color_updated_at IS '颜色更新时间';
