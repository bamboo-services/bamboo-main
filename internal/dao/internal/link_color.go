// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// LinkColorDao is the data access object for the table xf_link_color.
type LinkColorDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  LinkColorColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// LinkColorColumns defines and stores column names for the table xf_link_color.
type LinkColorColumns struct {
	ColorUuid      string // 颜色唯一标识符
	ColorName      string // 颜色名称
	ColorValue     string // 颜色值（如HEX值：#FFFFFF）
	ColorDesc      string // 颜色描述
	ColorCreatedAt string // 颜色创建时间
	ColorUpdatedAt string // 颜色更新时间
}

// linkColorColumns holds the columns for the table xf_link_color.
var linkColorColumns = LinkColorColumns{
	ColorUuid:      "color_uuid",
	ColorName:      "color_name",
	ColorValue:     "color_value",
	ColorDesc:      "color_desc",
	ColorCreatedAt: "color_created_at",
	ColorUpdatedAt: "color_updated_at",
}

// NewLinkColorDao creates and returns a new DAO object for table data access.
func NewLinkColorDao(handlers ...gdb.ModelHandler) *LinkColorDao {
	return &LinkColorDao{
		group:    "default",
		table:    "xf_link_color",
		columns:  linkColorColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *LinkColorDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *LinkColorDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *LinkColorDao) Columns() LinkColorColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *LinkColorDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *LinkColorDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *LinkColorDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
