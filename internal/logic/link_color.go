package logic

import (
	"context"
	"strconv"

	xError "github.com/bamboo-services/bamboo-base-go/common/error"
	xLog "github.com/bamboo-services/bamboo-base-go/common/log"
	xCtxUtil "github.com/bamboo-services/bamboo-base-go/common/utility/context"
	apiLinkColor "github.com/bamboo-services/bamboo-main/api/link"
	"github.com/bamboo-services/bamboo-main/internal/entity"
	"github.com/bamboo-services/bamboo-main/internal/models/base"
	"github.com/bamboo-services/bamboo-main/internal/repository"
	"github.com/gin-gonic/gin"
)

type linkColorRepo struct {
	color *repository.LinkColorRepo
	link  *repository.LinkRepo
}

type LinkColorLogic struct {
	logic
	repo linkColorRepo
}

func NewLinkColorLogic(ctx context.Context) *LinkColorLogic {
	db := xCtxUtil.MustGetDB(ctx)
	rdb := xCtxUtil.MustGetRDB(ctx)

	return &LinkColorLogic{
		logic: logic{
			db:  db,
			rdb: rdb,
			log: xLog.WithName(xLog.NamedLOGC, "LinkColorLogic"),
		},
		repo: linkColorRepo{
			color: repository.NewLinkColorRepo(db, rdb),
			link:  repository.NewLinkRepo(db, rdb),
		},
	}
}

func (l *LinkColorLogic) Add(ctx *gin.Context, req *apiLinkColor.ColorAddRequest) (*entity.LinkColor, *xError.Error) {
	if req.ColorType == 0 {
		if req.PrimaryColor == nil || req.SubColor == nil || req.HoverColor == nil {
			return nil, xError.NewError(ctx, xError.BadRequest, "普通颜色类型需要设置主颜色、副颜色和悬停颜色", false)
		}
		if *req.PrimaryColor == "" || *req.SubColor == "" || *req.HoverColor == "" {
			return nil, xError.NewError(ctx, xError.BadRequest, "普通颜色类型的颜色值不能为空", false)
		}
	}

	maxSort, xErr := l.repo.color.GetMaxSortOrder(ctx, nil)
	if xErr != nil {
		return nil, xErr
	}

	color := &entity.LinkColor{
		Name:      req.ColorName,
		Type:      req.ColorType,
		SortOrder: maxSort + 1,
		Status:    true,
	}
	if req.ColorType == 0 {
		color.PrimaryColor = req.PrimaryColor
		color.SubColor = req.SubColor
		color.HoverColor = req.HoverColor
	}

	_, xErr = l.repo.color.Create(ctx, color, nil)
	if xErr != nil {
		return nil, xErr
	}

	reloaded, found, xErr := l.repo.color.GetByID(ctx, color.ID, true, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "友链颜色不存在", false)
	}

	return reloaded, nil
}

func (l *LinkColorLogic) Update(ctx *gin.Context, colorIDStr string, req *apiLinkColor.ColorUpdateRequest) (*entity.LinkColor, *xError.Error) {
	colorID, err := strconv.ParseInt(colorIDStr, 10, 64)
	if err != nil {
		return nil, xError.NewError(ctx, xError.BadRequest, "无效的颜色ID", false)
	}

	color, found, xErr := l.repo.color.GetByID(ctx, colorID, false, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "友链颜色不存在", false)
	}

	if req.ColorName != nil {
		color.Name = *req.ColorName
	}
	if req.ColorType != nil {
		color.Type = *req.ColorType
	}
	if req.ColorOrder != nil {
		color.SortOrder = *req.ColorOrder
	}

	if req.PrimaryColor != nil {
		if *req.PrimaryColor == "" {
			color.PrimaryColor = nil
		} else {
			color.PrimaryColor = req.PrimaryColor
		}
	}
	if req.SubColor != nil {
		if *req.SubColor == "" {
			color.SubColor = nil
		} else {
			color.SubColor = req.SubColor
		}
	}
	if req.HoverColor != nil {
		if *req.HoverColor == "" {
			color.HoverColor = nil
		} else {
			color.HoverColor = req.HoverColor
		}
	}

	if color.Type == 0 {
		if color.PrimaryColor == nil || color.SubColor == nil || color.HoverColor == nil {
			return nil, xError.NewError(ctx, xError.BadRequest, "普通颜色类型需要设置主颜色、副颜色和悬停颜色", false)
		}
	}

	_, xErr = l.repo.color.Save(ctx, color, nil)
	if xErr != nil {
		return nil, xErr
	}

	reloaded, found, xErr := l.repo.color.GetByID(ctx, colorID, true, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "友链颜色不存在", false)
	}

	return reloaded, nil
}

