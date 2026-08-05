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

package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xSnowflake "github.com/bamboo-services/bamboo-base-go/common/snowflake"
	xCache "github.com/bamboo-services/bamboo-base-go/major/cache"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/pkg/constants"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// LinkRepo 友情链接数据访问层
//
// 收口友情链接实体的 CRUD、分页查询与缓存读写，缓存经 xCache.Manager 统一失效。
type LinkRepo struct {
	db     *gorm.DB
	kc     xCache.KeyCache[string, entity.LinkFriend]
	system *SystemRepo
	log    *xLog.LogNamedLogger
}

// FriendQuery 友情链接分页查询条件
//
// repository 层自有查询模型，由 logic 层从 transport DTO（api/link）转换而来，
// 避免数据访问层直接耦合传输层契约。
type FriendQuery struct {
	Page        int
	PageSize    int
	LinkName    string
	LinkStatus  *int
	LinkFail    *int
	LinkAnomaly *bool
	LinkGroupID xSnowflake.SnowflakeID
	UserID      xSnowflake.SnowflakeID
	SortBy      string
	SortOrder   string
}

// SortAssignment 友链排序赋值条目
//
// repository 层自有模型：承载「目标全局序号 + 目标分组」两项写入值，
// 由 logic 层纯函数 BuildSortAssignments 计算后交由 UpdateSortAndPosition 落库。
type SortAssignment struct {
	ID      xSnowflake.SnowflakeID  // 友链ID
	GroupID *xSnowflake.SnowflakeID // 目标分组ID（nil=未分组）
	Order   int                     // 目标全局排序值
}

// NewLinkRepo 创建 LinkRepo 实例
func NewLinkRepo(db *gorm.DB, m *xCache.Manager) *LinkRepo {
	return &LinkRepo{
		db:     db,
		kc:     xCache.KeyCacheOf[string, entity.LinkFriend](m),
		system: NewSystemRepo(db, m),
		log:    xLog.WithName(xLog.NamedREPO, "LinkRepo"),
	}
}

func (r *LinkRepo) pickDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}

// ensureFancyColor 若友链引用内置炫彩（color_id 为保留 ID），注入虚拟炫彩颜色对象。
//
// 炫彩为系统内置特殊颜色（不落库），数据库 Preload 无法命中该关联，
// 由查询层在返回前补齐 ColorFKey，保证下游渲染无需感知 ID 约定。
func ensureFancyColor(link *entity.LinkFriend) {
	if link.ColorID != nil && *link.ColorID == constants.BuiltinFancyColorID && link.ColorFKey == nil {
		link.ColorFKey = entity.NewFancyColor()
	}
}

// ensureBuiltinGroup 若友链引用内置「已失效」分组（group_id 为保留 ID），注入虚拟分组对象。
//
// 已失效为系统内置语义分组（不落库），数据库 Preload 无法命中该关联，
// 由查询层在返回前补齐 GroupFKey。名字/描述来自调用方读取 bm_system 配置后构造的
// builtinInvalid（nil 时回退默认「已失效」），保证下游渲染无需感知 ID 约定。
func ensureBuiltinGroup(link *entity.LinkFriend, builtinInvalid *entity.LinkGroup) {
	if link.GroupID != nil && link.GroupFKey == nil {
		if entity.IsBuiltinGroupID(*link.GroupID) {
			if builtinInvalid != nil {
				link.GroupFKey = builtinInvalid
			} else {
				link.GroupFKey = entity.NewDefaultBuiltinGroup()
			}
		}
	}
}

// buildBuiltinInvalid 读取内置「已失效」分组配置，读取失败或为空时回退默认值。
//
// 供各查询方法在循环注入前调用一次，避免对每条友链重复读取 system 配置（防 N+1）。
func (r *LinkRepo) buildBuiltinInvalid(ctx context.Context) *entity.LinkGroup {
	builtin, xErr := r.system.BuildBuiltinInvalidGroup(ctx)
	if xErr != nil || builtin == nil {
		return entity.NewDefaultBuiltinGroup()
	}
	return builtin
}

