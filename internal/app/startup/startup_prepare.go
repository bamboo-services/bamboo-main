package startup

import (
	"context"

	"github.com/bamboo-services/bamboo-main/internal/app/startup/prepare"

	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
)

func (r *reg) businessDataPrepare(ctx context.Context) (any, error) {
	log := xLog.WithName(xLog.NamedINIT)

	if err := prepare.New(log, ctx).Prepare(); err != nil {
		return nil, err
	}

	return nil, nil
}
