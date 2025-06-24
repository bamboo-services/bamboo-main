// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// LinkGroupDao is the data access object for the table xf_link_group.
type LinkGroupDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  LinkGroupColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// LinkGroupColumns defines and stores column names for the table xf_link_group.
type LinkGroupColumns struct {
	GroupUuid      string // 分组唯一标识符
	GroupName      string // 分组名称
	GroupDesc      string // 分组描述
	GroupOrder     string // 分组排序
	GroupCreatedAt string // 分组创建时间
	GroupUpdatedAt string // 分组更新时间
}

// linkGroupColumns holds the columns for the table xf_link_group.
var linkGroupColumns = LinkGroupColumns{
	GroupUuid:      "group_uuid",
	GroupName:      "group_name",
	GroupDesc:      "group_desc",
	GroupOrder:     "group_order",
	GroupCreatedAt: "group_created_at",
	GroupUpdatedAt: "group_updated_at",
}

// NewLinkGroupDao creates and returns a new DAO object for table data access.
func NewLinkGroupDao(handlers ...gdb.ModelHandler) *LinkGroupDao {
	return &LinkGroupDao{
		group:    "default",
		table:    "xf_link_group",
		columns:  linkGroupColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *LinkGroupDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *LinkGroupDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *LinkGroupDao) Columns() LinkGroupColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *LinkGroupDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *LinkGroupDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *LinkGroupDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
