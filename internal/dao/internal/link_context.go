// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// LinkContextDao is the data access object for the table xf_link_context.
type LinkContextDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  LinkContextColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// LinkContextColumns defines and stores column names for the table xf_link_context.
type LinkContextColumns struct {
	LinkUuid         string // 友链唯一标识符
	LinkName         string // 友链名称
	LinkUrl          string // 友链URL地址
	LinkAvatar       string // 友链头像URL
	LinkDesc         string // 友链描述
	LinkEmail        string // 友链联系邮箱
	LinkGroupUuid    string // 所属分组ID
	LinkColorUuid    string // 友链颜色ID
	LinkOrder        string // 友链排序
	LinkStatus       string // 友链状态（0: 待审核, 1: 已通过, 2: 已拒绝）
	LinkApplyRemark  string // 申请者备注
	LinkReviewRemark string // 审核备注
	LinkCreatedAt    string // 友链创建时间
	LinkUpdatedAt    string // 友链更新时间
}

// linkContextColumns holds the columns for the table xf_link_context.
var linkContextColumns = LinkContextColumns{
	LinkUuid:         "link_uuid",
	LinkName:         "link_name",
	LinkUrl:          "link_url",
	LinkAvatar:       "link_avatar",
	LinkDesc:         "link_desc",
	LinkEmail:        "link_email",
	LinkGroupUuid:    "link_group_uuid",
	LinkColorUuid:    "link_color_uuid",
	LinkOrder:        "link_order",
	LinkStatus:       "link_status",
	LinkApplyRemark:  "link_apply_remark",
	LinkReviewRemark: "link_review_remark",
	LinkCreatedAt:    "link_created_at",
	LinkUpdatedAt:    "link_updated_at",
}

// NewLinkContextDao creates and returns a new DAO object for table data access.
func NewLinkContextDao(handlers ...gdb.ModelHandler) *LinkContextDao {
	return &LinkContextDao{
		group:    "default",
		table:    "xf_link_context",
		columns:  linkContextColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *LinkContextDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *LinkContextDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *LinkContextDao) Columns() LinkContextColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *LinkContextDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *LinkContextDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *LinkContextDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