// Create 创建友情链接
func (r *LinkRepo) Create(ctx context.Context, link *entity.LinkFriend, tx *gorm.DB) (*entity.LinkFriend, *xError.Error) {
	r.log.Info(ctx, "Create - 创建友情链接")

	err := r.pickDB(tx).WithContext(ctx).Create(link).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "创建友情链接失败", false, err)
	}

	if cacheErr := r.kc.Set(ctx, constants.RedisLinkFriend.Get(link.ID).String(), link, xCache.WithTTL(15*time.Minute)); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return link, nil
}

// BindUserByEmail 将指定邮箱名下尚未归属的友链绑定到用户
//
// 用于游客提交友链后注册/登录时的归属关联：把所有 user_id 为空且联系邮箱匹配的友链归属到该用户。
func (r *LinkRepo) BindUserByEmail(ctx context.Context, userID xSnowflake.SnowflakeID, email string) *xError.Error {
	r.log.Info(ctx, "BindUserByEmail - 按邮箱绑定友链归属")

	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return nil
	}

	result := r.db.WithContext(ctx).Model(&entity.LinkFriend{}).
		Where("user_id IS NULL AND lower(email) = ?", email).
		Update("user_id", userID)
	if result.Error != nil {
		return xError.NewError(ctx, xError.DatabaseError, "绑定友链归属失败", false, result.Error)
	}
	return nil
}

// Save 保存友情链接（存在则更新）
//
// 经 Omit(clause.Associations) 收敛关联写入：实体可能来自缓存快照（含 GroupFKey/ColorFKey/UserFKey），
// 保留关联会让 GORM Save 用旧快照覆盖关联表并把外键改回旧值。关联引用不参与保存，交由外键字段落库。
func (r *LinkRepo) Save(ctx context.Context, link *entity.LinkFriend, tx *gorm.DB) (*entity.LinkFriend, *xError.Error) {
	r.log.Info(ctx, "Save - 保存友情链接")

	err := r.pickDB(tx).WithContext(ctx).Omit(clause.Associations).Save(link).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "保存友情链接失败", false, err)
	}

	if cacheErr := r.kc.Delete(ctx, constants.RedisLinkFriend.Get(link.ID).String()); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return link, nil
}

// GetByID 根据ID获取友情链接
func (r *LinkRepo) GetByID(ctx context.Context, id xSnowflake.SnowflakeID, withAssociations bool, tx *gorm.DB) (*entity.LinkFriend, bool, *xError.Error) {
	r.log.Info(ctx, "GetByID - 获取友情链接")

	if !withAssociations {
		if link, ok, err := r.kc.Get(ctx, constants.RedisLinkFriend.Get(id).String()); err != nil {
			return nil, false, xError.NewError(ctx, xError.CacheError, "获取友情链接缓存失败", true, err)
		} else if ok {
			return link, true, nil
		}
	}

	query := r.pickDB(tx).WithContext(ctx)
	if withAssociations {
		query = query.Preload("GroupFKey").Preload("ColorFKey")
	}

	var link entity.LinkFriend
	err := query.Where("id = ?", id).First(&link).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, xError.NewError(ctx, xError.DatabaseError, "查询友情链接失败", false, err)
	}

	// 内置虚拟项（不落库）：color_id/group_id 引用保留 ID 时补齐虚拟对象，随缓存一并持久化
	if withAssociations {
		ensureFancyColor(&link)
		ensureBuiltinGroup(&link, r.buildBuiltinInvalid(ctx))
	}

	if cacheErr := r.kc.Set(ctx, constants.RedisLinkFriend.Get(link.ID).String(), &link, xCache.WithTTL(15*time.Minute)); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return &link, true, nil
}

// DeleteByID 删除友情链接
func (r *LinkRepo) DeleteByID(ctx context.Context, id xSnowflake.SnowflakeID, tx *gorm.DB) (bool, *xError.Error) {
	r.log.Info(ctx, "DeleteByID - 删除友情链接")

	result := r.pickDB(tx).WithContext(ctx).Where("id = ?", id).Delete(&entity.LinkFriend{})
	if result.Error != nil {
		return false, xError.NewError(ctx, xError.DatabaseError, "删除友情链接失败", false, result.Error)
	}
	if result.RowsAffected == 0 {
		return false, nil
	}

	if cacheErr := r.kc.Delete(ctx, constants.RedisLinkFriend.Get(id).String()); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return true, nil
}

