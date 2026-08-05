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

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
	xUtil "github.com/bamboo-services/bamboo-base-go/common/utility"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/major/utility/context"
	xAsync "github.com/bamboo-services/bamboo-base-go/plugins/async"
	apiSponsor "github.com/bamboo-services/bamboo-main/api/sponsor"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/internal/models/base"
	"github.com/bamboo-services/bamboo-main/internal/repository"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
)

type sponsorRecordRepo struct {
	record  *repository.SponsorRecordRepo
	channel *repository.SponsorChannelRepo
	user    *repository.SystemUserRepo
}

// SponsorRecordLogic 赞助记录业务逻辑
type SponsorRecordLogic struct {
	logic
	repo sponsorRecordRepo
}

// NewSponsorRecordLogic 创建 SponsorRecordLogic 实例，从上下文获取数据库与缓存并初始化赞助记录、渠道与用户仓储依赖。
func NewSponsorRecordLogic(ctx context.Context) *SponsorRecordLogic {
	db := xCtxUtil.MustGetDB(ctx)
	m := xCtxUtil.MustGetCacheManager(ctx)

	return &SponsorRecordLogic{
		logic: logic{
			db:    db,
			cache: m,
			log:   xLog.WithName(xLog.NamedLOGC, "SponsorRecordLogic"),
		},
		repo: sponsorRecordRepo{
			record:  repository.NewSponsorRecordRepo(db, m),
			channel: repository.NewSponsorChannelRepo(db, m),
			user:    repository.NewSystemUserRepo(db, m),
		},
	}
}

// Add 添加赞助记录，校验指定渠道存在后创建记录。
func (l *SponsorRecordLogic) Add(ctx context.Context, req *apiSponsor.RecordAddRequest) (*apiSponsor.RecordEntityResponse, *xError.Error) {
	if req.ChannelID != nil {
		_, found, xErr := l.repo.channel.GetByID(ctx, *req.ChannelID)
		if xErr != nil {
			return nil, xErr
		}
		if !found {
			return nil, xError.NewError(ctx, xError.NotFound, "赞助渠道不存在", false)
		}
	}

	record, xErr := l.repo.record.Create(ctx, &entity.SponsorRecord{
		Nickname:    req.Nickname,
		RedirectURL: req.RedirectURL,
		Amount:      req.Amount,
		ChannelID:   req.ChannelID,
		Message:     req.Message,
		SponsorAt:   req.SponsorAt,
		SortOrder:   req.SortOrder,
		IsAnonymous: req.IsAnonymous,
		IsHidden:    req.IsHidden,
		Status:      constants.SponsorStatusApproved, // 管理员手动录入即视为已通过，直接前台展示
	}, nil)
	if xErr != nil {
		return nil, xErr
	}
	return buildRecordEntityResponse(record), nil
}

// Update 更新赞助记录，校验记录与渠道存在后按字段增量更新。
func (l *SponsorRecordLogic) Update(ctx context.Context, recordID xSnowflake.SnowflakeID, req *apiSponsor.RecordUpdateRequest) (*apiSponsor.RecordEntityResponse, *xError.Error) {
	_, found, xErr := l.repo.record.GetByID(ctx, recordID)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "赞助记录不存在", false)
	}

	if req.ChannelID != nil {
		_, found, xErr = l.repo.channel.GetByID(ctx, *req.ChannelID)
		if xErr != nil {
			return nil, xErr
		}
		if !found {
			return nil, xError.NewError(ctx, xError.NotFound, "赞助渠道不存在", false)
		}
	}

	updates := make(map[string]any)
	if req.Nickname != nil {
		updates["nickname"] = *req.Nickname
	}
	if req.RedirectURL != nil {
		updates["redirect_url"] = req.RedirectURL
	}
	if req.Amount != nil {
		updates["amount"] = *req.Amount
	}
	if req.ChannelID != nil {
		updates["channel_id"] = req.ChannelID
	}
	if req.Message != nil {
		updates["message"] = req.Message
	}
	if req.SponsorAt != nil {
		updates["sponsor_at"] = req.SponsorAt
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.IsAnonymous != nil {
		updates["is_anonymous"] = *req.IsAnonymous
	}
	if req.IsHidden != nil {
		updates["is_hidden"] = *req.IsHidden
	}

	var record *entity.SponsorRecord
	if len(updates) == 0 {
		record, found, xErr = l.repo.record.GetDetailByID(ctx, recordID)
		if xErr != nil {
			return nil, xErr
		}
		if !found {
			return nil, xError.NewError(ctx, xError.NotFound, "赞助记录不存在", false)
		}
		return buildRecordEntityResponse(record), nil
	}

	record, found, xErr = l.repo.record.UpdateByID(ctx, recordID, updates, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "赞助记录不存在", false)
	}
	return buildRecordEntityResponse(record), nil
}

