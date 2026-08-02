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
)

type linkRepo struct {
	link *repository.LinkRepo
	user *repository.SystemUserRepo
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
			link: repository.NewLinkRepo(db, m),
			user: repository.NewSystemUserRepo(db, m),
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
	if req.LinkGroupID != 0 {
		link.GroupID = &req.LinkGroupID
	}
	if req.LinkColorID != 0 {
		link.ColorID = &req.LinkColorID
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
	if req.LinkGroupID != 0 {
		link.GroupID = &req.LinkGroupID
	}
	if req.LinkColorID != 0 {
		link.ColorID = &req.LinkColorID
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
func (l *LinkLogic) List(ctx context.Context, req *apiLink.FriendQueryRequest) (*base.PaginationResponse[entity.LinkFriend], *xError.Error) {
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
		LinkGroupID: req.LinkGroupID,
		SortBy:      req.SortBy,
		SortOrder:   req.SortOrder,
	}, nil)
	if xErr != nil {
		return nil, xErr
	}

	return base.NewPaginationResponse(links, req.Page, req.PageSize, total), nil
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

// UpdateFailStatus 更新友情链接失效状态
func (l *LinkLogic) UpdateFailStatus(ctx context.Context, linkID xSnowflake.SnowflakeID, req *apiLink.FriendFailRequest) *xError.Error {
	ok, xErr := l.repo.link.UpdateFailureByID(ctx, linkID, req.LinkFail, req.LinkFailReason, nil)
	if xErr != nil {
		return xErr
	}
	if !ok {
		return xError.NewError(ctx, xError.NotFound, "友情链接不存在", false)
	}

	return nil
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
