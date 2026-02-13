package startup

import (
	"context"

	"github.com/bamboo-services/bamboo-main/internal/task"
	"github.com/bamboo-services/bamboo-main/pkg/constants"

	xCtx "github.com/bamboo-services/bamboo-base-go/context"
	xRegNode "github.com/bamboo-services/bamboo-base-go/register/node"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/utility/ctxutil"
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

	regNode = append(regNode, xRegNode.RegNodeList{Key: xCtx.DatabaseKey, Node: businessReg.databaseInit})
	regNode = append(regNode, xRegNode.RegNodeList{Key: xCtx.RedisClientKey, Node: businessReg.nosqlInit})
	regNode = append(regNode, xRegNode.RegNodeList{Key: constants.ContextMailWorker, Node: businessReg.mailWorkerInit})
	regNode = append(regNode, xRegNode.RegNodeList{Key: xCtx.Exec, Node: businessReg.businessDataPrepare})

	return businessReg.ctx, regNode
}

func GetMailWorker(ctx context.Context) *task.MailWorker {
	worker, err := xCtxUtil.Get[*task.MailWorker](ctx, constants.ContextMailWorker)
	if err != nil {
		return nil
	}
	return worker
}
