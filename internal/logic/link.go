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
	"github.com/bamboo-services/bamboo-main/internal/models/dto"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
	ctxUtil "github.com/bamboo-services/bamboo-main/pkg/util/ctx"

	xError "github.com/bamboo-services/bamboo-base-go/error"
	xLog "github.com/bamboo-services/bamboo-base-go/log"
	xUtil "github.com/bamboo-services/bamboo-base-go/utility"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/utility/ctxutil"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LinkLogic 友情链接业务逻辑
type LinkLogic struct {
	logic
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
	}
}

// Add 添加友情链接
func (l *LinkLogic) Add(ctx *gin.Context, req *apiLink.FriendAddRequest) (*dto.LinkFriendDetailDTO, *xError.Error) {
	// 获取数据库连接
	db := xCtxUtil.MustGetDB(ctx)

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
	err := db.Create(link).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "创建友情链接失败", false, err)
	}

	// 预加载关联数据
	err = db.Preload("GroupFKey").Preload("ColorFKey").First(link, "id = ?", link.ID).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询友情链接失败", false, err)
	}

	// 发送邮件通知管理员（异步，不阻断主流程）
	go l.sendApplyNotification(ctx, link)

	return convertLinkFriendToDTO(link), nil
}

// Update 更新友情链接
func (l *LinkLogic) Update(ctx *gin.Context, linkIDStr string, req *apiLink.FriendUpdateRequest) (*dto.LinkFriendDetailDTO, *xError.Error) {
	// 获取数据库连接
	db := xCtxUtil.MustGetDB(ctx)

	// 解析ID
	linkID, err := strconv.ParseInt(linkIDStr, 10, 64)
	if err != nil {
		return nil, xError.NewError(ctx, xError.BadRequest, "无效的友链ID", false)
	}

	// 查找友情链接
	var link entity.LinkFriend
	err = db.First(&link, "id = ?", linkID).Error
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
	err = db.Updates(&link).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "更新友情链接失败", false, err)
	}

	// 重新查询带关联数据
	err = db.Preload("GroupFKey").Preload("ColorFKey").First(&link, "id = ?", linkID).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询友情链接失败", false, err)
	}

	return convertLinkFriendToDTO(&link), nil
}

// Delete 删除友情链接
func (l *LinkLogic) Delete(ctx *gin.Context, linkIDStr string) *xError.Error {
	// 获取数据库连接
	db := xCtxUtil.MustGetDB(ctx)

	// 解析ID
	linkID, err := strconv.ParseInt(linkIDStr, 10, 64)
	if err != nil {
		return xError.NewError(ctx, xError.BadRequest, "无效的友链ID", false)
	}

	result := db.Where("id = ?", linkID).Delete(&entity.LinkFriend{})
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
	db := xCtxUtil.MustGetDB(ctx)

	// 解析ID
	linkID, err := strconv.ParseInt(linkIDStr, 10, 64)
	if err != nil {
		return nil, xError.NewError(ctx, xError.BadRequest, "无效的友链ID", false)
	}

	var link entity.LinkFriend
	err = db.Preload("GroupFKey").Preload("ColorFKey").First(&link, "id = ?", linkID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, xError.NewError(ctx, xError.NotFound, "友情链接不存在", false)
		}
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询友情链接失败", false, err)
	}

	return convertLinkFriendToDTO(&link), nil
}

// List 获取友情链接列表
func (l *LinkLogic) List(ctx *gin.Context, req *apiLink.FriendQueryRequest) (*base.PaginationResponse[dto.LinkFriendDetailDTO], *xError.Error) {
	// 获取数据库连接
	db := xCtxUtil.MustGetDB(ctx)

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 10
	}

	// 构建查询条件
	query := db.Model(&entity.LinkFriend{})

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
func (l *LinkLogic) UpdateStatus(ctx *gin.Context, linkIDStr string, req *apiLink.FriendStatusRequest) *xError.Error {
	// 获取数据库连接
	db := xCtxUtil.MustGetDB(ctx)

	// 解析ID
	linkID, err := strconv.ParseInt(linkIDStr, 10, 64)
	if err != nil {
		return xError.NewError(ctx, xError.BadRequest, "无效的友链ID", false)
	}

	// 先查询友链信息（用于发送邮件通知）
	var link entity.LinkFriend
	err = db.First(&link, "id = ?", linkID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return xError.NewError(ctx, xError.NotFound, "友情链接不存在", false)
		}
		return xError.NewError(ctx, xError.DatabaseError, "查询友情链接失败", false, err)
	}

	updates := map[string]interface{}{
		"status":        req.LinkStatus,
		"review_remark": req.LinkReviewRemark,
	}

	result := db.Model(&entity.LinkFriend{}).Where("id = ?", linkID).Updates(updates)
	if result.Error != nil {
		return xError.NewError(ctx, xError.DatabaseError, "更新友情链接状态失败", false, result.Error)
	}
	if result.RowsAffected == 0 {
		return xError.NewError(ctx, xError.NotFound, "友情链接不存在", false)
	}

	// 发送审核结果邮件通知（异步，不阻断主流程）
	go l.sendStatusNotification(ctx, &link, req.LinkStatus, req.LinkReviewRemark)

	return nil
}

// UpdateFailStatus 更新友情链接失效状态
func (l *LinkLogic) UpdateFailStatus(ctx *gin.Context, linkIDStr string, req *apiLink.FriendFailRequest) *xError.Error {
	// 获取数据库连接
	db := xCtxUtil.MustGetDB(ctx)

	// 解析ID
	linkID, err := strconv.ParseInt(linkIDStr, 10, 64)
	if err != nil {
		return xError.NewError(ctx, xError.BadRequest, "无效的友链ID", false)
	}

	updates := map[string]interface{}{
		"is_failure":  req.LinkFail,
		"fail_reason": req.LinkFailReason,
	}

	result := db.Model(&entity.LinkFriend{}).Where("id = ?", linkID).Updates(updates)
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
	db := xCtxUtil.MustGetDB(ctx)

	query := db.Where("status = ? AND is_failure = ?", constants.LinkStatusApproved, constants.LinkFailNormal)

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
		Avatar:       link.Avatar,
		RSS:          link.RSS,
		Description:  link.Description,
		Email:        link.Email,
		GroupID:      link.GroupID,
		ColorID:      link.ColorID,
		SortOrder:    link.SortOrder,
		Status:       link.Status,
		StatusText:   getLinkStatusText(link.Status),
		IsFailure:    link.IsFailure,
		FailureText:  getLinkFailText(link.IsFailure),
		FailReason:   link.FailReason,
		ApplyRemark:  link.ApplyRemark,
		ReviewRemark: link.ReviewRemark,
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
			ID:           link.ColorFKey.ID,
			Name:         link.ColorFKey.Name,
			Type:         link.ColorFKey.Type,
			PrimaryColor: link.ColorFKey.PrimaryColor,
			SubColor:     link.ColorFKey.SubColor,
			HoverColor:   link.ColorFKey.HoverColor,
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
