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
	SystemFieldsRedisKey           = "system:%s"                             // SystemFieldsRedisKey 系统字段Redis键，%s 为系统 name 字段
	UserTokenRedisKey              = "user:token"                            // UserTokenRedisKey 用户令牌Redis键，用于存储用户令牌信息
	NextSendEmailTimeRedisKey      = "email:next:%s"                         // NextSendEmailTimeRedisKey 下次发送邮件时间Redis键，%s 为用户邮箱地址
	SendEmailSmsRedisKey           = "email:sms:%s"                          // SendEmailSmsRedisKey 发送邮件验证码Redis键，%s 为用户邮箱地址
	LinkGroupByUUIDRedisKey        = "link_group:group:%s"                   // LinkGroupByUUIDRedisKey 链接组Redis键，%s 为链接组 UUID
	LinkGroupByNameRedisKey        = "link_group:group:name:%s"              // LinkGroupByNameRedisKey 链接组名称Redis键，%s 为链接组名称
	LinkGroupListRedisKey          = "link_group:group:list"                 // LinkGroupListRedisKey 链接组列表Redis键，用于存储所有链接组的列表
	LinkGroupListHasSearchRedisKey = "link_group:group:list:search:%d:%d:%s" // LinkGroupListHasSearchRedisKey 链接组列表搜索Redis键，%s 分别为 `页码:大小:关键词`
	LinkColorByUUIDRedisKey        = "xm:link:color:uuid:%s"                 // LinkColorByUUIDRedisKey 链接颜色通过 UUID 获取的 Redis 键
	LinkColorByNameRedisKey        = "xm:link:color:name:%s"                 // LinkColorByNameRedisKey 链接颜色通过名称获取的 Redis 键
	LinkColorListRedisKey          = "xm:link:color:list"                    // LinkColorListRedisKey 链接颜色列表的 Redis 键
	LinkColorListHasSearchRedisKey = "xm:link:color:list:%d:%d:%s"           // LinkColorListHasSearchRedisKey 链接颜色分页列表带搜索的 Redis 键，格式：页码:每页数量:搜索词
)
