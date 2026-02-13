package startup

import (
	"context"

	"github.com/bamboo-services/bamboo-main/internal/model/base"
	"github.com/bamboo-services/bamboo-main/internal/task"
	"github.com/bamboo-services/bamboo-main/pkg/constants"

	xCtxUtil "github.com/bamboo-services/bamboo-base-go/utility/ctxutil"
	"go.uber.org/zap"
)

func (r *reg) mailWorkerInit(ctx context.Context) (any, error) {
	config, err := xCtxUtil.Get[*base.BambooConfig](ctx, constants.ContextCustomConfig)
	if err != nil {
		return nil, err
	}

	rdb := xCtxUtil.MustGetRDB(ctx)

	worker := task.NewMailWorker(rdb, &config.Email, zap.NewNop().Sugar())
	worker.Start()
	return worker, nil
}
