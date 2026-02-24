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
