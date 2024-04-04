// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// XfLinkListDao is the data access object for table xf_link_list.
type XfLinkListDao struct {
	table   string            // table is the underlying table name of the DAO.
	group   string            // group is the database configuration group name of current DAO.
	columns XfLinkListColumns // columns contains all the column names of Table for convenient usage.
}

// XfLinkListColumns defines and stores column names for table xf_link_list.
type XfLinkListColumns struct {
	Id              string // 主键
	WebmasterEmail  string // 站长邮箱
	ServiceProvider string // 服务提供商
	SiteName        string // 站点名字
	SiteUrl         string // 站点地址
	SiteLogo        string // 站点 logo
	SiteDescription string // 站点描述
	SiteRssUrl      string // 站点订阅地址
	HasAdv          string // 是否有广告
	DesiredLocation string // 理想位置
	Location        string // 所在位置
	DesiredColor    string // 理想颜色
	Color           string // 颜色
	WebmasterRemark string // 站长留言
	Remark          string // 我的留言
	Status          string // 0 待审核，1 通过，-1 审核拒绝
	CreatedAt       string // 创建时间
	UpdatedAt       string // 修改时间
	DeletedAt       string // 删除时间
}

// xfLinkListColumns holds the columns for table xf_link_list.
var xfLinkListColumns = XfLinkListColumns{
	Id:              "id",
	WebmasterEmail:  "webmaster_email",
	ServiceProvider: "service_provider",
	SiteName:        "site_name",
	SiteUrl:         "site_url",
	SiteLogo:        "site_logo",
	SiteDescription: "site_description",
	SiteRssUrl:      "site_rss_url",
	HasAdv:          "has_adv",
	DesiredLocation: "desired_location",
	Location:        "location",
	DesiredColor:    "desired_color",
	Color:           "color",
	WebmasterRemark: "webmaster_remark",
	Remark:          "remark",
	Status:          "status",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
	DeletedAt:       "deleted_at",
}

// NewXfLinkListDao creates and returns a new DAO object for table data access.
func NewXfLinkListDao() *XfLinkListDao {
	return &XfLinkListDao{
		group:   "default",
		table:   "xf_link_list",
		columns: xfLinkListColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *XfLinkListDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *XfLinkListDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *XfLinkListDao) Columns() XfLinkListColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *XfLinkListDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *XfLinkListDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *XfLinkListDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
