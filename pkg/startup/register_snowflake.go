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
	xConsts "github.com/bamboo-services/bamboo-base-go/constants"
	"github.com/bwmarrin/snowflake"
)

// SnowflakeInit 初始化Snowflake算法节点，用于生成分布式唯一ID。
// 如果节点初始化失败，程序将触发 panic 并终止运行。
func (r *Reg) SnowflakeInit() {
	log := r.Serv.Logger.Named(xConsts.LogINIT).Sugar()
	nodeID := r.Config.Snowflake.NodeID
	if nodeID == nil {
		nodeID = new(int64)
		*nodeID = 1
		log.Warn("雪花节点ID未配置，使用默认值 1")
	} else {
		log.Infof("雪花节点ID: %d", *nodeID)
	}
	newNode, err := snowflake.NewNode(*nodeID)
	if err != nil {
		panic(err)
	}
	r.SnowflakeNode = newNode
}
