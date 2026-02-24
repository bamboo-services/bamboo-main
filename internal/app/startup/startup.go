package startup

import (
	"context"

	"github.com/bamboo-services/bamboo-main/pkg/constants"
	bSdkStartup "github.com/phalanx-labs/beacon-sso-sdk/startup"

	xCtx "github.com/bamboo-services/bamboo-base-go/defined/context"
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

func Init() (context.Context, []xRegNode.RegNodeList) {
	businessReg := newInit()
	var regNode []xRegNode.RegNodeList

	regNode = append(regNode, xRegNode.RegNodeList{Key: constants.ContextCustomConfig, Node: businessReg.configInit})
	regNode = append(regNode, xRegNode.RegNodeList{Key: xCtx.DatabaseKey, Node: businessReg.databaseInit})
	regNode = append(regNode, xRegNode.RegNodeList{Key: xCtx.RedisClientKey, Node: businessReg.nosqlInit})
	regNode = append(regNode, xRegNode.RegNodeList{Key: xCtx.Exec, Node: businessReg.businessDataPrepare})
	regNode = append(regNode, bSdkStartup.NewOAuthConfig()...)

	return businessReg.ctx, regNode
}
