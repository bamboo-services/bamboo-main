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
	// SystemFieldsRedisKey
	// 系统字段Redis键，%s 为系统 name 字段
	SystemFieldsRedisKey = "system:%s"

	// UserTokenRedisKey
	// 用户令牌Redis键，用于存储用户令牌信息
	UserTokenRedisKey = "user:token"
)

const (
	// NextSendEmailTimeRedisKey
	// 下次发送邮件时间Redis键，%s 为用户邮箱地址
	NextSendEmailTimeRedisKey = "email:next:%s"

	// SendEmailSmsRedisKey
	// 发送邮件验证码Redis键，%s 为用户邮箱地址
	SendEmailSmsRedisKey = "email:sms:%s"
)

const (
	// LinkGroupByUuidRedisKey
	// 链接组Redis键，%s 为链接组 UUID
	LinkGroupByUuidRedisKey = "link:group:%s"

	// LinkGroupByNameRedisKey
	// 链接组名称Redis键，%s 为链接组名称
	LinkGroupByNameRedisKey = "link:group:name:%s"

	// LinkGroupListRedisKey
	// 链接组列表Redis键，用于存储所有链接组的列表
	LinkGroupListRedisKey = "link:group:list"

	// LinkGroupListHasSearchRedisKey
	// 链接组列表搜索Redis键，%s 分别为 `页码:大小:关键词`
	LinkGroupListHasSearchRedisKey = "link:group:list:search:%d:%d:%s"
)

const (
	// LinkColorByUuidRedisKey
	// 链接颜色通过 UUID 获取的 Redis 键
	LinkColorByUuidRedisKey = "link:color:uuid:%s"

	// LinkColorByNameRedisKey
	// 链接颜色通过名称获取的 Redis 键
	LinkColorByNameRedisKey = "link:color:name:%s"

	// LinkColorListRedisKey
	// 链接颜色列表的 Redis 键
	LinkColorListRedisKey = "link:color:list"

	// LinkColorListHasSearchRedisKey
	// 链接颜色分页列表带搜索的 Redis 键，格式：页码:每页数量:搜索词
	LinkColorListHasSearchRedisKey = "link:color:list:%d:%d:%s"
)

const (
	// LinkContextByUuidRedisKey
	// 用于通过 UUID 存储获取的友链信息
	//
	// 1. %s 为友链 UUID
	LinkContextByUuidRedisKey = "link:context:uuid:%s"

	// LinkContextListRedisKey
	// 用于通过名称存储获取的友链信息
	LinkContextListRedisKey = "link:context:list"

	// LinkContextByUrlRedisKey
	// 用于通过 URL 存储获取的友链信息
	//
	// 1. %s 为经过 MD5 加密的友链 URL
	LinkContextByUrlRedisKey = "link:context:url:%s"

	// LinkContextListSearchRedisKey
	// 用于存储友链查询列表信息
	//
	// 1. %d 为页码
	// 2. %d 为每页数量
	// 3. %s 为搜索词
	LinkContextListSearchRedisKey = "link:context:list:search:%d:%d:%s" // %d:页码, %d:每页数量, %s:搜索词
)
