// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

package logic

import (
	"context"

	apiDashboard "github.com/bamboo-services/bamboo-main/api/dashboard"
	"github.com/bamboo-services/bamboo-main/internal/repository"
	"github.com/bamboo-services/bamboo-main/pkg/constants"

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/major/utility/context"
)

type dashboardRepo struct {
	link *repository.LinkRepo
}

// DashboardLogic 仪表盘业务逻辑
type DashboardLogic struct {
	logic
	repo dashboardRepo
}

// NewDashboardLogic 创建 DashboardLogic 实例，从上下文获取数据库与缓存并初始化友链仓储依赖。
func NewDashboardLogic(ctx context.Context) *DashboardLogic {
	db := xCtxUtil.MustGetDB(ctx)
	m := xCtxUtil.MustGetCacheManager(ctx)

	return &DashboardLogic{
		logic: logic{
			db:    db,
			cache: m,
			log:   xLog.WithName(xLog.NamedLOGC, "DashboardLogic"),
		},
		repo: dashboardRepo{link: repository.NewLinkRepo(db, m)},
	}
}

// Stats 获取仪表盘统计数据：友链计数与最近待审核申请
func (d *DashboardLogic) Stats(ctx context.Context) (*apiDashboard.StatsResponse, *xError.Error) {
	total, xErr := d.repo.link.CountByStatus(ctx, -1, nil)
	if xErr != nil {
		return nil, xErr
	}
	pending, xErr := d.repo.link.CountByStatus(ctx, constants.LinkStatusPending, nil)
	if xErr != nil {
		return nil, xErr
	}
	approved, xErr := d.repo.link.CountByStatus(ctx, constants.LinkStatusApproved, nil)
	if xErr != nil {
		return nil, xErr
	}

	// 最近 5 条待审核申请（按创建时间倒序）
	pendingStatus := constants.LinkStatusPending
	recentLinks, _, xErr := d.repo.link.List(ctx, &repository.FriendQuery{
		Page:       1,
		PageSize:   5,
		LinkStatus: &pendingStatus,
		SortBy:     "created_at",
		SortOrder:  "desc",
	}, nil)
	if xErr != nil {
		return nil, xErr
	}

	recent := make([]apiDashboard.RecentApplicationItem, 0, len(recentLinks))
	for _, link := range recentLinks {
		avatar := ""
		if link.Avatar != nil {
			avatar = *link.Avatar
		}
		recent = append(recent, apiDashboard.RecentApplicationItem{
			ID:        link.ID,
			Name:      link.Name,
			Avatar:    avatar,
			URL:       link.URL,
			CreatedAt: link.CreatedAt,
		})
	}

	return &apiDashboard.StatsResponse{
		TotalLinks:         total,
		PendingLinks:       pending,
		ApprovedLinks:      approved,
		RecentApplications: recent,
	}, nil
}
