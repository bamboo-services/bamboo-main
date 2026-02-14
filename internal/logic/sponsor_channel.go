package logic

import (
	"context"
	"strconv"

	xError "github.com/bamboo-services/bamboo-base-go/error"
	xLog "github.com/bamboo-services/bamboo-base-go/log"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/utility/ctxutil"
	apiSponsor "github.com/bamboo-services/bamboo-main/api/sponsor"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/internal/models/base"
	"github.com/bamboo-services/bamboo-main/internal/repository"
	"github.com/gin-gonic/gin"
)

type sponsorChannelRepo struct {
	channel *repository.SponsorChannelRepo
}

type SponsorChannelLogic struct {
	logic
	repo sponsorChannelRepo
}

func NewSponsorChannelLogic(ctx context.Context) *SponsorChannelLogic {
	db := xCtxUtil.MustGetDB(ctx)
	rdb := xCtxUtil.MustGetRDB(ctx)

	return &SponsorChannelLogic{
		logic: logic{
			db:  db,
			rdb: rdb,
			log: xLog.WithName(xLog.NamedLOGC, "SponsorChannelLogic"),
		},
		repo: sponsorChannelRepo{
			channel: repository.NewSponsorChannelRepo(db, rdb),
		},
	}
}

func (l *SponsorChannelLogic) Add(ctx *gin.Context, req *apiSponsor.ChannelAddRequest) (*apiSponsor.ChannelEntityResponse, *xError.Error) {
	channel, xErr := l.repo.channel.Create(ctx, &entity.SponsorChannel{
		Name:        req.Name,
		Icon:        req.Icon,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		Status:      true,
	})
	if xErr != nil {
		return nil, xErr
	}

	sponsorCount, xErr := l.repo.channel.CountSponsorsByChannelID(ctx, channel.ID)
	if xErr != nil {
		return nil, xErr
	}
	return buildChannelEntityResponse(channel, int(sponsorCount)), nil
}

func (l *SponsorChannelLogic) Update(ctx *gin.Context, idStr string, req *apiSponsor.ChannelUpdateRequest) (*apiSponsor.ChannelEntityResponse, *xError.Error) {
	channelID, xErr := parseSponsorChannelID(ctx, idStr)
	if xErr != nil {
		return nil, xErr
	}

	channel, found, xErr := l.repo.channel.GetByID(ctx, channelID)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "赞助渠道不存在", false)
	}

	updates := make(map[string]any)
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Icon != nil {
		updates["icon"] = req.Icon
	}
	if req.Description != nil {
		updates["description"] = req.Description
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}

	if len(updates) > 0 {
		channel, found, xErr = l.repo.channel.UpdateByID(ctx, channelID, updates)
		if xErr != nil {
			return nil, xErr
		}
		if !found {
			return nil, xError.NewError(ctx, xError.NotFound, "赞助渠道不存在", false)
		}
	}

	sponsorCount, xErr := l.repo.channel.CountSponsorsByChannelID(ctx, channelID)
	if xErr != nil {
		return nil, xErr
	}
	return buildChannelEntityResponse(channel, int(sponsorCount)), nil
}

func (l *SponsorChannelLogic) UpdateStatus(ctx *gin.Context, idStr string, req *apiSponsor.ChannelStatusRequest) (bool, *xError.Error) {
	channelID, xErr := parseSponsorChannelID(ctx, idStr)
	if xErr != nil {
		return false, xErr
	}

	_, found, xErr := l.repo.channel.UpdateByID(ctx, channelID, map[string]any{"status": req.Status})
	if xErr != nil {
		return false, xErr
	}
	if !found {
		return false, xError.NewError(ctx, xError.NotFound, "赞助渠道不存在", false)
	}
	return req.Status, nil
}

func (l *SponsorChannelLogic) Delete(ctx *gin.Context, idStr string) *xError.Error {
	channelID, xErr := parseSponsorChannelID(ctx, idStr)
	if xErr != nil {
		return xErr
	}

	_, found, xErr := l.repo.channel.GetByID(ctx, channelID)
	if xErr != nil {
		return xErr
	}
	if !found {
		return xError.NewError(ctx, xError.NotFound, "赞助渠道不存在", false)
	}

	sponsorCount, xErr := l.repo.channel.CountSponsorsByChannelID(ctx, channelID)
	if xErr != nil {
		return xErr
	}
	if sponsorCount > 0 {
		return xError.NewError(ctx, xError.BadRequest, "该渠道下存在赞助记录，无法删除", false)
	}

	deleted, xErr := l.repo.channel.HardDeleteByID(ctx, channelID)
	if xErr != nil {
		return xErr
	}
	if !deleted {
		return xError.NewError(ctx, xError.NotFound, "赞助渠道不存在", false)
	}
	return nil
}