// List 分页查询友情链接
func (r *LinkRepo) List(ctx context.Context, req *FriendQuery, tx *gorm.DB) ([]entity.LinkFriend, int64, *xError.Error) {
	r.log.Info(ctx, "List - 查询友情链接分页列表")

	query := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkFriend{})

	if req.LinkName != "" {
		query = query.Where("name ILIKE ?", "%"+req.LinkName+"%")
	}
	if req.LinkStatus != nil {
		query = query.Where("status = ?", *req.LinkStatus)
	}
	if req.LinkFail != nil {
		query = query.Where("is_failure = ?", *req.LinkFail)
	}
	if req.LinkAnomaly != nil && *req.LinkAnomaly {
		// 异常：非待审核/已通过（拒绝/下架待审核/已下架）或已失效
		query = query.Where("(status NOT IN (0, 1) OR is_failure = 1)")
	}
	if req.LinkGroupID != 0 {
		query = query.Where("group_id = ?", req.LinkGroupID)
	}
	if req.UserID != 0 {
		query = query.Where("user_id = ?", req.UserID)
	}

	orderBy := "created_at"
	switch req.SortBy {
	case "created_at":
		orderBy = "created_at"
	case "updated_at":
		orderBy = "updated_at"
	case "link_order":
		orderBy = "sort_order"
	case "link_name":
		orderBy = "name"
	}

	sortOrder := "DESC"
	if strings.EqualFold(req.SortOrder, "asc") {
		sortOrder = "ASC"
	}

	query = query.Order(orderBy + " " + sortOrder)

	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, xError.NewError(ctx, xError.DatabaseError, "统计友情链接数量失败", true, err)
	}

	offset := (req.Page - 1) * req.PageSize
	var links []entity.LinkFriend
	err = query.Preload("GroupFKey").Preload("ColorFKey").Offset(offset).Limit(req.PageSize).Find(&links).Error
	if err != nil {
		return nil, 0, xError.NewError(ctx, xError.DatabaseError, "查询友情链接列表失败", false, err)
	}

	// 内置虚拟项（不落库）：color_id/group_id 引用保留 ID 时补齐虚拟对象
	builtinInvalid := r.buildBuiltinInvalid(ctx)
	for i := range links {
		ensureFancyColor(&links[i])
		ensureBuiltinGroup(&links[i], builtinInvalid)
	}

	return links, total, nil
}

// UpdateStatusByID 更新友情链接的审核状态
func (r *LinkRepo) UpdateStatusByID(ctx context.Context, id xSnowflake.SnowflakeID, status int, reviewRemark string, tx *gorm.DB) (bool, *xError.Error) {
	r.log.Info(ctx, "UpdateStatusByID - 更新友情链接状态")

	updates := map[string]any{
		"status":        status,
		"review_remark": reviewRemark,
	}

	result := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkFriend{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return false, xError.NewError(ctx, xError.DatabaseError, "更新友情链接状态失败", false, result.Error)
	}
	if result.RowsAffected == 0 {
		return false, nil
	}

	if cacheErr := r.kc.Delete(ctx, constants.RedisLinkFriend.Get(id).String()); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return true, nil
}

// UpdateFailureByID 更新友情链接的失效状态
//
// invalidGroupID 为内置「已失效」分组的保留 ID：标记失效时自动归入该分组，
// 恢复时若当前分组即已失效分组则清空（未分组），否则保持原分组不变。
func (r *LinkRepo) UpdateFailureByID(ctx context.Context, id xSnowflake.SnowflakeID, isFailure int, failReason string, invalidGroupID *xSnowflake.SnowflakeID, tx *gorm.DB) (bool, *xError.Error) {
	r.log.Info(ctx, "UpdateFailureByID - 更新友情链接失效状态")

	updates := buildFailureUpdates(isFailure, failReason, invalidGroupID)

	result := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkFriend{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return false, xError.NewError(ctx, xError.DatabaseError, "更新友情链接失效状态失败", false, result.Error)
	}
	if result.RowsAffected == 0 {
		return false, nil
	}

	if cacheErr := r.kc.Delete(ctx, constants.RedisLinkFriend.Get(id).String()); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}

	return true, nil
}

