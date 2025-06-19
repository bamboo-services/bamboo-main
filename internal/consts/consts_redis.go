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

package consts

const (
	SystemFieldsRedisKey = "system:%s"  // SystemFieldsRedisKey 系统字段Redis键，%s 为系统 name 字段
	UserTokenRedisKey    = "user:token" // UserTokenRedisKey 用户令牌Redis键，用于存储用户令牌信息
)
