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

package logic

import (
	"context"
	"errors"
	"fmt"

	apiLink "github.com/bamboo-services/bamboo-main/api/link"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	logcHelper "github.com/bamboo-services/bamboo-main/internal/logic/helper"
	"github.com/bamboo-services/bamboo-main/internal/models/base"
	"github.com/bamboo-services/bamboo-main/internal/repository"
	"github.com/bamboo-services/bamboo-main/internal/service/screenshot"
	"github.com/bamboo-services/bamboo-main/pkg/constants"

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
	xUtil "github.com/bamboo-services/bamboo-base-go/common/utility"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/major/utility/context"
	xAsync "github.com/bamboo-services/bamboo-base-go/plugins/async"
	"gorm.io/gorm"
)

type linkRepo struct {
	link   *repository.LinkRepo
	user   *repository.SystemUserRepo
	group  *repository.LinkGroupRepo
	system *repository.SystemRepo
}

// LinkLogic 友情链接业务逻辑
type LinkLogic struct {
	logic
	repo linkRepo
}

// NewLinkLogic 创建 LinkLogic 实例，从上下文获取数据库与缓存并初始化友链仓储依赖。
func NewLinkLogic(ctx context.Context) *LinkLogic {
	db := xCtxUtil.MustGetDB(ctx)
	m := xCtxUtil.MustGetCacheManager(ctx)

	return &LinkLogic{
		logic: logic{
			db:    db,
			cache: m,
			log:   xLog.WithName(xLog.NamedLOGC, "LinkLogic"),
		},
		repo: linkRepo{
			link:   repository.NewLinkRepo(db, m),
			user:   repository.NewSystemUserRepo(db, m),
			group:  repository.NewLinkGroupRepo(db, m),
			system: repository.NewSystemRepo(db, m),
		},
	}
}

// Add 添加友情链接
func (l *LinkLogic) Add(ctx context.Context, req *apiLink.FriendAddRequest) (*entity.LinkFriend, *xError.Error) {
	// 创建友情链接实体
	link := &entity.LinkFriend{
		Name:        req.LinkName,
		URL:         req.LinkURL,
		Avatar:      xUtil.Ptr(req.LinkAvatar),
		RSS:         xUtil.Ptr(req.LinkRSS),
		Description: xUtil.Ptr(req.LinkDesc),
		Email:       xUtil.Ptr(req.LinkEmail),
		SortOrder:   req.LinkOrder,
		Status:      constants.LinkStatusPending, // 默认待审核
		IsFailure:   constants.LinkFailNormal,    // 默认正常
		Level:       req.LinkLevel,
		ApplyRemark: xUtil.Ptr(req.LinkApplyRemark),
	}

	// 设置ID外键
	if req.LinkGroupID.Provided() {
		link.GroupID = req.LinkGroupID.Value()
	}
	if req.LinkColorID.Provided() {
		link.ColorID = req.LinkColorID.Value()
	}

	_, xErr := l.repo.link.Create(ctx, link, nil)
	if xErr != nil {
		return nil, xErr
	}

	reloaded, found, xErr := l.repo.link.GetByID(ctx, link.ID, true, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "友情链接不存在", false)
	}

	// 发送邮件通知管理员（xAsync 解耦请求上下文，不阻断主流程）
	xAsync.Async(ctx, func(asyncCtx context.Context) {
		l.sendApplyNotification(asyncCtx, reloaded)
	}, xAsync.WithName("MAIL"))

	return reloaded, nil
}

