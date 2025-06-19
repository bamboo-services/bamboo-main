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

-- xf_logs 操作日志表
CREATE TABLE xf_logs
(
    log_uuid UUID NOT NULL PRIMARY KEY,
    log_type SMALLINT DEFAULT 0 NOT NULL,
    log_content TEXT NOT NULL,
    log_created_at TIMESTAMP DEFAULT NOW() NOT NULL
);

-- 表注释
COMMENT ON TABLE xf_logs IS '操作日志表';
COMMENT ON COLUMN xf_logs.log_uuid IS '日志唯一标识符';
COMMENT ON COLUMN xf_logs.log_type IS '日志类型';
COMMENT ON COLUMN xf_logs.log_content IS '日志内容';
COMMENT ON COLUMN xf_logs.log_created_at IS '日志创建时间';