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

package startup

import (
	"bamboo-main/internal/model/base"

	xInit "github.com/bamboo-services/bamboo-base-go/init"
	"github.com/bwmarrin/snowflake"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Reg 是一个类型别名，用于表示注册操作的整数类型。
type Reg struct {
	DB            *gorm.DB           // 数据库连接实例
	Rdb           *redis.Client      // Redis 客户端实例
	Config        *base.BambooConfig // 客制化配置文件
	Serv          *xInit.Reg         // 初始服务注册
	SnowflakeNode *snowflake.Node    // Snowflake 节点实例
}

// New 创建并返回一个未初始化的 `Reg` 实例。
//
// 该函数仅分配内存并返回 `Reg` 类型的初始值，
// 需要调用者进一步初始化相关字段。
//
// 返回值:
//   - `*Reg`: 返回一个新的 `Reg` 实例。
func New(serv *xInit.Reg) *Reg {
	return &Reg{Serv: serv}
}

// Register 创建一个经过完整初始化的 `Reg` 实例。
//
// 该函数负责调用一系列初始化方法，包括配置文件加载（ConfigInit）、数据库连接（DatabaseInit）、
// Redis 连接池初始化（RedisInit）以及系统上下文配置（SystemContextInit）。
// 最终返回一个已准备好的 `Reg` 实例用于后续操作。
//
// 参数说明:
//   - serv: 表示服务注册的实例 `*xInit.Reg`。
//
// 返回值:
//   - 返回初始化完毕的 `*Reg` 实例。
func Register(serv *xInit.Reg) *Reg {
	reg := New(serv)

	reg.ConfigInit()       // 配置文件增量处理
	reg.SnowflakeInit()    // 初始化 Snowflake 节点
	reg.DatabaseInit()     // 初始化数据库连接
	reg.RedisInit()        // 初始化 Redis 连接池
	reg.DatabaseUserInit() // 初始化系统用户

	// 初始化系统上下文
	reg.SystemContextInit()

	return reg
}