// buildFailureUpdates 组装失效状态更新字段（纯函数，便于表驱动单测）。
//
// invalidGroupID 为空时不操作分组；标记失效（isFailure 为失效标志）时归入内置
// 「已失效」分组，恢复时以单条 CASE 原子判断「当前分组为已失效分组则清空」，
// 避免读-改-写竞态。
func buildFailureUpdates(isFailure int, failReason string, invalidGroupID *xSnowflake.SnowflakeID) map[string]any {
	updates := map[string]any{
		"is_failure":  isFailure,
		"fail_reason": failReason,
	}

	if invalidGroupID != nil {
		if isFailure == constants.LinkFailBroken {
			updates["group_id"] = *invalidGroupID
		} else {
			updates["group_id"] = gorm.Expr(
				"CASE WHEN group_id = ? THEN NULL ELSE group_id END", *invalidGroupID,
			)
		}
	}

	return updates
}

// ListPublic 查询公开的友情链接列表
//
// 排序语义「数字越小权重越大」：章节序遵循分组排序值（bm_link_group.sort_order ASC，
// 未分组 NULL 经 NULLS LAST 置底），章内遵循友链排序值（bm_link_friend.sort_order ASC），
// 同序回退创建时间倒序。LEFT JOIN 后 status 于两表同名，条件列一律带表限定；
// 表名含 bm_ 前缀且为单数（xOptionDB.WithTablePrefix + SingularTable，见 main.go）。
func (r *LinkRepo) ListPublic(ctx context.Context, groupID *xSnowflake.SnowflakeID, approvedStatus int, normalFail int, tx *gorm.DB) ([]entity.LinkFriend, *xError.Error) {
	r.log.Info(ctx, "ListPublic - 查询公开友情链接")

	query := r.pickDB(tx).WithContext(ctx).
		Joins("LEFT JOIN bm_link_group ON bm_link_group.id = bm_link_friend.group_id").
		Where("bm_link_friend.status = ? AND bm_link_friend.is_failure = ?", approvedStatus, normalFail)

	if groupID != nil {
		query = query.Where("bm_link_friend.group_id = ?", *groupID)
	}

	var links []entity.LinkFriend
	err := query.Preload("GroupFKey").Preload("ColorFKey").
		Order("bm_link_group.sort_order ASC NULLS LAST, bm_link_friend.sort_order ASC, bm_link_friend.created_at DESC").
		Find(&links).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询公开友情链接失败", false, err)
	}

	// 内置虚拟项（不落库）：color_id/group_id 引用保留 ID 时补齐虚拟对象
	builtinInvalid := r.buildBuiltinInvalid(ctx)
	for i := range links {
		ensureFancyColor(&links[i])
		ensureBuiltinGroup(&links[i], builtinInvalid)
	}

	return links, nil
}

// ListFailed 查询已失效的公开友链（is_failure 为失效标志且状态为已通过）。
//
// 供公开「已失效」章节独立展示，按友链排序值升序、创建时间倒序，
// 返回前补齐内置虚拟关联，与 ListPublic 共用同一套注入逻辑。
func (r *LinkRepo) ListFailed(ctx context.Context, approvedStatus int, brokenFail int, tx *gorm.DB) ([]entity.LinkFriend, *xError.Error) {
	r.log.Info(ctx, "ListFailed - 查询已失效公开友链")

	var links []entity.LinkFriend
	err := r.pickDB(tx).WithContext(ctx).
		Where("is_failure = ? AND status = ?", brokenFail, approvedStatus).
		Preload("GroupFKey").Preload("ColorFKey").
		Order("sort_order ASC, created_at DESC").
		Find(&links).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询已失效友链失败", false, err)
	}

	// 内置虚拟项（不落库）：color_id/group_id 引用保留 ID 时补齐虚拟对象
	builtinInvalid := r.buildBuiltinInvalid(ctx)
	for i := range links {
		ensureFancyColor(&links[i])
		ensureBuiltinGroup(&links[i], builtinInvalid)
	}

	return links, nil
}

// GetByIDs 按 ID 列表批量查询友链（无预加载、不走单条缓存，供排序前存在性校验）
func (r *LinkRepo) GetByIDs(ctx context.Context, ids []xSnowflake.SnowflakeID, tx *gorm.DB) ([]entity.LinkFriend, *xError.Error) {
	r.log.Info(ctx, "GetByIDs - 批量查询友情链接")

	var links []entity.LinkFriend
	err := r.pickDB(tx).WithContext(ctx).Where("id IN ?", ids).Find(&links).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "批量查询友情链接失败", false, err)
	}

	return links, nil
}

