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

-- xf_link_group 友链分组表
CREATE TABLE xf_link_group
(
    group_uuid UUID NOT NULL PRIMARY KEY,
    group_name VARCHAR(50) NOT NULL,
    group_desc TEXT,
    group_order INT DEFAULT 0 NOT NULL,
    group_created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    group_updated_at TIMESTAMP DEFAULT NOW() NOT NULL
);

-- 表注释
COMMENT ON TABLE xf_link_group IS '友链分组表';
COMMENT ON COLUMN xf_link_group.group_uuid IS '分组唯一标识符';
COMMENT ON COLUMN xf_link_group.group_name IS '分组名称';
COMMENT ON COLUMN xf_link_group.group_desc IS '分组描述';
COMMENT ON COLUMN xf_link_group.group_order IS '分组排序';
COMMENT ON COLUMN xf_link_group.group_created_at IS '分组创建时间';
COMMENT ON COLUMN xf_link_group.group_updated_at IS '分组更新时间';
