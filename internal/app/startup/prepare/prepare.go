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

package prepare

import (
	"context"

	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	"gorm.io/gorm"
)

// Prepare 封装预置数据初始化上下文，供子步骤复用日志、DB 连接与请求上下文。
//
// 由 [DefaultData] 在 xOptionDB.WithPrepare 回调中构造，不对外暴露构造函数。
type Prepare struct {
	log *xLog.LogNamedLogger
	db  *gorm.DB
	ctx context.Context
}

// DefaultData 构造业务预置数据的 xOptionDB.PrepareFunc。
//
// 该回调在 xOption 数据库节点完成 AutoMigrate 后执行，依次初始化：
//   - 默认管理员用户（幂等，基于 system.admin.id 标记判定是否已初始化）
//   - 默认站点信息配置（幂等，按 key 跳过已存在条目）
//
// 任一子步骤返回 error 会中断启动流程。失败信息以 Warn 记录后向上传递。
func DefaultData(ctx context.Context, db *gorm.DB) error {
	log := xLog.WithName(xLog.NamedINIT)
	p := &Prepare{log: log, db: db, ctx: ctx}

	if err := p.prepareDefaultUser(); err != nil {
		return err
	}

	if err := p.prepareDefaultInfo(); err != nil {
		return err
	}

	log.Info(ctx, "业务预置数据初始化完成")
	return nil
}
