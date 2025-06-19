// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// LogsDao is the data access object for the table xf_logs.
type LogsDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  LogsColumns        // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// LogsColumns defines and stores column names for the table xf_logs.
type LogsColumns struct {
	LogUuid      string // 日志唯一标识符
	LogType      string // 日志类型
	LogContent   string // 日志内容
	LogCreatedAt string // 日志创建时间
}

// logsColumns holds the columns for the table xf_logs.
var logsColumns = LogsColumns{
	LogUuid:      "log_uuid",
	LogType:      "log_type",
	LogContent:   "log_content",
	LogCreatedAt: "log_created_at",
}

// NewLogsDao creates and returns a new DAO object for table data access.
func NewLogsDao(handlers ...gdb.ModelHandler) *LogsDao {
	return &LogsDao{
		group:    "default",
		table:    "xf_logs",
		columns:  logsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *LogsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *LogsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *LogsDao) Columns() LogsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *LogsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *LogsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *LogsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
