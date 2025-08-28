package logic

import (
	"fmt"

	"bamboo-main/internal/model/dto"
	"bamboo-main/internal/model/entity"
	"bamboo-main/internal/model/request"
	"bamboo-main/pkg/constants"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/utility/ctxutil"
	"gorm.io/gorm"
)

// stringToPtr 字符串转指针辅助函数
func stringToPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// LinkFriendLogic 友情链接业务逻辑
type LinkFriendLogic struct {
}

// NewLinkFriendLogic 创建友情链接业务逻辑实例
func NewLinkFriendLogic() *LinkFriendLogic {
	return &LinkFriendLogic{}
}

// Add 添加友情链接
func (l *LinkFriendLogic) Add(ctx *gin.Context, req *request.LinkFriendAddReq) (*dto.LinkFriendDTO, error) {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)
	
	// 创建友情链接实体
	link := &entity.LinkFriend{
		Name:         req.LinkName,
		URL:          req.LinkURL,
		Avatar:       stringToPtr(req.LinkAvatar),
		RSS:          stringToPtr(req.LinkRSS),
		Description:  stringToPtr(req.LinkDesc),
		Email:        stringToPtr(req.LinkEmail),
		SortOrder:    req.LinkOrder,
		Status:       constants.LinkStatusPending, // 默认待审核
		IsFailure:    constants.LinkFailNormal,    // 默认正常
		ApplyRemark:  stringToPtr(req.LinkApplyRemark),
	}

	// 设置UUID外键
	if req.LinkGroupUUID != "" {
		if groupUUID, err := uuid.Parse(req.LinkGroupUUID); err == nil {
			link.GroupUUID = &groupUUID
		}
	}
	if req.LinkColorUUID != "" {
		if colorUUID, err := uuid.Parse(req.LinkColorUUID); err == nil {
			link.ColorUUID = &colorUUID
		}
	}

	// 保存到数据库
	err := db.WithContext(ctx).Create(link).Error
	if err != nil {
		return nil, fmt.Errorf("创建友情链接失败: %w", err)
	}

	// 预加载关联数据
	err = db.WithContext(ctx).Preload("GroupFKey").Preload("ColorFKey").First(link, "uuid = ?", link.UUID).Error
	if err != nil {
		return nil, fmt.Errorf("查询友情链接失败: %w", err)
	}

	return convertLinkFriendToDTO(link), nil
}

// Update 更新友情链接
func (l *LinkFriendLogic) Update(ctx *gin.Context, linkUUID string, req *request.LinkFriendUpdateReq) (*dto.LinkFriendDTO, error) {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)
	
	// 查找友情链接
	var link entity.LinkFriend
	err := db.WithContext(ctx.Request.Context()).First(&link, "link_uuid = ?", linkUUID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("友情链接不存在")
		}
		return nil, fmt.Errorf("查询友情链接失败: %w", err)
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.LinkName != "" {
		updates["link_name"] = req.LinkName
	}
	if req.LinkURL != "" {
		updates["link_url"] = req.LinkURL
	}
	if req.LinkAvatar != "" {
		updates["link_avatar"] = req.LinkAvatar
	}
	if req.LinkRSS != "" {
		updates["link_rss"] = req.LinkRSS
	}
	if req.LinkDesc != "" {
		updates["link_desc"] = req.LinkDesc
	}
	if req.LinkEmail != "" {
		updates["link_email"] = req.LinkEmail
	}
	if req.LinkGroupUUID != "" {
		updates["link_group_uuid"] = req.LinkGroupUUID
	}
	if req.LinkColorUUID != "" {
		updates["link_color_uuid"] = req.LinkColorUUID
	}
	if req.LinkOrder != nil {
		updates["link_order"] = *req.LinkOrder
	}
	if req.LinkApplyRemark != "" {
		updates["link_apply_remark"] = req.LinkApplyRemark
	}

	// 执行更新
	err = db.WithContext(ctx).Model(&link).Updates(updates).Error
	if err != nil {
		return nil, fmt.Errorf("更新友情链接失败: %w", err)
	}

	// 重新查询带关联数据
	err = db.WithContext(ctx).Preload("LinkGroup").Preload("LinkColor").First(&link, "link_uuid = ?", linkUUID).Error
	if err != nil {
		return nil, fmt.Errorf("查询友情链接失败: %w", err)
	}

	return convertLinkFriendToDTO(&link), nil
}

// Delete 删除友情链接
func (l *LinkFriendLogic) Delete(ctx *gin.Context, linkUUID string) error {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)
	
	result := db.WithContext(ctx.Request.Context()).Where("link_uuid = ?", linkUUID).Delete(&entity.LinkFriend{})
	if result.Error != nil {
		return fmt.Errorf("删除友情链接失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("友情链接不存在")
	}
	return nil
}

// Get 获取友情链接详情
func (l *LinkFriendLogic) Get(ctx *gin.Context, linkUUID string) (*dto.LinkFriendDTO, error) {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)
	
	var link entity.LinkFriend
	err := db.WithContext(ctx.Request.Context()).Preload("LinkGroup").Preload("LinkColor").First(&link, "link_uuid = ?", linkUUID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("友情链接不存在")
		}
		return nil, fmt.Errorf("查询友情链接失败: %w", err)
	}

	return convertLinkFriendToDTO(&link), nil
}

