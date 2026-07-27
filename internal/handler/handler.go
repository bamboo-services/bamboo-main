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

package handler

import (
	"context"

	"github.com/bamboo-services/bamboo-main/internal/logic"
	bSdkLogic "github.com/phalanx-labs/beacon-sso-sdk/logic"
)

type service struct {
	authLogic           *logic.AuthLogic
	infoLogic           *logic.InfoLogic
	linkLogic           *logic.LinkLogic
	linkColorLogic      *logic.LinkColorLogic
	linkGroupLogic      *logic.LinkGroupLogic
	sponsorChannelLogic *logic.SponsorChannelLogic
	sponsorRecordLogic  *logic.SponsorRecordLogic
	dashboardLogic      *logic.DashboardLogic
	publicLogic         *logic.PublicLogic
	oauthLogic          *bSdkLogic.BusinessLogic
}

type handler struct {
	service *service
}

// IHandler handler 泛型约束，限定由 NewHandler 构造的处理器结构
type IHandler interface {
	~struct {
		service *service
	}
}

// NewHandler 泛型构造器，装配全部 logic 依赖并返回具体 handler 实例
func NewHandler[T IHandler](ctx context.Context) *T {
	return &T{
		service: &service{
			authLogic:           logic.NewAuthLogic(ctx),
			infoLogic:           logic.NewInfoLogic(ctx),
			linkLogic:           logic.NewLinkLogic(ctx),
			linkColorLogic:      logic.NewLinkColorLogic(ctx),
			linkGroupLogic:      logic.NewLinkGroupLogic(ctx),
			sponsorChannelLogic: logic.NewSponsorChannelLogic(ctx),
			sponsorRecordLogic:  logic.NewSponsorRecordLogic(ctx),
			dashboardLogic:      logic.NewDashboardLogic(ctx),
			publicLogic:         logic.NewPublicLogic(ctx),
			oauthLogic:          bSdkLogic.NewBusiness(ctx),
		},
	}
}

// AuthHandler 认证接口处理器
type AuthHandler handler

// InfoHandler 站点信息接口处理器
type InfoHandler handler

// LinkHandler 友情链接接口处理器
type LinkHandler handler

// LinkColorHandler 友链颜色接口处理器
type LinkColorHandler handler

// LinkGroupHandler 友链分组接口处理器
type LinkGroupHandler handler

// SponsorChannelHandler 赞助渠道接口处理器
type SponsorChannelHandler handler

// SponsorRecordHandler 赞助记录接口处理器
type SponsorRecordHandler handler

// DashboardHandler 仪表盘接口处理器
type DashboardHandler handler

// PublicHandler 公开接口处理器
type PublicHandler handler