// Update 更新友情链接
func (l *LinkLogic) Update(ctx context.Context, linkID xSnowflake.SnowflakeID, req *apiLink.FriendUpdateRequest) (*entity.LinkFriend, *xError.Error) {
	link, found, xErr := l.repo.link.GetByID(ctx, linkID, false, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "友情链接不存在", false)
	}

	// 直接更新实体字段
	if req.LinkName != "" {
		link.Name = req.LinkName
	}
	if req.LinkURL != "" {
		link.URL = req.LinkURL
	}
	if req.LinkAvatar != "" {
		link.Avatar = xUtil.Ptr(req.LinkAvatar)
	}
	if req.LinkRSS != "" {
		link.RSS = xUtil.Ptr(req.LinkRSS)
	}
	if req.LinkDesc != "" {
		link.Description = xUtil.Ptr(req.LinkDesc)
	}
	urlChanged := req.LinkURL != "" && req.LinkURL != link.URL
	if req.LinkEmail != "" {
		link.Email = xUtil.Ptr(req.LinkEmail)
	}
	if req.LinkGroupID.Provided() {
		link.GroupID = req.LinkGroupID.Value()
	}
	if req.LinkColorID.Provided() {
		link.ColorID = req.LinkColorID.Value()
	}
	if req.LinkOrder != nil {
		link.SortOrder = *req.LinkOrder
	}
	if req.LinkLevel != nil {
		link.Level = *req.LinkLevel
	}
	if req.LinkApplyRemark != "" {
		link.ApplyRemark = xUtil.Ptr(req.LinkApplyRemark)
	}

	_, xErr = l.repo.link.Save(ctx, link, nil)
	if xErr != nil {
		return nil, xErr
	}

	reloaded, found, xErr := l.repo.link.GetByID(ctx, linkID, true, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "友情链接不存在", false)
	}

	// URL 变更：旧截图与站点不匹配，触发重新截图
	if urlChanged {
		if manager := screenshot.GetManager(ctx); manager != nil {
			manager.Enqueue(linkID)
		}
	}

	return reloaded, nil
}

// Delete 删除友情链接
func (l *LinkLogic) Delete(ctx context.Context, linkID xSnowflake.SnowflakeID) *xError.Error {
	ok, xErr := l.repo.link.DeleteByID(ctx, linkID, nil)
	if xErr != nil {
		return xErr
	}
	if !ok {
		return xError.NewError(ctx, xError.NotFound, "友情链接不存在", false)
	}
	return nil
}

// Get 获取友情链接详情
func (l *LinkLogic) Get(ctx context.Context, linkID xSnowflake.SnowflakeID) (*entity.LinkFriend, *xError.Error) {
	link, found, xErr := l.repo.link.GetByID(ctx, linkID, true, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "友情链接不存在", false)
	}

	return link, nil
}

// List 获取友情链接列表
func (l *LinkLogic) List(ctx context.Context, req *apiLink.FriendQueryRequest) (*apiLink.FriendListResponse, *xError.Error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 10
	}

	links, total, xErr := l.repo.link.List(ctx, &repository.FriendQuery{
		Page:        req.Page,
		PageSize:    req.PageSize,
		LinkName:    req.LinkName,
		LinkStatus:  req.LinkStatus,
		LinkFail:    req.LinkFail,
		LinkAnomaly: req.LinkAnomaly,
		LinkGroupID: req.LinkGroupID,
		SortBy:      req.SortBy,
		SortOrder:   req.SortOrder,
	}, nil)
	if xErr != nil {
		return nil, xErr
	}

	// 附加计数：待审核 / 异常，供管理端入口徽章展示
	pendingCount, xErr := l.repo.link.CountByStatus(ctx, constants.LinkStatusPending, nil)
	if xErr != nil {
		return nil, xErr
	}
	anomalyCount, xErr := l.repo.link.CountAnomaly(ctx, nil)
	if xErr != nil {
		return nil, xErr
	}

	result := base.NewPaginationResponse(links, req.Page, req.PageSize, total)
	return &apiLink.FriendListResponse{
		PaginationResponse: *result,
		PendingCount:       pendingCount,
		AnomalyCount:       anomalyCount,
	}, nil
}

// UpdateStatus 更新友情链接状态
func (l *LinkLogic) UpdateStatus(ctx context.Context, linkID xSnowflake.SnowflakeID, req *apiLink.FriendStatusRequest) *xError.Error {
	// 先查询友链信息（用于发送邮件通知）
	link, found, xErr := l.repo.link.GetByID(ctx, linkID, false, nil)
	if xErr != nil {
		return xErr
	}
	if !found {
		return xError.NewError(ctx, xError.NotFound, "友情链接不存在", false)
	}

	ok, xErr := l.repo.link.UpdateStatusByID(ctx, linkID, req.LinkStatus, req.LinkReviewRemark, nil)
	if xErr != nil {
		return xErr
	}
	if !ok {
		return xError.NewError(ctx, xError.NotFound, "友情链接不存在", false)
	}

	// 发送审核结果邮件通知（xAsync 解耦请求上下文，不阻断主流程）
	xAsync.Async(ctx, func(asyncCtx context.Context) {
		l.sendStatusNotification(asyncCtx, link, req.LinkStatus, req.LinkReviewRemark)
	}, xAsync.WithName("MAIL"))

	// 审核通过：触发友链首次站点截图（入队后由截图 worker 串行处理）
	if req.LinkStatus == constants.LinkStatusApproved {
		if manager := screenshot.GetManager(ctx); manager != nil {
			manager.Enqueue(linkID)
		}
	}

	return nil
}

