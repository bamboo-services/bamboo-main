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
	"bamboo-main/internal/model/base"
	"bamboo-main/internal/model/dto"
	"bamboo-main/internal/model/entity"
	"bamboo-main/internal/model/request"
	"bamboo-main/pkg/constants"
	"strconv"

	xError "github.com/bamboo-services/bamboo-base-go/error"
	xUtil "github.com/bamboo-services/bamboo-base-go/utility"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/utility/ctxutil"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LinkLogic 友情链接业务逻辑
type LinkLogic struct {
}

// Add 添加友情链接
func (l *LinkLogic) Add(ctx *gin.Context, req *request.LinkFriendAddReq) (*dto.LinkFriendDetailDTO, *xError.Error) {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)

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

	// 保存到数据库
	err := db.WithContext(ctx).Create(link).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "创建友情链接失败", false, err)
	}

	// 预加载关联数据
	err = db.WithContext(ctx).Preload("GroupFKey").Preload("ColorFKey").First(link, "id = ?", link.ID).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询友情链接失败", false, err)
	}

	return convertLinkFriendToDTO(link), nil
}

// Update 更新友情链接
func (l *LinkLogic) Update(ctx *gin.Context, linkIDStr string, req *request.LinkFriendUpdateReq) (*dto.LinkFriendDetailDTO, *xError.Error) {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)

	// 解析ID
	linkID, err := strconv.ParseInt(linkIDStr, 10, 64)
	if err != nil {
		return nil, xError.NewError(ctx, xError.BadRequest, "无效的友链ID", false)
	}

	// 查找友情链接
	var link entity.LinkFriend
	err = db.WithContext(ctx.Request.Context()).First(&link, "id = ?", linkID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, xError.NewError(ctx, xError.NotFound, "友情链接不存在", false)
		}
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询友情链接失败", false, err)
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

	// 执行更新
	err = db.WithContext(ctx).Updates(&link).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "更新友情链接失败", false, err)
	}

	// 重新查询带关联数据
	err = db.WithContext(ctx).Preload("GroupFKey").Preload("ColorFKey").First(&link, "id = ?", linkID).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询友情链接失败", false, err)
	}

	return convertLinkFriendToDTO(&link), nil
}

// Delete 删除友情链接
func (l *LinkLogic) Delete(ctx *gin.Context, linkIDStr string) *xError.Error {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)

	// 解析ID
	linkID, err := strconv.ParseInt(linkIDStr, 10, 64)
	if err != nil {
		return xError.NewError(ctx, xError.BadRequest, "无效的友链ID", false)
	}

	result := db.WithContext(ctx.Request.Context()).Where("id = ?", linkID).Delete(&entity.LinkFriend{})
	if result.Error != nil {
		return xError.NewError(ctx, xError.DatabaseError, "删除友情链接失败", false, result.Error)
	}
	if result.RowsAffected == 0 {
		return xError.NewError(ctx, xError.NotFound, "友情链接不存在", false)
	}
	return nil
}

// Get 获取友情链接详情
func (l *LinkLogic) Get(ctx *gin.Context, linkIDStr string) (*dto.LinkFriendDetailDTO, *xError.Error) {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)

	// 解析ID
	linkID, err := strconv.ParseInt(linkIDStr, 10, 64)
	if err != nil {
		return nil, xError.NewError(ctx, xError.BadRequest, "无效的友链ID", false)
	}

	var link entity.LinkFriend
	err = db.WithContext(ctx.Request.Context()).Preload("GroupFKey").Preload("ColorFKey").First(&link, "id = ?", linkID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, xError.NewError(ctx, xError.NotFound, "友情链接不存在", false)
		}
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询友情链接失败", false, err)
	}

	return convertLinkFriendToDTO(&link), nil
}

// List 获取友情链接列表
func (l *LinkLogic) List(ctx *gin.Context, req *request.LinkFriendQueryReq) (*base.PaginationResponse[dto.LinkFriendDetailDTO], *xError.Error) {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 10
	}

	// 构建查询条件
	query := db.WithContext(ctx.Request.Context()).Model(&entity.LinkFriend{})

	if req.LinkName != "" {
		query = query.Where("link_name ILIKE ?", "%"+req.LinkName+"%")
	}
	if req.LinkStatus != nil {
		query = query.Where("link_status = ?", *req.LinkStatus)
	}
	if req.LinkFail != nil {
		query = query.Where("is_failure = ?", *req.LinkFail)
	}
	if req.LinkGroupID != 0 {
		query = query.Where("group_id = ?", req.LinkGroupID)
	}

	// 排序
	orderBy := "link_created_at"
	if req.SortBy != "" {
		switch req.SortBy {
		case "created_at":
			orderBy = "link_created_at"
		case "updated_at":
			orderBy = "link_updated_at"
		case "link_order":
			orderBy = "link_order"
		case "link_name":
			orderBy = "link_name"
		}
	}

	sortOrder := "DESC"
	if req.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	query = query.Order(orderBy + " " + sortOrder)

	// 获取总数
	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "统计友情链接数量失败", true, err)
	}

	// 分页查询
	var links []entity.LinkFriend
	offset := (req.Page - 1) * req.PageSize
	err = query.Preload("GroupFKey").Preload("ColorFKey").Offset(offset).Limit(req.PageSize).Find(&links).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询友情链接列表失败", false, err)
	}

	// 转换为 DTO
	var linkDTOs []dto.LinkFriendDetailDTO
	for _, link := range links {
		linkDTOs = append(linkDTOs, *convertLinkFriendToDTO(&link))
	}

	return base.NewPaginationResponse(linkDTOs, req.Page, req.PageSize, total), nil
}