// linkSortRowTemplate 友链批量排序 VALUES 行模板（id, sort_order, group_id）
//
// 占位符逐项显式 cast 目标类型：id/group_id 为 bigint（xSnowflake.SnowflakeID=int64），
// sort_order 为 int；避免 PostgreSQL 将纯参数列默认推断为 text 触发 SQLSTATE 42883。
// 占位符数量必须与 SortAssignment 字段一一对应，由 TestLinkSortRowTemplate 兜底。
const linkSortRowTemplate = "(?::bigint, ?::int, ?::bigint)"

// UpdateSortAndPosition 批量更新友链的全局排序值与分组归属
//
// 单条 UPDATE ... FROM (VALUES ...) 一次落库全部条目，消除逐行 N+1 往返；
// 排序视为内容更新，统一刷新 updated_at（与分组/颜色重排行为对齐）；
// RowsAffected==len(items) 依赖 PostgreSQL「匹配行」语义（值未变化也计数），
// 作为原子存在性兜底，可同时捕获 GetByIDs 预检与写入之间的并发删除（TOCTOU）；
// VALUES 语法与 ::cast 为 PostgreSQL 专属，勿在其他驱动下使用。
// 写后逐条失效单条缓存。
func (r *LinkRepo) UpdateSortAndPosition(ctx context.Context, items []SortAssignment, tx *gorm.DB) *xError.Error {
	r.log.Info(ctx, "UpdateSortAndPosition - 批量更新友链排序与位置")

	if len(items) == 0 {
		return nil
	}

	db := r.pickDB(tx).WithContext(ctx)
	values := make([]string, 0, len(items))
	args := make([]any, 0, len(items)*3)
	for _, it := range items {
		values = append(values, linkSortRowTemplate)
		args = append(args, it.ID, it.Order, it.GroupID)
	}

	result := db.Exec(
		fmt.Sprintf(`UPDATE bm_link_friend
			SET sort_order = v.sort_order, group_id = v.group_id, updated_at = now()
			FROM (VALUES %s) AS v(id, sort_order, group_id)
			WHERE bm_link_friend.id = v.id`, strings.Join(values, ",")),
		args...,
	)
	if result.Error != nil {
		return xError.NewError(ctx, xError.DatabaseError, "更新友链排序失败", false, result.Error)
	}
	if result.RowsAffected != int64(len(items)) {
		return xError.NewError(ctx, xError.NotFound, "友情链接不存在", false)
	}

	for _, it := range items {
		if cacheErr := r.kc.Delete(ctx, constants.RedisLinkFriend.Get(it.ID).String()); cacheErr != nil {
			r.log.Warn(ctx, cacheErr.Error())
		}
	}

	return nil
}

// ListApprovedForScreenshot 查询全部已通过且未失效的友链（按排序与创建时间，供截图全量入队）
func (r *LinkRepo) ListApprovedForScreenshot(ctx context.Context, tx *gorm.DB) ([]entity.LinkFriend, *xError.Error) {
	r.log.Info(ctx, "ListApprovedForScreenshot - 查询待截图友链")

	var links []entity.LinkFriend
	err := r.pickDB(tx).WithContext(ctx).
		Where("status = ? AND is_failure = ?", constants.LinkStatusApproved, constants.LinkFailNormal).
		Order("sort_order ASC, created_at DESC").
		Find(&links).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询待截图友链失败", false, err)
	}

	return links, nil
}

// UpdateScreenshot 更新友链站点截图信息（URL 与最近截图时间）
func (r *LinkRepo) UpdateScreenshot(ctx context.Context, id xSnowflake.SnowflakeID, url string, at time.Time, tx *gorm.DB) *xError.Error {
	r.log.Info(ctx, "UpdateScreenshot - 更新友链截图信息")

	result := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkFriend{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"screenshot_url": url,
			"screenshot_at":  at,
		})
	if result.Error != nil {
		return xError.NewError(ctx, xError.DatabaseError, "更新友链截图信息失败", false, result.Error)
	}
	if result.RowsAffected == 0 {
		return xError.NewError(ctx, xError.NotFound, "友情链接不存在", false)
	}

	if cacheErr := r.kc.Delete(ctx, constants.RedisLinkFriend.Get(id).String()); cacheErr != nil {
		r.log.Warn(ctx, cacheErr.Error())
	}
	return nil
}