// UpdateFailStatus 更新友情链接失效状态。
//
// 标记失效时自动归入内置「已失效」分组（保留 ID）；恢复时若当前分组为已失效分组则清空为未分组。
func (l *LinkLogic) UpdateFailStatus(ctx context.Context, linkID xSnowflake.SnowflakeID, req *apiLink.FriendFailRequest) *xError.Error {
	invalidGroupID := xSnowflake.SnowflakeID(constants.BuiltinGroupInvalidID)
	ok, xErr := l.repo.link.UpdateFailureByID(ctx, linkID, req.LinkFail, req.LinkFailReason, &invalidGroupID, nil)
	if xErr != nil {
		return xErr
	}
	if !ok {
		return xError.NewError(ctx, xError.NotFound, "友情链接不存在", false)
	}

	return nil
}

// BuildSortAssignments 根据载荷顺序计算全局排序赋值（纯函数，无 ctx/DB，便于表驱动单测）
//
// 以 links 建立 id→现有 GroupID 映射；按 items 载荷顺序分配全局序号 0..N-1；
// GroupID 三态：item 显式提供则取提供值（null 即置未分组），否则保持原分组。
// 载荷含重复 ID 时返回 error，由调用方包装为参数错误。
func BuildSortAssignments(links []entity.LinkFriend, items []apiLink.FriendSortItem) ([]repository.SortAssignment, error) {
	existing := make(map[xSnowflake.SnowflakeID]*xSnowflake.SnowflakeID, len(links))
	for i := range links {
		existing[links[i].ID] = links[i].GroupID
	}

	assignments := make([]repository.SortAssignment, 0, len(items))
	seen := make(map[xSnowflake.SnowflakeID]struct{}, len(items))
	for order, item := range items {
		if _, dup := seen[item.ID]; dup {
			return nil, errors.New("items 含重复友链ID")
		}
		seen[item.ID] = struct{}{}

		groupID := existing[item.ID]
		if item.GroupID.Provided() {
			groupID = item.GroupID.Value()
		}

		assignments = append(assignments, repository.SortAssignment{
			ID:      item.ID,
			GroupID: groupID,
			Order:   order,
		})
	}

	return assignments, nil
}

// UpdateSort 批量更新友链全局排序与分组归属
//
// 事务内按载荷顺序重写全局 sort_order 为 0..N-1，并随载荷携带三态 group_id
// （省略=保持原组 / null=置未分组 / 值=移入该组），写后由 repo 逐条失效缓存。
func (l *LinkLogic) UpdateSort(ctx context.Context, req *apiLink.FriendSortRequest) *xError.Error {
	// 收集友链 ID 并校验全部存在
	idSet := make(map[xSnowflake.SnowflakeID]struct{}, len(req.Items))
	ids := make([]xSnowflake.SnowflakeID, 0, len(req.Items))
	for _, item := range req.Items {
		if _, ok := idSet[item.ID]; !ok {
			idSet[item.ID] = struct{}{}
			ids = append(ids, item.ID)
		}
	}
	links, xErr := l.repo.link.GetByIDs(ctx, ids, nil)
	if xErr != nil {
		return xErr
	}
	if len(links) != len(ids) {
		return xError.NewError(ctx, xError.NotFound, "友情链接不存在", false)
	}

	// 收集显式提供的目标分组 ID，校验全部存在且启用
	groupIDSet := make(map[xSnowflake.SnowflakeID]struct{})
	for _, item := range req.Items {
		if item.GroupID.Provided() && item.GroupID.Value() != nil {
			// 内置分组恒存在且启用（不落库），跳过数据库存在性校验
			if entity.IsBuiltinGroupID(*item.GroupID.Value()) {
				continue
			}
			groupIDSet[*item.GroupID.Value()] = struct{}{}
		}
	}
	if len(groupIDSet) > 0 {
		groupIDs := make([]xSnowflake.SnowflakeID, 0, len(groupIDSet))
		for gid := range groupIDSet {
			groupIDs = append(groupIDs, gid)
		}
		groups, xErr := l.repo.group.GetByIDs(ctx, groupIDs, nil)
		if xErr != nil {
			return xErr
		}
		if len(groups) != len(groupIDs) {
			return xError.NewError(ctx, xError.NotFound, "友链分组不存在或已禁用", false)
		}
		for i := range groups {
			if !groups[i].Status {
				return xError.NewError(ctx, xError.NotFound, "友链分组不存在或已禁用", false)
			}
		}
	}

	// 计算排序赋值（纯函数），重复 ID 等非法载荷包装为参数错误
	assignments, err := BuildSortAssignments(links, req.Items)
	if err != nil {
		return xError.NewError(ctx, xError.ParameterError, xError.ErrMessage(err.Error()), false, err)
	}

	// 事务内批量写入（关联引用由 repo.Save 经 Omit(clause.Associations) 收敛，见 LinkRepo.Save）
	return l.withTx(ctx, func(tx *gorm.DB) *xError.Error {
		return l.repo.link.UpdateSortAndPosition(ctx, assignments, tx)
	})
}