func (l *LinkColorLogic) UpdateSort(ctx *gin.Context, req *apiLinkColor.ColorSortRequest) *xError.Error {
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

	if xErr := l.repo.color.UpdateSortByIDs(ctx, req.ColorIDs, startSort, tx); xErr != nil {
		tx.Rollback()
		return xErr
	}

	if err := tx.Commit().Error; err != nil {
		return xError.NewError(ctx, xError.DatabaseError, "提交排序更新失败", false, err)
	}

	return nil
}

func (l *LinkColorLogic) UpdateStatus(ctx *gin.Context, colorIDStr string, req *apiLinkColor.ColorStatusRequest) *xError.Error {
	colorID, err := strconv.ParseInt(colorIDStr, 10, 64)
	if err != nil {
		return xError.NewError(ctx, xError.BadRequest, "无效的颜色ID", false)
	}

	ok, xErr := l.repo.color.UpdateStatusByID(ctx, colorID, req.Status, nil)
	if xErr != nil {
		return xErr
	}
	if !ok {
		return xError.NewError(ctx, xError.NotFound, "友链颜色不存在", false)
	}

	return nil
}

func (l *LinkColorLogic) Delete(ctx *gin.Context, colorIDStr string, req *apiLinkColor.ColorDeleteRequest) ([]entity.LinkFriend, *xError.Error) {
	colorID, err := strconv.ParseInt(colorIDStr, 10, 64)
	if err != nil {
		return nil, xError.NewError(ctx, xError.BadRequest, "无效的颜色ID", false)
	}

	_, found, xErr := l.repo.color.GetByID(ctx, colorID, false, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "友链颜色不存在", false)
	}

	linkCount, xErr := l.repo.link.CountByColorID(ctx, colorID, nil)
	if xErr != nil {
		return nil, xErr
	}

	if linkCount > 0 && !req.Force {
		conflictLinks, xErr := l.repo.link.ListByColorID(ctx, colorID, 10, nil)
		if xErr != nil {
			return nil, xErr
		}
		return conflictLinks, xError.NewError(ctx, xError.BadRequest, "颜色下存在友链，无法删除", false)
	}

	tx := l.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if req.Force && linkCount > 0 {
		if xErr = l.repo.link.ClearColorID(ctx, colorID, tx); xErr != nil {
			tx.Rollback()
			return nil, xErr
		}
	}

	ok, xErr := l.repo.color.DeleteByID(ctx, colorID, tx)
	if xErr != nil {
		tx.Rollback()
		return nil, xErr
	}
	if !ok {
		tx.Rollback()
		return nil, xError.NewError(ctx, xError.NotFound, "友链颜色不存在", false)
	}

	if err = tx.Commit().Error; err != nil {
		return nil, xError.NewError(ctx, xError.DatabaseError, "提交删除操作失败", false, err)
	}

	return nil, nil
}

func (l *LinkColorLogic) Get(ctx *gin.Context, colorIDStr string) (*entity.LinkColor, *xError.Error) {
	colorID, err := strconv.ParseInt(colorIDStr, 10, 64)
	if err != nil {
		return nil, xError.NewError(ctx, xError.BadRequest, "无效的颜色ID", false)
	}

	color, found, xErr := l.repo.color.GetByID(ctx, colorID, true, nil)
	if xErr != nil {
		return nil, xErr
	}
	if !found {
		return nil, xError.NewError(ctx, xError.NotFound, "友链颜色不存在", false)
	}

	return color, nil
}

func (l *LinkColorLogic) GetList(ctx *gin.Context, req *apiLinkColor.ColorListRequest) ([]entity.LinkColor, *xError.Error) {
	return l.repo.color.List(ctx, req, nil)
}

func (l *LinkColorLogic) GetPage(ctx *gin.Context, req *apiLinkColor.ColorPageRequest) (*base.PaginationResponse[entity.LinkColor], *xError.Error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 10
	}

	colors, total, xErr := l.repo.color.Page(ctx, req, nil)
	if xErr != nil {
		return nil, xErr
	}

	return base.NewPaginationResponse(colors, req.Page, req.PageSize, total), nil
}
