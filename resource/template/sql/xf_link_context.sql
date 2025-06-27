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


-- xf_link_context 友链内容表
CREATE TABLE xf_link_context
(
    link_uuid UUID NOT NULL PRIMARY KEY,
    link_name VARCHAR(100) NOT NULL,
    link_url VARCHAR(255) NOT NULL,
    link_avatar VARCHAR(255),
    link_desc TEXT,
    link_email VARCHAR(100),
    link_group_uuid UUID NOT NULL,
    link_color_uuid UUID,
    link_order INT DEFAULT 0 NOT NULL,
    link_status SMALLINT DEFAULT 0 NOT NULL,  -- 0: 待审核, 1: 已通过, 2: 已失效
    link_apply_remark TEXT,                  -- 申请者备注
    link_review_remark TEXT,                 -- 审核备注
    link_created_at TIMESTAMP DEFAULT NOW() NOT NULL,
    link_updated_at TIMESTAMP DEFAULT NOW() NOT NULL,
    FOREIGN KEY (link_group_uuid) REFERENCES xf_link_group(group_uuid),
    FOREIGN KEY (link_color_uuid) REFERENCES xf_link_color(color_uuid)
);

-- 表注释
COMMENT ON TABLE xf_link_context IS '友链内容表';
COMMENT ON COLUMN xf_link_context.link_uuid IS '友链唯一标识符';
COMMENT ON COLUMN xf_link_context.link_name IS '友链名称';
COMMENT ON COLUMN xf_link_context.link_url IS '友链URL地址';
COMMENT ON COLUMN xf_link_context.link_avatar IS '友链头像URL';
COMMENT ON COLUMN xf_link_context.link_desc IS '友链描述';
COMMENT ON COLUMN xf_link_context.link_email IS '友链联系邮箱';
COMMENT ON COLUMN xf_link_context.link_group_uuid IS '所属分组ID';
COMMENT ON COLUMN xf_link_context.link_color_uuid IS '友链颜色ID';
COMMENT ON COLUMN xf_link_context.link_order IS '友链排序';
COMMENT ON COLUMN xf_link_context.link_status IS '友链状态（0: 待审核, 1: 已通过, 2: 已拒绝）';
COMMENT ON COLUMN xf_link_context.link_apply_remark IS '申请者备注';
COMMENT ON COLUMN xf_link_context.link_review_remark IS '审核备注';
COMMENT ON COLUMN xf_link_context.link_created_at IS '友链创建时间';
COMMENT ON COLUMN xf_link_context.link_updated_at IS '友链更新时间';