// ReScreenshot 手动触发友链重新截图（仅已通过且未失效的友链）
//
// 仅将任务加入截图队列，由截图 worker 串行处理，接口本身不等待截图完成。
func (l *LinkLogic) ReScreenshot(ctx context.Context, linkID xSnowflake.SnowflakeID) *xError.Error {
	link, found, xErr := l.repo.link.GetByID(ctx, linkID, false, nil)
	if xErr != nil {
		return xErr
	}
	if !found {
		return xError.NewError(ctx, xError.NotFound, "友情链接不存在", false)
	}
	if link.Status != constants.LinkStatusApproved {
		return xError.NewError(ctx, xError.BadRequest, "仅已通过的友链支持重新截图", false)
	}

	manager := screenshot.GetManager(ctx)
	if manager == nil {
		return xError.NewError(ctx, xError.ServerInternalError, "截图服务不可用", false)
	}
	manager.Enqueue(linkID)
	return nil
}

// GetPublicLinks 获取公开的友情链接列表
func (l *LinkLogic) GetPublicLinks(ctx context.Context, groupIDStr string) ([]entity.LinkFriend, *xError.Error) {
	var groupID *xSnowflake.SnowflakeID
	if groupIDStr != "" {
		parsedID, err := logcHelper.ParseSnowflakeID(groupIDStr)
		if err == nil {
			groupID = &parsedID
		}
	}

	return l.repo.link.ListPublic(ctx, groupID, constants.LinkStatusApproved, constants.LinkFailNormal, nil)
}

// GetFailedLinks 获取公开「已失效」分组及其下的失效友链。
//
// 与 ListPublic（仅正常友链）独立成档：内置「已失效」分组信息来自 bm_system 配置，
// 读取失败时降级默认分组不阻断；友链限定已通过且失效，避免暴露非公开审核状态。
func (l *LinkLogic) GetFailedLinks(ctx context.Context) (*apiLink.FriendFailedResponse, *xError.Error) {
	invalidGroup, xErr := l.repo.system.BuildBuiltinInvalidGroup(ctx)
	if xErr != nil {
		invalidGroup = entity.NewDefaultBuiltinGroup()
	}

	links, xErr := l.repo.link.ListFailed(ctx, constants.LinkStatusApproved, constants.LinkFailBroken, nil)
	if xErr != nil {
		return nil, xErr
	}

	return &apiLink.FriendFailedResponse{
		Group: invalidGroup,
		Links: links,
	}, nil
}

