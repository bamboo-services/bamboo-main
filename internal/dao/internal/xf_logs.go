// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// XfLogsDao is the data access object for table xf_logs.
type XfLogsDao struct {
	table   string        // table is the underlying table name of the DAO.
	group   string        // group is the database configuration group name of current DAO.
	columns XfLogsColumns // columns contains all the column names of Table for convenient usage.
}

// XfLogsColumns defines and stores column names for table xf_logs.
type XfLogsColumns struct {
	Id        string // 主键
	Type      string // 日志类型
	Log       string // 日志内容
	CreatedAt string // 日志时间
}

// xfLogsColumns holds the columns for table xf_logs.
var xfLogsColumns = XfLogsColumns{
	Id:        "id",
	Type:      "type",
	Log:       "log",
	CreatedAt: "created_at",
}

// NewXfLogsDao creates and returns a new DAO object for table data access.
func NewXfLogsDao() *XfLogsDao {
	return &XfLogsDao{
		group:   "default",
		table:   "xf_logs",
		columns: xfLogsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *XfLogsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *XfLogsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *XfLogsDao) Columns() XfLogsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *XfLogsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *XfLogsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *XfLogsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
