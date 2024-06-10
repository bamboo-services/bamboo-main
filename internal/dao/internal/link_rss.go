/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋(https://www.x-lf.com)
 *
 * 本文件包含 XiaoMain 的源代码，该项目的所有源代码均遵循MIT开源许可证协议。
 * --------------------------------------------------------------------------------
 * 许可证声明：
 *
 * 版权所有 (c) 2016-2024 筱锋。保留所有权利。
 *
 * 本软件是“按原样”提供的，没有任何形式的明示或暗示的保证，包括但不限于
 * 对适销性、特定用途的适用性和非侵权性的暗示保证。在任何情况下，
 * 作者或版权持有人均不承担因软件或软件的使用或其他交易而产生的、
 * 由此引起的或以任何方式与此软件有关的任何索赔、损害或其他责任。
 *
 * 使用本软件即表示您了解此声明并同意其条款。
 *
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 * 免责声明：
 *
 * 使用本软件的风险由用户自担。作者或版权持有人在法律允许的最大范围内，
 * 对因使用本软件内容而导致的任何直接或间接的损失不承担任何责任。
 * --------------------------------------------------------------------------------
 */

// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// LinkRssDao is the data access object for table xf_link_rss.
type LinkRssDao struct {
	table   string         // table is the underlying table name of the DAO.
	group   string         // group is the database configuration group name of current DAO.
	columns LinkRssColumns // columns contains all the column names of Table for convenient usage.
}

// LinkRssColumns defines and stores column names for table xf_link_rss.
type LinkRssColumns struct {
	LinkId    string // 链接 id
	RssJson   string // Rss 内容解析
	CheckTime string // 检查时间
}

// linkRssColumns holds the columns for table xf_link_rss.
var linkRssColumns = LinkRssColumns{
	LinkId:    "link_id",
	RssJson:   "rss_json",
	CheckTime: "check_time",
}

// NewLinkRssDao creates and returns a new DAO object for table data access.
func NewLinkRssDao() *LinkRssDao {
	return &LinkRssDao{
		group:   "default",
		table:   "xf_link_rss",
		columns: linkRssColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *LinkRssDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *LinkRssDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *LinkRssDao) Columns() LinkRssColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *LinkRssDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *LinkRssDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *LinkRssDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