// Delete 删除赞助记录，校验存在后硬删除。
func (l *SponsorRecordLogic) Delete(ctx context.Context, recordID xSnowflake.SnowflakeID) *xError.Error {
	_, found, xErr := l.repo.record.GetByID(ctx, recordID)
	if xErr != nil {
		return xErr
	}
	if !found {
		return xError.NewError(ctx, xError.NotFound, "赞助记录不存在", false)
	}

	deleted, xErr := l.repo.record.HardDeleteByID(ctx, recordID, nil)
	if xErr != nil {
		return xErr
	}
	if !deleted {
		return xError.NewError(ctx, xError.NotFound, "赞助记录不存在", false)
	}
	return nil
}

// Get 获取赞助记录详情（含渠道信息）。
func (l *SponsorRecordLogic) Get(ctx context.Context, recordID xSnowflake.SnowflakeID) (*apiSponsor.RecordEntityResponse, *xError.Error) {
	record, found, xErr := l.repo.record.GetDetailByID(ctx, recordID)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "赞助记录不存在", false)
	}
	return buildRecordEntityResponse(record), nil
}

// GetPage 分页获取赞助记录（管理端，附带待审核数量供入口徽章展示）。
func (l *SponsorRecordLogic) GetPage(ctx context.Context, req *apiSponsor.RecordPageRequest) (*apiSponsor.RecordAdminPageResponse, *xError.Error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 10
	}

	query := repository.SponsorRecordPageQuery{
		Page:        req.Page,
		PageSize:    req.PageSize,
		ChannelID:   req.ChannelID,
		Nickname:    stringPointerValue(req.Nickname),
		IsAnonymous: req.IsAnonymous,
		IsHidden:    req.IsHidden,
		Status:      req.Status,
		OrderBy:     stringPointerValue(req.OrderBy),
		Order:       stringPointerValue(req.Order),
	}

	records, total, xErr := l.repo.record.Page(ctx, query)
	if xErr != nil {
		return nil, xErr
	}

	resp := make([]apiSponsor.RecordEntityResponse, 0, len(records))
	for _, item := range records {
		row := buildRecordEntityResponse(&item)
		if row != nil {
			resp = append(resp, *row)
		}
	}

	// 附带待审核数量供管理端入口徽章展示（统计失败不阻断列表返回）
	pendingCount, _ := l.repo.record.CountByStatus(ctx, constants.SponsorStatusPending, nil)

	result := base.NewPaginationResponse(resp, req.Page, req.PageSize, total)
	return &apiSponsor.RecordAdminPageResponse{
		PaginationResponse: *result,
		PendingCount:       pendingCount,
	}, nil
}

