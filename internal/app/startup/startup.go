/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

package startup

import (
	"context"

	"github.com/bamboo-services/bamboo-main/pkg/constants"
	bSdkStartup "github.com/phalanx-labs/beacon-sso-sdk/startup"

	xRegNode "github.com/bamboo-services/bamboo-base-go/major/register/node"
)

type reg struct {
	ctx context.Context
}

func newInit() *reg {
	return &reg{
		ctx: context.Background(),
	}
}

// Init 返回 startup 自定义节点列表与根上下文。
//
// 自 v1.0.4 起，数据库与缓存由 main.go 的 xOption 声明式配置装配，
// 本函数仅注册框架未覆盖的业务节点：
//   - Email 业务配置（供 worker_mail 读取）
//   - SSO SDK 依赖节点（OAuth 配置 + gRPC 客户端）
func Init() (context.Context, []xRegNode.RegNodeList) {
	businessReg := newInit()
	var regNode []xRegNode.RegNodeList

	regNode = append(regNode, xRegNode.RegNodeList{Key: constants.ContextCustomConfig, Node: businessReg.emailConfigInit})
	regNode = append(regNode, bSdkStartup.NewStartupConfig()...)

	return businessReg.ctx, regNode
}
