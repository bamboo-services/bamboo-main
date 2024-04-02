// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// XfTokenDao is the data access object for table xf_token.
type XfTokenDao struct {
	table   string         // table is the underlying table name of the DAO.
	group   string         // group is the database configuration group name of current DAO.
	columns XfTokenColumns // columns contains all the column names of Table for convenient usage.
}

// XfTokenColumns defines and stores column names for table xf_token.
type XfTokenColumns struct {
	Id        string // 主键
	UserUuid  string // 用户 UUID
	UserToken string // 用户 TOKEN
	CreatedAt string // 创建时间
	ExpiredAt string // 过期时间
}

// xfTokenColumns holds the columns for table xf_token.
var xfTokenColumns = XfTokenColumns{
	Id:        "id",
	UserUuid:  "user_uuid",
	UserToken: "user_token",
	CreatedAt: "created_at",
	ExpiredAt: "expired_at",
}

// NewXfTokenDao creates and returns a new DAO object for table data access.
func NewXfTokenDao() *XfTokenDao {
	return &XfTokenDao{
		group:   "default",
		table:   "xf_token",
		columns: xfTokenColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *XfTokenDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *XfTokenDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *XfTokenDao) Columns() XfTokenColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *XfTokenDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *XfTokenDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *XfTokenDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