// Apply 访客自助申请友情链接
//
// 面向游客与登录用户的公开申请入口：组装实体时固定为待审核、一般级别，
// 若申请邮箱已对应注册用户则即时绑定归属，创建后异步通知管理员。
func (l *LinkLogic) Apply(ctx context.Context, req *apiLink.FriendApplyRequest) (*entity.LinkFriend, *xError.Error) {
	link := &entity.LinkFriend{
		Name:        req.LinkName,
		URL:         req.LinkURL,
		Avatar:      xUtil.Ptr(req.LinkAvatar),
		RSS:         xUtil.Ptr(req.LinkRSS),
		Description: xUtil.Ptr(req.LinkDesc),
		Email:       xUtil.Ptr(req.LinkEmail),
		Status:      constants.LinkStatusPending, // 默认待审核
		IsFailure:   constants.LinkFailNormal,    // 默认正常
		Level:       constants.LinkLevelRegular,  // 默认一般级别，由管理员审核时调整
		ApplyRemark: xUtil.Ptr(req.LinkApplyRemark),
	}

	// 申请邮箱若已对应注册用户，则即时绑定归属（否则保持为空，待该用户注册/登录时按邮箱绑定）
	if user, found, xErr := l.repo.user.GetByEmail(ctx, req.LinkEmail); xErr == nil && found {
		link.UserID = &user.ID
	}

	_, xErr := l.repo.link.Create(ctx, link, nil)
	if xErr != nil {
		return nil, xErr
	}

	reloaded, found, xErr := l.repo.link.GetByID(ctx, link.ID, true, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "友情链接不存在", false)
	}

	// 发送邮件通知管理员（xAsync 解耦请求上下文，不阻断主流程）
	xAsync.Async(ctx, func(asyncCtx context.Context) {
		l.sendApplyNotification(asyncCtx, reloaded)
	}, xAsync.WithName("MAIL"))

	return reloaded, nil
}

// ListMine 获取当前用户名下的友情链接列表
func (l *LinkLogic) ListMine(ctx context.Context, userID xSnowflake.SnowflakeID, req *apiLink.FriendUserQueryRequest) (*base.PaginationResponse[entity.LinkFriend], *xError.Error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 10
	}

	links, total, xErr := l.repo.link.List(ctx, &repository.FriendQuery{
		Page:       req.Page,
		PageSize:   req.PageSize,
		LinkStatus: req.LinkStatus,
		UserID:     userID,
		SortBy:     "created_at",
		SortOrder:  "desc",
	}, nil)
	if xErr != nil {
		return nil, xErr
	}

	return base.NewPaginationResponse(links, req.Page, req.PageSize, total), nil
}

// GetMine 获取当前用户名下的友情链接详情（含归属校验）
func (l *LinkLogic) GetMine(ctx context.Context, userID xSnowflake.SnowflakeID, linkID xSnowflake.SnowflakeID) (*entity.LinkFriend, *xError.Error) {
	link, found, xErr := l.repo.link.GetByID(ctx, linkID, true, nil)
	if xErr != nil {
		return nil, xErr
	}
	// 不存在或不归属当前用户时一律返回 NotFound，避免泄露他人友链信息
	if !found || link.UserID == nil || *link.UserID != userID {
		return nil, xError.NewError(ctx, xError.NotFound, "友情链接不存在", false)
	}
	return link, nil
}

// UpdateMine 更新当前用户名下的友情链接（仅站点基础信息，审核状态保持不变）
func (l *LinkLogic) UpdateMine(ctx context.Context, userID xSnowflake.SnowflakeID, linkID xSnowflake.SnowflakeID, req *apiLink.FriendUserUpdateRequest) (*entity.LinkFriend, *xError.Error) {
	link, xErr := l.GetMine(ctx, userID, linkID)
	if xErr != nil {
		return nil, xErr
	}

	if req.LinkName != "" {
		link.Name = req.LinkName
	}
	if req.LinkURL != "" {
		link.URL = req.LinkURL
	}
	if req.LinkAvatar != "" {
		link.Avatar = xUtil.Ptr(req.LinkAvatar)
	}
	if req.LinkRSS != "" {
		link.RSS = xUtil.Ptr(req.LinkRSS)
	}
	if req.LinkDesc != "" {
		link.Description = xUtil.Ptr(req.LinkDesc)
	}
	if req.LinkEmail != "" {
		link.Email = xUtil.Ptr(req.LinkEmail)
	}
	if req.LinkApplyRemark != "" {
		link.ApplyRemark = xUtil.Ptr(req.LinkApplyRemark)
	}

	_, xErr = l.repo.link.Save(ctx, link, nil)
	if xErr != nil {
		return nil, xErr
	}

	reloaded, found, xErr := l.repo.link.GetByID(ctx, linkID, true, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "友情链接不存在", false)
	}
	return reloaded, nil
}

