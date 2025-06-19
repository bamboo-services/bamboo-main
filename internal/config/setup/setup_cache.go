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

package setup

import (
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
)

// ProjectSetupCache
//
// 项目缓存初始化；
// 主要用于初始化 Redis 缓存适配器，确保项目中使用的缓存系统正常工作。
func (setup *projectSetup) ProjectSetupCache() {
	blog.BambooNotice(setup.ctx, "ProjectSetupCache", "项目缓存初始化")

	// 初始化 Redis 缓存
	getAdapter := gcache.NewAdapterRedis(g.Redis())
	g.DB().GetCache().SetAdapter(getAdapter)
}