// GetPublicPage 分页获取公开的赞助记录，匿名记录隐藏昵称与跳转链接。
func (l *SponsorRecordLogic) GetPublicPage(ctx context.Context, req *apiSponsor.RecordPublicPageRequest) (*base.PaginationResponse[apiSponsor.RecordPublicItemResponse], *xError.Error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 50 {
		req.PageSize = 20
	}

	query := repository.SponsorRecordPublicPageQuery{
		Page:      req.Page,
		PageSize:  req.PageSize,
		ChannelID: req.ChannelID,
		OrderBy:   stringPointerValue(req.OrderBy),
		Order:     stringPointerValue(req.Order),
	}

	records, total, xErr := l.repo.record.PublicPage(ctx, query)
	if xErr != nil {
		return nil, xErr
	}

	resp := make([]apiSponsor.RecordPublicItemResponse, 0, len(records))
	for _, item := range records {
		resp = append(resp, buildRecordPublicItemResponse(&item))
	}
	return base.NewPaginationResponse(resp, req.Page, req.PageSize, total), nil
}

// Apply 访客自助申请赞助展示
//
// 面向游客与登录用户的公开申请入口：组装实体时固定为待审核状态，
// 若申请邮箱已对应注册用户则即时绑定归属，创建后异步通知管理员。
func (l *SponsorRecordLogic) Apply(ctx context.Context, req *apiSponsor.SponsorApplyRequest) (*apiSponsor.RecordEntityResponse, *xError.Error) {
	// 校验渠道存在
	if req.ChannelID != nil {
		_, found, xErr := l.repo.channel.GetByID(ctx, *req.ChannelID)
		if xErr != nil {
			return nil, xErr
		}
		if !found {
			return nil, xError.NewError(ctx, xError.NotFound, "赞助渠道不存在", false)
		}
	}

	record := &entity.SponsorRecord{
		Nickname:    req.Nickname,
		RedirectURL: req.RedirectURL,
		Amount:      req.Amount,
		ChannelID:   req.ChannelID,
		Message:     req.Message,
		SponsorAt:   req.SponsorAt,
		Email:       xUtil.Ptr(req.Email),
		IsAnonymous: req.IsAnonymous,
		ApplyRemark: req.ApplyRemark,
		Status:      constants.SponsorStatusPending, // 默认待审核，经管理员审核后前台展示
	}

	// 申请邮箱若已对应注册用户，则即时绑定归属（否则保持为空，待该用户注册/登录时按邮箱绑定）
	if user, found, xErr := l.repo.user.GetByEmail(ctx, req.Email); xErr == nil && found {
		record.UserID = &user.ID
	}

	_, xErr := l.repo.record.Create(ctx, record, nil)
	if xErr != nil {
		return nil, xErr
	}

	reloaded, found, xErr := l.repo.record.GetDetailByID(ctx, record.ID)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "赞助记录不存在", false)
	}

	// 发送邮件通知管理员（xAsync 解耦请求上下文，不阻断主流程）
	xAsync.Async(ctx, func(asyncCtx context.Context) {
		l.sendApplyNotification(asyncCtx, reloaded)
	}, xAsync.WithName("MAIL"))

	return buildRecordEntityResponse(reloaded), nil
}

// UpdateStatus 更新赞助记录审核状态
func (l *SponsorRecordLogic) UpdateStatus(ctx context.Context, recordID xSnowflake.SnowflakeID, req *apiSponsor.SponsorStatusRequest) *xError.Error {
	// 先查询赞助记录信息（用于发送邮件通知）
	record, found, xErr := l.repo.record.GetDetailByID(ctx, recordID)
	if xErr != nil {
		return xErr
	}
	if !found {
		return xError.NewError(ctx, xError.NotFound, "赞助记录不存在", false)
	}

	ok, xErr := l.repo.record.UpdateStatusByID(ctx, recordID, req.SponsorStatus, req.SponsorReviewRemark, nil)
	if xErr != nil {
		return xErr
	}
	if !ok {
		return xError.NewError(ctx, xError.NotFound, "赞助记录不存在", false)
	}

	// 发送审核结果邮件通知（xAsync 解耦请求上下文，不阻断主流程）
	xAsync.Async(ctx, func(asyncCtx context.Context) {
		l.sendStatusNotification(asyncCtx, record, req.SponsorStatus, req.SponsorReviewRemark)
	}, xAsync.WithName("MAIL"))

	return nil
}

