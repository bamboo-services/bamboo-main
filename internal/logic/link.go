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
		repo: linkRepo{link: repository.NewLinkRepo(db, m)},
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
	default:
		// 非通过/拒绝状态不发送邮件
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
