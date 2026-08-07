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
	"strings"

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/major/utility/context"
	apiAuth "github.com/bamboo-services/bamboo-main/api/auth"
	"github.com/bamboo-services/bamboo-main/internal/models/base"
	"github.com/bamboo-services/bamboo-main/internal/repository"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
)

type systemUserRepo struct {
	user   *repository.SystemUserRepo
	system *repository.SystemRepo
}

// SystemUserLogic 系统用户管理业务逻辑
//
// 收口管理端用户查询与状态管理：列表批量填充管理员身份计算字段（防 N+1），
// 状态变更时保护系统唯一管理员（以 bm_system 配置表的 system.admin.id 为权威判定）。
type SystemUserLogic struct {
	logic
	repo systemUserRepo
}

// NewSystemUserLogic 创建 SystemUserLogic 实例，从上下文获取数据库与缓存管理器并初始化仓储依赖。
func NewSystemUserLogic(ctx context.Context) *SystemUserLogic {
	db := xCtxUtil.MustGetDB(ctx)
	m := xCtxUtil.MustGetCacheManager(ctx)

	return &SystemUserLogic{
		logic: logic{
			db:    db,
			cache: m,
			log:   xLog.WithName(xLog.NamedLOGC, "SystemUserLogic"),
		},
		repo: systemUserRepo{
			user:   repository.NewSystemUserRepo(db, m),
			system: repository.NewSystemRepo(db, m),
		},
	}
}

// ListUsers 分页查询系统用户，批量填充管理员身份计算字段。
//
// 管理员身份由 bm_system 配置表的 system.admin.id 唯一标记，一次性读取配置
// 后逐条比对填充，避免对每名用户触发 IsAdmin 单查造成的 N+1 查询。
func (s *SystemUserLogic) ListUsers(ctx context.Context, req *apiAuth.UserQueryRequest) (*apiAuth.UserListResponse, *xError.Error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 10
	}

	users, total, xErr := s.repo.user.List(ctx, &repository.UserQuery{
		Page:      req.Page,
		PageSize:  req.PageSize,
		Keyword:   req.Keyword,
		Status:    req.Status,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	})
	if xErr != nil {
		return nil, xErr
	}

	// 读取一次管理员 ID 批量填充 IsAdmin（配置缺失/异常按非管理员处理，不阻断主流程）
	adminID := s.adminID(ctx)
	for i := range users {
		users[i].IsAdmin = adminID != "" && users[i].ID.String() == adminID
	}

	result := base.NewPaginationResponse(users, req.Page, req.PageSize, total)
	return &apiAuth.UserListResponse{PaginationResponse: *result}, nil
}

// UpdateUserStatus 启用/禁用指定用户。
//
// 系统唯一管理员不可被禁用：以 bm_system 配置表的 system.admin.id 为权威判定
// （不可依赖缓存读回的用户实体 IsAdmin 字段——该字段经 JSON 序列化后恒为 false）。
// 状态更新复用 UpdateFieldsByID，内部已处理缓存失效与回读重建。
func (s *SystemUserLogic) UpdateUserStatus(ctx context.Context, userID xSnowflake.SnowflakeID, req *apiAuth.UserStatusRequest) *xError.Error {
	user, found, xErr := s.repo.user.GetByID(ctx, userID)
	if xErr != nil {
		return xErr
	}
	if !found {
		return xError.NewError(ctx, xError.NotFound, "用户不存在", false)
	}

	// 保护系统唯一管理员（配置判定）：仅拦截「禁用」方向，「启用」管理员无害
	if req.Status == 0 && user.ID.String() == s.adminID(ctx) {
		return xError.NewError(ctx, xError.OperationInvalid, "不能禁用系统管理员", false)
	}

	_, xErr = s.repo.user.UpdateFieldsByID(ctx, userID, map[string]any{"status": req.Status})
	return xErr
}

// adminID 读取系统唯一管理员的用户 ID。
//
// 配置缺失、值为空串或读取失败时返回空字符串，调用方按无管理员处理，不影响主流程。
func (s *SystemUserLogic) adminID(ctx context.Context) string {
	config, found, xErr := s.repo.system.GetByKey(ctx, constants.KeySystemAdminID)
	if xErr != nil || !found || config.Value == nil {
		return ""
	}
	return strings.TrimSpace(*config.Value)
}
