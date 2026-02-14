package logic

import (
	"context"
	"strconv"

	xError "github.com/bamboo-services/bamboo-base-go/error"
	xLog "github.com/bamboo-services/bamboo-base-go/log"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/utility/ctxutil"
	apiLinkGroup "github.com/bamboo-services/bamboo-main/api/link"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/internal/models/base"
	"github.com/bamboo-services/bamboo-main/internal/repository"
	"github.com/gin-gonic/gin"
)

type linkGroupRepo struct {
	group *repository.LinkGroupRepo
	link  *repository.LinkRepo
}

type LinkGroupLogic struct {
	logic
	repo linkGroupRepo
}

func NewLinkGroupLogic(ctx context.Context) *LinkGroupLogic {
	db := xCtxUtil.MustGetDB(ctx)
	rdb := xCtxUtil.MustGetRDB(ctx)

	return &LinkGroupLogic{
		logic: logic{
			db:  db,
			rdb: rdb,
			log: xLog.WithName(xLog.NamedLOGC, "LinkGroupLogic"),
		},
		repo: linkGroupRepo{
			group: repository.NewLinkGroupRepo(db, rdb),
			link:  repository.NewLinkRepo(db, rdb),
		},
	}
}

func (l *LinkGroupLogic) Add(ctx *gin.Context, req *apiLinkGroup.GroupAddRequest) (*entity.LinkGroup, *xError.Error) {
	maxSort, xErr := l.repo.group.GetMaxSortOrder(ctx, nil)
	if xErr != nil {
		return nil, xErr
	}

	group := &entity.LinkGroup{
		Name:        req.GroupName,
		Description: &req.GroupDesc,
		SortOrder:   maxSort + 1,
		Status:      true,
	}

	_, xErr = l.repo.group.Create(ctx, group, nil)
	if xErr != nil {
		return nil, xErr
	}

	reloaded, found, xErr := l.repo.group.GetByID(ctx, group.ID, true, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "友链分组不存在", false)
	}

	return reloaded, nil
}

func (l *LinkGroupLogic) Update(ctx *gin.Context, groupIDStr string, req *apiLinkGroup.GroupUpdateRequest) (*entity.LinkGroup, *xError.Error) {
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		return nil, xError.NewError(ctx, xError.BadRequest, "无效的分组ID", false)
	}

	group, found, xErr := l.repo.group.GetByID(ctx, groupID, false, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "友链分组不存在", false)
	}

	if req.GroupName != "" {
		group.Name = req.GroupName
	}
	if req.GroupDesc != "" {
		group.Description = &req.GroupDesc
	}
	if req.GroupOrder != nil {
		group.SortOrder = *req.GroupOrder
	}
	if req.GroupStatus != nil {
		group.Status = *req.GroupStatus == 1
	}

	_, xErr = l.repo.group.Save(ctx, group, nil)
	if xErr != nil {
		return nil, xErr
	}

	reloaded, found, xErr := l.repo.group.GetByID(ctx, groupID, true, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "友链分组不存在", false)
	}

	return reloaded, nil
}

func (l *LinkGroupLogic) UpdateSort(ctx *gin.Context, req *apiLinkGroup.GroupSortRequest) *xError.Error {
	tx := l.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	startSort := 0
	if req.SortOrder != nil && *req.SortOrder > 0 {
		startSort = *req.SortOrder
	}

	if xErr := l.repo.group.UpdateSortByIDs(ctx, req.GroupIDs, startSort, tx); xErr != nil {
		tx.Rollback()
		return xErr
	}

	if err := tx.Commit().Error; err != nil {
		return xError.NewError(ctx, xError.DatabaseError, "提交排序更新失败", false, err)
	}

	return nil
}

func (l *LinkGroupLogic) UpdateStatus(ctx *gin.Context, groupIDStr string, req *apiLinkGroup.GroupStatusRequest) *xError.Error {
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		return xError.NewError(ctx, xError.BadRequest, "无效的分组ID", false)
	}

	ok, xErr := l.repo.group.UpdateStatusByID(ctx, groupID, req.Status, nil)
	if xErr != nil {
		return xErr
	}
	if !ok {
		return xError.NewError(ctx, xError.NotFound, "友链分组不存在", false)
	}

	return nil
}

func (l *LinkGroupLogic) Delete(ctx *gin.Context, groupIDStr string, req *apiLinkGroup.GroupDeleteRequest) ([]entity.LinkFriend, *xError.Error) {
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		return nil, xError.NewError(ctx, xError.BadRequest, "无效的分组ID", false)
	}

	_, found, xErr := l.repo.group.GetByID(ctx, groupID, false, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "友链分组不存在", false)
	}

	linkCount, xErr := l.repo.link.CountByGroupID(ctx, groupID, nil)
	if xErr != nil {
		return nil, xErr
	}

	if linkCount > 0 && !req.Force {
		conflictLinks, xErr := l.repo.link.ListByGroupID(ctx, groupID, 10, nil)
		if xErr != nil {
			return nil, xErr
		}
		return conflictLinks, xError.NewError(ctx, xError.BadRequest, "分组下存在友链，无法删除", false)
	}

	tx := l.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if req.Force && linkCount > 0 {
		if xErr = l.repo.link.ClearGroupID(ctx, groupID, tx); xErr != nil {
			tx.Rollback()
			return nil, xErr
		}
	}

	ok, xErr := l.repo.group.DeleteByID(ctx, groupID, tx)
	if xErr != nil {
		tx.Rollback()
		return nil, xErr
	}
	if !ok {
		tx.Rollback()
		return nil, xError.NewError(ctx, xError.NotFound, "友链分组不存在", false)
	}

	if err = tx.Commit().Error; err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "提交删除操作失败", false, err)
	}

	return nil, nil
}

func (l *LinkGroupLogic) Get(ctx *gin.Context, groupIDStr string) (*entity.LinkGroup, *xError.Error) {
	groupID, err := strconv.ParseInt(groupIDStr, 10, 64)
	if err != nil {
		return nil, xError.NewError(ctx, xError.BadRequest, "无效的分组ID", false)
	}

	group, found, xErr := l.repo.group.GetByID(ctx, groupID, true, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "友链分组不存在", false)
	}

	return group, nil
}

func (l *LinkGroupLogic) GetList(ctx *gin.Context, req *apiLinkGroup.GroupListRequest) ([]entity.LinkGroup, *xError.Error) {
	return l.repo.group.List(ctx, req, nil)
}

func (l *LinkGroupLogic) GetPage(ctx *gin.Context, req *apiLinkGroup.GroupPageRequest) (*base.PaginationResponse[entity.LinkGroup], *xError.Error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 10
	}

	groups, total, xErr := l.repo.group.Page(ctx, req, nil)
	if xErr != nil {
		return nil, xErr
	}

	return base.NewPaginationResponse(groups, req.Page, req.PageSize, total), nil
}