// ListMine 获取当前用户名下的赞助记录列表
func (l *SponsorRecordLogic) ListMine(ctx context.Context, userID xSnowflake.SnowflakeID, req *apiSponsor.SponsorUserQueryRequest) (*base.PaginationResponse[apiSponsor.RecordEntityResponse], *xError.Error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 10
	}

	records, total, xErr := l.repo.record.Page(ctx, repository.SponsorRecordPageQuery{
		Page:     req.Page,
		PageSize: req.PageSize,
		Status:   req.SponsorStatus,
		UserID:   userID,
		OrderBy:  "created_at",
		Order:    "desc",
	})
	if xErr != nil {
		return nil, xErr
	}

	resp := make([]apiSponsor.RecordEntityResponse, 0, len(records))
	for _, item := range records {
		row := buildRecordEntityResponse(&item)
		if row != nil {
			resp = append(resp, *row)
		}
	}
	return base.NewPaginationResponse(resp, req.Page, req.PageSize, total), nil
}

// GetMine 获取当前用户名下的赞助记录详情（含归属校验）
func (l *SponsorRecordLogic) GetMine(ctx context.Context, userID xSnowflake.SnowflakeID, recordID xSnowflake.SnowflakeID) (*apiSponsor.RecordEntityResponse, *xError.Error) {
	record, found, xErr := l.repo.record.GetDetailByID(ctx, recordID)
	if xErr != nil {
		return nil, xErr
	}
	// 不存在或不归属当前用户时一律返回 NotFound，避免泄露他人赞助信息
	if !found || record.UserID == nil || *record.UserID != userID {
		return nil, xError.NewError(ctx, xError.NotFound, "赞助记录不存在", false)
	}
	return buildRecordEntityResponse(record), nil
}

// UpdateMine 更新当前用户名下的赞助记录（仅展示类基础字段，审核状态与金额/渠道保持不变）
func (l *SponsorRecordLogic) UpdateMine(ctx context.Context, userID xSnowflake.SnowflakeID, recordID xSnowflake.SnowflakeID, req *apiSponsor.SponsorUserUpdateRequest) (*apiSponsor.RecordEntityResponse, *xError.Error) {
	// 先做归属校验（非本人一律 NotFound）
	origin, xErr := l.GetMine(ctx, userID, recordID)
	if xErr != nil {
		return nil, xErr
	}

	updates := make(map[string]any)
	if req.Nickname != nil {
		updates["nickname"] = *req.Nickname
	}
	if req.RedirectURL != nil {
		updates["redirect_url"] = req.RedirectURL
	}
	if req.Message != nil {
		updates["message"] = req.Message
	}
	if req.SponsorAt != nil {
		updates["sponsor_at"] = req.SponsorAt
	}
	if req.IsAnonymous != nil {
		updates["is_anonymous"] = *req.IsAnonymous
	}
	if req.ApplyRemark != nil {
		updates["apply_remark"] = req.ApplyRemark
	}

	if len(updates) == 0 {
		return origin, nil
	}

	updated, found, xErr := l.repo.record.UpdateByID(ctx, recordID, updates, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "赞助记录不存在", false)
	}
	return buildRecordEntityResponse(updated), nil
}

// sendApplyNotification 发送赞助申请通知邮件给管理员
//
// 此函数应在 xAsync 异步任务中调用，ctx 为解耦后的独立上下文，不会阻断主流程
func (l *SponsorRecordLogic) sendApplyNotification(ctx context.Context, record *entity.SponsorRecord) {
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
	email := ""
	if record.Email != nil {
		email = *record.Email
	}
	channelName := ""
	if record.ChannelFKey != nil {
		channelName = record.ChannelFKey.Name
	}

	variables := map[string]string{
		"Username": record.Nickname,
		"Amount":   fmt.Sprintf("%.2f", float64(record.Amount)/100),
		"Channel":  channelName,
		"Message":  stringPointerValue(record.Message),
		"Email":    email,
		"AdminURL": "", // 可后续配置后台管理链接
	}

	// 发送邮件
	mailLogic := NewMailLogic()
	err := mailLogic.SendWithTemplate(
		ctx,
		constants.SponsorEmailTypeApply,
		[]string{config.Email.AdminEmail},
		"【赞助申请】收到新的赞助展示申请",
		variables,
	)
	if err != nil {
		logger.Warn(ctx, fmt.Sprintf("发送赞助申请通知邮件失败: %v", err))
	}
}

