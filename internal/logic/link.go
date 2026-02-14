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

package logic

import (
	"context"
	"fmt"
	"strconv"

	apiLink "github.com/bamboo-services/bamboo-main/api/link"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	logcHelper "github.com/bamboo-services/bamboo-main/internal/logic/helper"
	"github.com/bamboo-services/bamboo-main/internal/models/base"
	"github.com/bamboo-services/bamboo-main/internal/repository"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
	ctxUtil "github.com/bamboo-services/bamboo-main/pkg/util/ctx"

	xError "github.com/bamboo-services/bamboo-base-go/error"
	xLog "github.com/bamboo-services/bamboo-base-go/log"
	xUtil "github.com/bamboo-services/bamboo-base-go/utility"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/utility/ctxutil"
	"github.com/gin-gonic/gin"
)

type linkRepo struct {
	link *repository.LinkRepo
}

// LinkLogic 友情链接业务逻辑
type LinkLogic struct {
	logic
	repo linkRepo
}

func NewLinkLogic(ctx context.Context) *LinkLogic {
	db := xCtxUtil.MustGetDB(ctx)
	rdb := xCtxUtil.MustGetRDB(ctx)

	return &LinkLogic{
		logic: logic{
			db:  db,
			rdb: rdb,
			log: xLog.WithName(xLog.NamedLOGC, "LinkLogic"),
		},
		repo: linkRepo{link: repository.NewLinkRepo(db, rdb)},
	}
}

// Add 添加友情链接
func (l *LinkLogic) Add(ctx *gin.Context, req *apiLink.FriendAddRequest) (*entity.LinkFriend, *xError.Error) {
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

	// 发送邮件通知管理员（异步，不阻断主流程）
	go l.sendApplyNotification(ctx, reloaded)

	return reloaded, nil
}

// Update 更新友情链接
func (l *LinkLogic) Update(ctx *gin.Context, linkIDStr string, req *apiLink.FriendUpdateRequest) (*entity.LinkFriend, *xError.Error) {
	// 解析ID
	linkID, err := strconv.ParseInt(linkIDStr, 10, 64)
	if err != nil {
		return nil, xError.NewError(ctx, xError.BadRequest, "无效的友链ID", false)
	}

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
func (l *LinkLogic) Delete(ctx *gin.Context, linkIDStr string) *xError.Error {
	// 解析ID
	linkID, err := strconv.ParseInt(linkIDStr, 10, 64)
	if err != nil {
		return xError.NewError(ctx, xError.BadRequest, "无效的友链ID", false)
	}

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
func (l *LinkLogic) Get(ctx *gin.Context, linkIDStr string) (*entity.LinkFriend, *xError.Error) {
	// 解析ID
	linkID, err := strconv.ParseInt(linkIDStr, 10, 64)
	if err != nil {
		return nil, xError.NewError(ctx, xError.BadRequest, "无效的友链ID", false)
	}

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
func (l *LinkLogic) List(ctx *gin.Context, req *apiLink.FriendQueryRequest) (*base.PaginationResponse[entity.LinkFriend], *xError.Error) {
	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 10
	}

	links, total, xErr := l.repo.link.List(ctx, req, nil)
	if xErr != nil {
		return nil, xErr
	}

	return base.NewPaginationResponse(links, req.Page, req.PageSize, total), nil
}

// UpdateStatus 更新友情链接状态
func (l *LinkLogic) UpdateStatus(ctx *gin.Context, linkIDStr string, req *apiLink.FriendStatusRequest) *xError.Error {
	// 解析ID
	linkID, err := strconv.ParseInt(linkIDStr, 10, 64)
	if err != nil {
		return xError.NewError(ctx, xError.BadRequest, "无效的友链ID", false)
	}

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

	// 发送审核结果邮件通知（异步，不阻断主流程）
	go l.sendStatusNotification(ctx, link, req.LinkStatus, req.LinkReviewRemark)

	return nil
}

// UpdateFailStatus 更新友情链接失效状态
func (l *LinkLogic) UpdateFailStatus(ctx *gin.Context, linkIDStr string, req *apiLink.FriendFailRequest) *xError.Error {
	// 解析ID
	linkID, err := strconv.ParseInt(linkIDStr, 10, 64)
	if err != nil {
		return xError.NewError(ctx, xError.BadRequest, "无效的友链ID", false)
	}

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
func (l *LinkLogic) GetPublicLinks(ctx *gin.Context, groupIDStr string) ([]entity.LinkFriend, *xError.Error) {
	var groupID *int64
	if groupIDStr != "" {
		parsedID, err := strconv.ParseInt(groupIDStr, 10, 64)
		if err == nil {
			groupID = &parsedID
		}
	}

	return l.repo.link.ListPublic(ctx, groupID, constants.LinkStatusApproved, constants.LinkFailNormal, nil)
}

// sendApplyNotification 发送友链申请通知邮件给管理员
//
// 此函数应在 goroutine 中异步调用，不会阻断主流程
func (l *LinkLogic) sendApplyNotification(ctx *gin.Context, link *entity.LinkFriend) {
	logger := xLog.WithName(xLog.NamedLOGC, "MAIL")

	// 获取配置
	config := ctxUtil.GetConfig(ctx)
	if config == nil {
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
		"FromName": config.Email.FromName,
	}

	// 发送邮件
	mailLogic := &MailLogic{TemplateService: &logcHelper.MailTemplateLogic{}, MaxRetry: 3}
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
// 此函数应在 goroutine 中异步调用，不会阻断主流程
func (l *LinkLogic) sendStatusNotification(ctx *gin.Context, link *entity.LinkFriend, status int, reviewRemark string) {
	logger := xLog.WithName(xLog.NamedLOGC, "MAIL")

	// 检查友链是否有邮箱
	if link.Email == nil || *link.Email == "" {
		logger.Info(ctx, fmt.Sprintf("友链 %s 无联系邮箱，跳过发送审核通知", link.Name))
		return
	}

	// 获取配置
	config := ctxUtil.GetConfig(ctx)
	if config == nil {
		logger.Warn(ctx, "无法获取配置，跳过发送审核通知邮件")
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
		"FromName":     config.Email.FromName,
	}

	// 发送邮件
	mailLogic := &MailLogic{TemplateService: &logcHelper.MailTemplateLogic{}, MaxRetry: 3}
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