// UpdateStatus 更新友情链接状态
func (l *LinkLogic) UpdateStatus(ctx *gin.Context, linkIDStr string, req *request.LinkFriendStatusReq) *xError.Error {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)

	// 解析ID
	linkID, err := strconv.ParseInt(linkIDStr, 10, 64)
	if err != nil {
		return xError.NewError(ctx, xError.BadRequest, "无效的友链ID", false)
	}

	updates := map[string]interface{}{
		"status":        req.LinkStatus,
		"review_remark": req.LinkReviewRemark,
	}

	result := db.WithContext(ctx.Request.Context()).Model(&entity.LinkFriend{}).Where("id = ?", linkID).Updates(updates)
	if result.Error != nil {
		return xError.NewError(ctx, xError.DatabaseError, "更新友情链接状态失败", false, result.Error)
	}
	if result.RowsAffected == 0 {
		return xError.NewError(ctx, xError.NotFound, "友情链接不存在", false)
	}

	return nil
}

// UpdateFailStatus 更新友情链接失效状态
func (l *LinkLogic) UpdateFailStatus(ctx *gin.Context, linkIDStr string, req *request.LinkFriendFailReq) *xError.Error {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)

	// 解析ID
	linkID, err := strconv.ParseInt(linkIDStr, 10, 64)
	if err != nil {
		return xError.NewError(ctx, xError.BadRequest, "无效的友链ID", false)
	}

	updates := map[string]interface{}{
		"is_failure":  req.LinkFail,
		"fail_reason": req.LinkFailReason,
	}

	result := db.WithContext(ctx.Request.Context()).Model(&entity.LinkFriend{}).Where("id = ?", linkID).Updates(updates)
	if result.Error != nil {
		return xError.NewError(ctx, xError.DatabaseError, "更新友情链接失效状态失败", false, result.Error)
	}
	if result.RowsAffected == 0 {
		return xError.NewError(ctx, xError.NotFound, "友情链接不存在", false)
	}

	return nil
}

// GetPublicLinks 获取公开的友情链接列表
func (l *LinkLogic) GetPublicLinks(ctx *gin.Context, groupIDStr string) ([]dto.LinkFriendDetailDTO, *xError.Error) {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)

	query := db.WithContext(ctx.Request.Context()).Where("status = ? AND is_failure = ?", constants.LinkStatusApproved, constants.LinkFailNormal)

	if groupIDStr != "" {
		groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
		if err == nil {
			query = query.Where("group_id = ?", groupID)
		}
	}

	var links []entity.LinkFriend
	err := query.Preload("GroupFKey").Preload("ColorFKey").Order("sort_order ASC, created_at DESC").Find(&links).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询公开友情链接失败", false, err)
	}

	var linkDTOs []dto.LinkFriendDetailDTO
	for _, link := range links {
		linkDTOs = append(linkDTOs, *convertLinkFriendToDTO(&link))
	}

	return linkDTOs, nil
}

// 辅助函数：将友链实体转换为详细DTO
func convertLinkFriendToDTO(link *entity.LinkFriend) *dto.LinkFriendDetailDTO {
	if link == nil {
		return nil
	}

	linkDTO := &dto.LinkFriendDetailDTO{
		ID:           link.ID,
		Name:         link.Name,
		URL:          link.URL,
		Avatar:       link.Avatar,      // 直接赋值指针 *string → *string
		RSS:          link.RSS,         // 直接赋值指针 *string → *string
		Description:  link.Description, // 直接赋值指针 *string → *string
		Email:        link.Email,       // 直接赋值指针 *string → *string
		GroupID:      link.GroupID,     // 直接赋值指针 *int64 → *int64
		ColorID:      link.ColorID,     // 直接赋值指针 *int64 → *int64
		SortOrder:    link.SortOrder,
		Status:       link.Status,
		StatusText:   getLinkStatusText(link.Status),
		IsFailure:    link.IsFailure,
		FailureText:  getLinkFailText(link.IsFailure),
		FailReason:   link.FailReason,   // 直接赋值指针 *string → *string
		ApplyRemark:  link.ApplyRemark,  // 直接赋值指针 *string → *string
		ReviewRemark: link.ReviewRemark, // 直接赋值指针 *string → *string
		CreatedAt:    link.CreatedAt,
		UpdatedAt:    link.UpdatedAt,
	}

	// 转换关联的分组信息
	if link.GroupFKey != nil {
		linkDTO.GroupInfo = &dto.LinkGroupSimpleDTO{
			ID:   link.GroupFKey.ID,
			Name: link.GroupFKey.Name,
		}
	}

	// 转换关联的颜色信息
	if link.ColorFKey != nil {
		linkDTO.ColorInfo = &dto.LinkColorSimpleDTO{
			ID:    link.ColorFKey.ID,
			Name:  link.ColorFKey.Name,
			Value: link.ColorFKey.Value,
		}
	}

	return linkDTO
}

// getLinkStatusText 获取链接状态文本
func getLinkStatusText(status int) string {
	switch status {
	case constants.LinkStatusPending:
		return "待审核"
	case constants.LinkStatusApproved:
		return "已通过"
	case constants.LinkStatusRejected:
		return "已拒绝"
	default:
		return "未知状态"
	}
}

// getLinkFailText 获取链接失效状态文本
func getLinkFailText(fail int) string {
	switch fail {
	case constants.LinkFailNormal:
		return "正常"
	case constants.LinkFailBroken:
		return "失效"
	default:
		return "未知状态"
	}
}
