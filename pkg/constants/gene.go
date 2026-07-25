// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

package constants

import (
	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
)

// 业务级 Gene 基因常量（区间 16-63，由 xSnowflake 预留给业务自定义）。
//
// 基因类型用于雪花 ID 生成时按实体类型做命名空间隔离，
// 不同实体的 ID 不会在全局冲突，便于跨表识别 ID 归属。
// 系统级实体（User/Config/Log）复用 xSnowflake 内置 GeneUser/GeneConfig/GeneLog。
const (
	GeneLink           xSnowflake.Gene = 16 // 友情链接实体
	GeneLinkGroup      xSnowflake.Gene = 17 // 友链分组实体
	GeneLinkColor      xSnowflake.Gene = 18 // 友链颜色实体
	GeneSponsorChannel xSnowflake.Gene = 19 // 赞助渠道实体
	GeneSponsorRecord  xSnowflake.Gene = 20 // 赞助记录实体
)
