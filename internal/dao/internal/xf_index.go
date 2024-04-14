// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// XfIndexDao is the data access object for table xf_index.sql.
type XfIndexDao struct {
	table   string         // table is the underlying table name of the DAO.
	group   string         // group is the database configuration group name of current DAO.
	columns XfIndexColumns // columns contains all the column names of Table for convenient usage.
}

// XfIndexColumns defines and stores column names for table xf_index.sql.
type XfIndexColumns struct {
	Id        string // 主键
	Key       string // 键
	Value     string // 值
	CreatedAt string // 创建时间
	UpdatedAt string // 修改时间
}

// xfIndexColumns holds the columns for table xf_index.sql.
var xfIndexColumns = XfIndexColumns{
	Id:        "id",
	Key:       "key",
	Value:     "value",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
}

// NewXfIndexDao creates and returns a new DAO object for table data access.
func NewXfIndexDao() *XfIndexDao {
	return &XfIndexDao{
		group:   "default",
		table:   "xf_index.sql",
		columns: xfIndexColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *XfIndexDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *XfIndexDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *XfIndexDao) Columns() XfIndexColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *XfIndexDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *XfIndexDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *XfIndexDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