// List 获取友情链接列表
func (l *LinkFriendLogic) List(ctx *gin.Context, req *request.LinkFriendQueryReq) (*dto.PaginationDTO[dto.LinkFriendDTO], error) {
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
		query = query.Where("link_fail = ?", *req.LinkFail)
	}
	if req.LinkGroupUUID != "" {
		query = query.Where("link_group_uuid = ?", req.LinkGroupUUID)
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
	if req.SortOrder != "" && req.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	query = query.Order(orderBy + " " + sortOrder)

	// 获取总数
	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, fmt.Errorf("统计友情链接数量失败: %w", err)
	}

	// 分页查询
	var links []entity.LinkFriend
	offset := (req.Page - 1) * req.PageSize
	err = query.Preload("LinkGroup").Preload("LinkColor").Offset(offset).Limit(req.PageSize).Find(&links).Error
	if err != nil {
		return nil, fmt.Errorf("查询友情链接列表失败: %w", err)
	}

	// 转换为 DTO
	var linkDTOs []dto.LinkFriendDTO
	for _, link := range links {
		linkDTOs = append(linkDTOs, *convertLinkFriendToDTO(&link))
	}

	return &dto.PaginationDTO[dto.LinkFriendDTO]{
		Data:       linkDTOs,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: calculateTotalPages(total, req.PageSize),
	}, nil
}

// UpdateStatus 更新友情链接状态
func (l *LinkFriendLogic) UpdateStatus(ctx *gin.Context, linkUUID string, req *request.LinkFriendStatusReq) error {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)
	
	updates := map[string]interface{}{
		"link_status":        req.LinkStatus,
		"link_review_remark": req.LinkReviewRemark,
	}

	result := db.WithContext(ctx.Request.Context()).Model(&entity.LinkFriend{}).Where("link_uuid = ?", linkUUID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("更新友情链接状态失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("友情链接不存在")
	}

	return nil
}

// UpdateFailStatus 更新友情链接失效状态
func (l *LinkFriendLogic) UpdateFailStatus(ctx *gin.Context, linkUUID string, req *request.LinkFriendFailReq) error {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)
	
	updates := map[string]interface{}{
		"link_fail":        req.LinkFail,
		"link_fail_reason": req.LinkFailReason,
	}

	result := db.WithContext(ctx.Request.Context()).Model(&entity.LinkFriend{}).Where("link_uuid = ?", linkUUID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("更新友情链接失效状态失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("友情链接不存在")
	}

	return nil
}

// GetPublicLinks 获取公开的友情链接列表
func (l *LinkFriendLogic) GetPublicLinks(ctx *gin.Context, groupUUID string) ([]dto.LinkFriendDTO, error) {
	// 获取数据库连接
	db := xCtxUtil.GetDB(ctx)
	
	query := db.WithContext(ctx.Request.Context()).Where("link_status = ? AND link_fail = ?", constants.LinkStatusApproved, constants.LinkFailNormal)

	if groupUUID != "" {
		query = query.Where("link_group_uuid = ?", groupUUID)
	}

	var links []entity.LinkFriend
	err := query.Preload("LinkGroup").Preload("LinkColor").Order("link_order ASC, link_created_at DESC").Find(&links).Error
	if err != nil {
		return nil, fmt.Errorf("查询公开友情链接失败: %w", err)
	}

	var linkDTOs []dto.LinkFriendDTO
	for _, link := range links {
		linkDTOs = append(linkDTOs, *convertLinkFriendToDTO(&link))
	}

	return linkDTOs, nil
}

// 辅助函数：将友链实体转换为详细DTO
func convertLinkFriendToDTO(link *entity.LinkFriend) *dto.LinkFriendDTO {
	if link == nil {
		return nil
	}

	linkDTO := &dto.LinkFriendDTO{
		UUID:          link.UUID.String(),
		Name:          link.Name,
		URL:           link.URL,
		Avatar:        safeStringValue(link.Avatar),
		RSS:           safeStringValue(link.RSS),
		Description:   safeStringValue(link.Description),
		Email:         safeStringValue(link.Email),
		GroupUUID:     safeUUIDValue(link.GroupUUID),
		ColorUUID:     safeUUIDValue(link.ColorUUID),
		SortOrder:     link.SortOrder,
		Status:        link.Status,
		StatusText:    getLinkStatusText(link.Status),
		IsFailure:     link.IsFailure,
		FailureText:   getLinkFailText(link.IsFailure),
		FailReason:    safeStringValue(link.FailReason),
		ApplyRemark:   safeStringValue(link.ApplyRemark),
		ReviewRemark:  safeStringValue(link.ReviewRemark),
		CreatedAt:     link.CreatedAt,
		UpdatedAt:     link.UpdatedAt,
	}

	// 转换关联的分组信息
	if link.GroupFKey != nil {
		linkDTO.GroupInfo = &dto.LinkGroupSimpleDTO{
			UUID: link.GroupFKey.UUID.String(),
			Name: link.GroupFKey.Name,
		}
	}

	// 转换关联的颜色信息
	if link.ColorFKey != nil {
		linkDTO.ColorInfo = &dto.LinkColorSimpleDTO{
			UUID:  link.ColorFKey.UUID.String(),
			Name:  link.ColorFKey.Name,
			Value: link.ColorFKey.Value,
		}
	}

	return linkDTO
}

// safeUUIDValue 安全转换UUID指针为字符串
func safeUUIDValue(u *uuid.UUID) string {
	if u == nil {
		return ""
	}
	return u.String()
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

// calculateTotalPages 计算总页数
func calculateTotalPages(total int64, pageSize int) int {
	if pageSize <= 0 {
		return 0
	}
	return int((total + int64(pageSize) - 1) / int64(pageSize))
}