// CountByStatus 按审核状态统计友情链接数量；status 为负数时统计全部
func (r *LinkRepo) CountByStatus(ctx context.Context, status int, tx *gorm.DB) (int64, *xError.Error) {
	query := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkFriend{})
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, xError.NewError(ctx, xError.DatabaseError, "统计友情链接数量失败", true, err)
	}
	return count, nil
}

// CountAnomaly 统计异常友链数量（status 非 0/1 或已失效）
func (r *LinkRepo) CountAnomaly(ctx context.Context, tx *gorm.DB) (int64, *xError.Error) {
	var count int64
	err := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkFriend{}).
		Where("(status NOT IN (0, 1) OR is_failure = 1)").
		Count(&count).Error
	if err != nil {
		return 0, xError.NewError(ctx, xError.DatabaseError, "统计异常友链数量失败", true, err)
	}
	return count, nil
}

// CountByGroupID 统计指定分组下友情链接数量
func (r *LinkRepo) CountByGroupID(ctx context.Context, groupID xSnowflake.SnowflakeID, tx *gorm.DB) (int64, *xError.Error) {
	var count int64
	err := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkFriend{}).Where("group_id = ?", groupID).Count(&count).Error
	if err != nil {
		return 0, xError.NewError(ctx, xError.DatabaseError, "查询关联友链失败", false, err)
	}
	return count, nil
}

// CountByColorID 统计指定颜色下友情链接数量
func (r *LinkRepo) CountByColorID(ctx context.Context, colorID xSnowflake.SnowflakeID, tx *gorm.DB) (int64, *xError.Error) {
	var count int64
	err := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkFriend{}).Where("color_id = ?", colorID).Count(&count).Error
	if err != nil {
		return 0, xError.NewError(ctx, xError.DatabaseError, "查询关联友链失败", false, err)
	}
	return count, nil
}

// ListByGroupID 根据分组ID查询友情链接列表
func (r *LinkRepo) ListByGroupID(ctx context.Context, groupID xSnowflake.SnowflakeID, limit int, tx *gorm.DB) ([]entity.LinkFriend, *xError.Error) {
	query := r.pickDB(tx).WithContext(ctx).Where("group_id = ?", groupID)
	if limit > 0 {
		query = query.Limit(limit)
	}

	var links []entity.LinkFriend
	err := query.Find(&links).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询冲突友链失败", false, err)
	}

	return links, nil
}

// ListByColorID 根据颜色ID查询友情链接列表
func (r *LinkRepo) ListByColorID(ctx context.Context, colorID xSnowflake.SnowflakeID, limit int, tx *gorm.DB) ([]entity.LinkFriend, *xError.Error) {
	query := r.pickDB(tx).WithContext(ctx).Where("color_id = ?", colorID)
	if limit > 0 {
		query = query.Limit(limit)
	}

	var links []entity.LinkFriend
	err := query.Find(&links).Error
	if err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "查询冲突友链失败", false, err)
	}

	return links, nil
}

// ClearGroupID 清空指定分组下友情链接的分组关联
func (r *LinkRepo) ClearGroupID(ctx context.Context, groupID xSnowflake.SnowflakeID, tx *gorm.DB) *xError.Error {
	result := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkFriend{}).Where("group_id = ?", groupID).Update("group_id", nil)
	if result.Error != nil {
		return xError.NewError(ctx, xError.DatabaseError, "清空友链分组关联失败", false, result.Error)
	}
	return nil
}

// ClearColorID 清空指定颜色下友情链接的颜色关联
func (r *LinkRepo) ClearColorID(ctx context.Context, colorID xSnowflake.SnowflakeID, tx *gorm.DB) *xError.Error {
	result := r.pickDB(tx).WithContext(ctx).Model(&entity.LinkFriend{}).Where("color_id = ?", colorID).Update("color_id", nil)
	if result.Error != nil {
		return xError.NewError(ctx, xError.DatabaseError, "清空友链颜色关联失败", false, result.Error)
	}
	return nil
}
