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

	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
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
	publicLogic         *logic.PublicLogic
	oauthLogic          *bSdkLogic.BusinessLogic
}

type handler struct {
	name    string
	log     *xLog.LogNamedLogger
	service *service
}

type IHandler interface {
	~struct {
		name    string
		log     *xLog.LogNamedLogger
		service *service
	}
}

func NewHandler[T IHandler](ctx context.Context, handlerName string) *T {
	return &T{
		name: handlerName,
		log:  xLog.WithName(xLog.NamedCONT, handlerName),
		service: &service{
			authLogic:           logic.NewAuthLogic(ctx),
			infoLogic:           logic.NewInfoLogic(ctx),
			linkLogic:           logic.NewLinkLogic(ctx),
			linkColorLogic:      logic.NewLinkColorLogic(ctx),
			linkGroupLogic:      logic.NewLinkGroupLogic(ctx),
			sponsorChannelLogic: logic.NewSponsorChannelLogic(ctx),
			sponsorRecordLogic:  logic.NewSponsorRecordLogic(ctx),
			publicLogic:         logic.NewPublicLogic(ctx),
			oauthLogic:          bSdkLogic.NewBusiness(ctx),
		},
	}
}

type AuthHandler handler
type InfoHandler handler
type LinkHandler handler
type LinkColorHandler handler
type LinkGroupHandler handler
type SponsorChannelHandler handler
type SponsorRecordHandler handler
type PublicHandler handler