func (l *SponsorChannelLogic) Get(ctx *gin.Context, idStr string) (*apiSponsor.ChannelEntityResponse, *xError.Error) {
	channelID, xErr := parseSponsorChannelID(ctx, idStr)
	if xErr != nil {
		return nil, xErr
	}

	channel, found, xErr := l.repo.channel.GetByID(ctx, channelID)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "赞助渠道不存在", false)
	}

	sponsorCount, xErr := l.repo.channel.CountSponsorsByChannelID(ctx, channelID)
	if xErr != nil {
		return nil, xErr
	}
	return buildChannelEntityResponse(channel, int(sponsorCount)), nil
}

func (l *SponsorChannelLogic) GetList(ctx *gin.Context, req *apiSponsor.ChannelListRequest) ([]apiSponsor.ChannelListItemResponse, *xError.Error) {
	listQuery := repository.SponsorChannelListQuery{
		OnlyEnabled: req.OnlyEnabled != nil && *req.OnlyEnabled,
		OrderBy:     stringPointerValue(req.OrderBy),
		Order:       stringPointerValue(req.Order),
	}
	if req.Status != nil {
		listQuery.Status = req.Status
	}
	if req.Name != nil {
		listQuery.Name = *req.Name
	}

	channels, xErr := l.repo.channel.List(ctx, listQuery)
	if xErr != nil {
		return nil, xErr
	}

	sponsorCounts, xErr := l.repo.channel.CountSponsorsByChannelIDs(ctx, collectChannelIDs(channels))
	if xErr != nil {
		return nil, xErr
	}

	resp := make([]apiSponsor.ChannelListItemResponse, 0, len(channels))
	for _, channel := range channels {
		resp = append(resp, buildChannelListItemResponse(&channel, int(sponsorCounts[channel.ID])))
	}
	return resp, nil
}

func (l *SponsorChannelLogic) GetPage(ctx *gin.Context, req *apiSponsor.ChannelPageRequest) (*base.PaginationResponse[apiSponsor.ChannelEntityResponse], *xError.Error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 10
	}

	pageQuery := repository.SponsorChannelPageQuery{
		Page:     req.Page,
		PageSize: req.PageSize,
		OrderBy:  stringPointerValue(req.OrderBy),
		Order:    stringPointerValue(req.Order),
	}
	if req.Status != nil {
		pageQuery.Status = req.Status
	}
	if req.Name != nil {
		pageQuery.Name = *req.Name
	}

	channels, total, xErr := l.repo.channel.Page(ctx, pageQuery)
	if xErr != nil {
		return nil, xErr
	}

	sponsorCounts, xErr := l.repo.channel.CountSponsorsByChannelIDs(ctx, collectChannelIDs(channels))
	if xErr != nil {
		return nil, xErr
	}

	resp := make([]apiSponsor.ChannelEntityResponse, 0, len(channels))
	for _, channel := range channels {
		item := buildChannelEntityResponse(&channel, int(sponsorCounts[channel.ID]))
		if item != nil {
			resp = append(resp, *item)
		}
	}

	return base.NewPaginationResponse(resp, req.Page, req.PageSize, total), nil
}

func (l *SponsorChannelLogic) GetPublicList(ctx *gin.Context) ([]apiSponsor.ChannelListItemResponse, *xError.Error) {
	channels, xErr := l.repo.channel.PublicList(ctx)
	if xErr != nil {
		return nil, xErr
	}

	sponsorCounts, xErr := l.repo.channel.CountSponsorsByChannelIDs(ctx, collectChannelIDs(channels))
	if xErr != nil {
		return nil, xErr
	}

	resp := make([]apiSponsor.ChannelListItemResponse, 0, len(channels))
	for _, channel := range channels {
		resp = append(resp, buildChannelListItemResponse(&channel, int(sponsorCounts[channel.ID])))
	}
	return resp, nil
}

func parseSponsorChannelID(ctx *gin.Context, idStr string) (int64, *xError.Error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, xError.NewError(ctx, xError.BadRequest, "无效的渠道ID", false)
	}
	return id, nil
}

func collectChannelIDs(channels []entity.SponsorChannel) []int64 {
	ids := make([]int64, 0, len(channels))
	for _, item := range channels {
		ids = append(ids, item.ID)
	}
	return ids
}

func buildChannelEntityResponse(channel *entity.SponsorChannel, sponsorCount int) *apiSponsor.ChannelEntityResponse {
	if channel == nil {
		return nil
	}
	return &apiSponsor.ChannelEntityResponse{
		ID:           channel.ID,
		Name:         channel.Name,
		Icon:         channel.Icon,
		Description:  channel.Description,
		SortOrder:    channel.SortOrder,
		Status:       channel.Status,
		SponsorCount: sponsorCount,
		CreatedAt:    channel.CreatedAt,
		UpdatedAt:    channel.UpdatedAt,
	}
}

func buildChannelListItemResponse(channel *entity.SponsorChannel, sponsorCount int) apiSponsor.ChannelListItemResponse {
	return apiSponsor.ChannelListItemResponse{
		ID:           channel.ID,
		Name:         channel.Name,
		Icon:         channel.Icon,
		SortOrder:    channel.SortOrder,
		Status:       channel.Status,
		SponsorCount: sponsorCount,
	}
}

func stringPointerValue(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
