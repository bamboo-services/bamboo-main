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

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/major/utility/context"
	apiSponsor "github.com/bamboo-services/bamboo-main/api/sponsor"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/internal/models/base"
	"github.com/bamboo-services/bamboo-main/internal/repository"
)

type sponsorRecordRepo struct {
	record  *repository.SponsorRecordRepo
	channel *repository.SponsorChannelRepo
}

// SponsorRecordLogic 赞助记录业务逻辑
type SponsorRecordLogic struct {
	logic
	repo sponsorRecordRepo
}

// NewSponsorRecordLogic 创建 SponsorRecordLogic 实例，从上下文获取数据库与缓存并初始化赞助记录与渠道仓储依赖。
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
	})
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

	record, found, xErr = l.repo.record.UpdateByID(ctx, recordID, updates)
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

	deleted, xErr := l.repo.record.HardDeleteByID(ctx, recordID)
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

// GetPage 分页获取赞助记录。
func (l *SponsorRecordLogic) GetPage(ctx context.Context, req *apiSponsor.RecordPageRequest) (*base.PaginationResponse[apiSponsor.RecordEntityResponse], *xError.Error) {
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
	return base.NewPaginationResponse(resp, req.Page, req.PageSize, total), nil
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

func buildRecordEntityResponse(record *entity.SponsorRecord) *apiSponsor.RecordEntityResponse {
	if record == nil {
		return nil
	}
	return &apiSponsor.RecordEntityResponse{
		ID:          record.ID,
		Nickname:    record.Nickname,
		RedirectURL: record.RedirectURL,
		Amount:      record.Amount,
		ChannelID:   record.ChannelID,
		Message:     record.Message,
		SponsorAt:   record.SponsorAt,
		SortOrder:   record.SortOrder,
		IsAnonymous: record.IsAnonymous,
		IsHidden:    record.IsHidden,
		CreatedAt:   record.CreatedAt,
		UpdatedAt:   record.UpdatedAt,
		Channel:     buildRecordChannelResponse(record.ChannelFKey),
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