// sendStatusNotification 发送审核结果通知邮件给申请者
//
// 此函数应在 xAsync 异步任务中调用，ctx 为解耦后的独立上下文，不会阻断主流程
func (l *SponsorRecordLogic) sendStatusNotification(ctx context.Context, record *entity.SponsorRecord, status int, reviewRemark string) {
	logger := xLog.WithName(xLog.NamedLOGC, "MAIL")

	// 检查赞助记录是否有邮箱
	if record.Email == nil || *record.Email == "" {
		logger.Info(ctx, fmt.Sprintf("赞助记录 %s 无联系邮箱，跳过发送审核通知", record.Nickname))
		return
	}

	// 根据状态选择模板和主题
	var templateName, subject string
	switch status {
	case constants.SponsorStatusApproved:
		templateName = constants.SponsorEmailTypeApproved
		subject = "🎉 您的赞助展示申请已通过"
	case constants.SponsorStatusRejected:
		templateName = constants.SponsorEmailTypeRejected
		subject = "📋 您的赞助展示申请审核结果"
	default:
		// 非通过/拒绝状态不发送邮件
		return
	}

	// 构建模板变量
	variables := map[string]string{
		"Username":     record.Nickname,
		"Amount":       fmt.Sprintf("%.2f", float64(record.Amount)/100),
		"RejectReason": reviewRemark,
	}

	// 发送邮件
	mailLogic := NewMailLogic()
	err := mailLogic.SendWithTemplate(
		ctx,
		templateName,
		[]string{*record.Email},
		subject,
		variables,
	)
	if err != nil {
		logger.Warn(ctx, fmt.Sprintf("发送赞助审核通知邮件失败: %v", err))
	}
}

func buildRecordEntityResponse(record *entity.SponsorRecord) *apiSponsor.RecordEntityResponse {
	if record == nil {
		return nil
	}
	return &apiSponsor.RecordEntityResponse{
		ID:           record.ID,
		Nickname:     record.Nickname,
		RedirectURL:  record.RedirectURL,
		Amount:       record.Amount,
		ChannelID:    record.ChannelID,
		Message:      record.Message,
		SponsorAt:    record.SponsorAt,
		SortOrder:    record.SortOrder,
		IsAnonymous:  record.IsAnonymous,
		IsHidden:     record.IsHidden,
		Status:       record.Status,
		Email:        record.Email,
		UserID:       record.UserID,
		ApplyRemark:  record.ApplyRemark,
		ReviewRemark: record.ReviewRemark,
		CreatedAt:    record.CreatedAt,
		UpdatedAt:    record.UpdatedAt,
		Channel:      buildRecordChannelResponse(record.ChannelFKey),
	}
}

func buildRecordPublicItemResponse(record *entity.SponsorRecord) apiSponsor.RecordPublicItemResponse {
	nickname := record.Nickname
	redirectURL := record.RedirectURL
	if record.IsAnonymous {
		nickname = "匿名用户"
		redirectURL = nil
	}

	return apiSponsor.RecordPublicItemResponse{
		ID:          record.ID,
		Nickname:    nickname,
		RedirectURL: redirectURL,
		Amount:      record.Amount,
		Message:     record.Message,
		SponsorAt:   record.SponsorAt,
		Channel:     buildRecordChannelResponse(record.ChannelFKey),
	}
}

func buildRecordChannelResponse(channel *entity.SponsorChannel) *apiSponsor.SponsorChannelSimpleResponse {
	if channel == nil {
		return nil
	}
	return &apiSponsor.SponsorChannelSimpleResponse{
		ID:   channel.ID,
		Name: channel.Name,
		Icon: channel.Icon,
	}
}