// RequestTakedown 当前用户申请下架自己名下的友情链接
//
// 仅已通过的友链可申请下架，置为下架待审核状态，等待管理员审核。
func (l *LinkLogic) RequestTakedown(ctx context.Context, userID xSnowflake.SnowflakeID, linkID xSnowflake.SnowflakeID) *xError.Error {
	link, xErr := l.GetMine(ctx, userID, linkID)
	if xErr != nil {
		return xErr
	}
	if link.Status != constants.LinkStatusApproved {
		return xError.NewError(ctx, xError.OperationInvalid, "仅已通过的友链可申请下架", false)
	}

	ok, xErr := l.repo.link.UpdateStatusByID(ctx, linkID, constants.LinkStatusTakedownPending, "用户申请下架", nil)
	if xErr != nil {
		return xErr
	}
	if !ok {
		return xError.NewError(ctx, xError.NotFound, "友情链接不存在", false)
	}

	// 异步通知管理员处理下架申请（复用申请通知模板）
	xAsync.Async(ctx, func(asyncCtx context.Context) {
		l.sendApplyNotification(asyncCtx, link)
	}, xAsync.WithName("MAIL"))

	return nil
}

// sendApplyNotification 发送友链申请通知邮件给管理员
//
// 此函数应在 xAsync 异步任务中调用，ctx 为解耦后的独立上下文，不会阻断主流程
func (l *LinkLogic) sendApplyNotification(ctx context.Context, link *entity.LinkFriend) {
	logger := xLog.WithName(xLog.NamedLOGC, "MAIL")

	// 获取配置
	config, xerr := xCtxUtil.Get[*base.BambooConfig](ctx, constants.ContextCustomConfig)
	if xerr != nil {
		logger.Warn(ctx, "无法获取配置，跳过发送申请通知邮件")
		return
	}

	// 检查管理员邮箱是否配置
	if config.Email.AdminEmail == "" {
		logger.Warn(ctx, "管理员邮箱未配置，跳过发送申请通知邮件")
		return
	}

	// 构建模板变量
	linkDesc := ""
	if link.Description != nil {
		linkDesc = *link.Description
	}
	linkEmail := ""
	if link.Email != nil {
		linkEmail = *link.Email
	}

	variables := map[string]string{
		"Username": link.Name,
		"LinkName": link.Name,
		"LinkURL":  link.URL,
		"LinkDesc": linkDesc,
		"Email":    linkEmail,
		"AdminURL": "", // 可后续配置后台管理链接
	}

	// 发送邮件
	mailLogic := NewMailLogic()
	err := mailLogic.SendWithTemplate(
		ctx,
		"apply",
		[]string{config.Email.AdminEmail},
		"【友链申请】收到新的友情链接申请",
		variables,
	)
	if err != nil {
		logger.Warn(ctx, fmt.Sprintf("发送友链申请通知邮件失败: %v", err))
	}
}

// sendStatusNotification 发送审核结果通知邮件给申请者
//
// 此函数应在 xAsync 异步任务中调用，ctx 为解耦后的独立上下文，不会阻断主流程
func (l *LinkLogic) sendStatusNotification(ctx context.Context, link *entity.LinkFriend, status int, reviewRemark string) {
	logger := xLog.WithName(xLog.NamedLOGC, "MAIL")

	// 检查友链是否有邮箱
	if link.Email == nil || *link.Email == "" {
		logger.Info(ctx, fmt.Sprintf("友链 %s 无联系邮箱，跳过发送审核通知", link.Name))
		return
	}

	// 根据状态选择模板和主题
	var templateName, subject string
	switch status {
	case constants.LinkStatusApproved:
		templateName = "approved"
		subject = "🎉 您的友链申请已通过"
	case constants.LinkStatusRejected:
		templateName = "rejected"
		subject = "📋 您的友链申请审核结果"
	case constants.LinkStatusTakenDown:
		templateName = "takedown"
		subject = "📤 您的友链已下架"
	default:
		// 非通过/拒绝/下架状态不发送邮件
		return
	}

	// 构建模板变量
	variables := map[string]string{
		"Username":     link.Name,
		"LinkName":     link.Name,
		"LinkURL":      link.URL,
		"RejectReason": reviewRemark,
	}

	// 发送邮件
	mailLogic := NewMailLogic()
	err := mailLogic.SendWithTemplate(
		ctx,
		templateName,
		[]string{*link.Email},
		subject,
		variables,
	)
	if err != nil {
		logger.Warn(ctx, fmt.Sprintf("发送友链审核通知邮件失败: %v", err))
	}
}
