package prepare

import (
	"context"

	xLog "github.com/bamboo-services/bamboo-base-go/log"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/utility/ctxutil"
	"gorm.io/gorm"
)

type Prepare struct {
	log *xLog.LogNamedLogger
	db  *gorm.DB
	ctx context.Context
}

func New(log *xLog.LogNamedLogger, ctx context.Context) *Prepare {
	return &Prepare{
		log: log,
		db:  xCtxUtil.MustGetDB(ctx),
		ctx: ctx,
	}
}

func (p *Prepare) Prepare() error {
	if err := p.prepareDefaultUser(); err != nil {
		return err
	}

	if err := p.prepareDefaultInfo(); err != nil {
		return err
	}

	p.log.Info(p.ctx, "业务预置数据初始化完成")
	return nil
}
