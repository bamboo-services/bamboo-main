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

-- xf_system 系统表
CREATE TABLE xf_system (
    system_uuid UUID PRIMARY KEY,
    system_name VARCHAR NOT NULL,
    system_value JSONB NOT NULL,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 表注释
COMMENT ON TABLE xf_system IS '系统表';
COMMENT ON COLUMN xf_system.system_uuid IS '系统UUID';
COMMENT ON COLUMN xf_system.system_name IS '系统名称';
COMMENT ON COLUMN xf_system.system_value IS '系统值';
COMMENT ON COLUMN xf_system.created_at IS '创建时间';
COMMENT ON COLUMN xf_system.updated_at IS '更新时间';
