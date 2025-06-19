/*
 * --------------------------------------------------------------------------------
 * Copyright (c) 2016-NOW(至今) 筱锋
 * Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
 * --------------------------------------------------------------------------------
 * 许可证声明：版权所有 (c) 2016-2025 筱锋。保留所有权利。
 * 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
 * https://opensource.org/licenses/MIT
 * --------------------------------------------------------------------------------
 */

package setup

import (
	"context"
	"github.com/XiaoLFeng/bamboo-utils/berror"
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// ProjectSetupForSQL
//
// 初始化项目的 SQL 数据库；
// 当对应的 SQL 数据库表不存在时，将会从模板中读取 SQL 文件内容并执行插入操作。
func (setup *projectSetup) ProjectSetupForSQL() {
	blog.BambooNotice(setup.ctx, "ProjectSetupForSQL", "正在进行项目的 SQL 数据库初始化设置...")

	for _, sqlName := range SQLNameList {
		record, daoErr := g.DB().Model("information_schema.tables").Where("table_name = ?", sqlName).One()
		if daoErr != nil {
			blog.BambooError(setup.ctx, "ProjectSetupForSQL", "查询 SQL 数据库表失败 %v:%v", sqlName, daoErr.Error())
			panic("查询 SQL 数据库表失败")
		}
		if record.IsEmpty() {
			errorCode := getSQLDataAndInsert(setup.ctx, sqlName)
			if errorCode != nil {
				blog.BambooError(setup.ctx, "ProjectSetupForSQL", "SQL 数据库初始化设置失败 %v:%v", sqlName, errorCode.Error())
			}
		}
	}
}

var SQLNameList = []string{
	"xf_system",
}

// getSQLDataAndInsert
//
// 根据名称获取 SQL 数据并插入到数据库中
func getSQLDataAndInsert(ctx context.Context, name string) *berror.ErrorCode {
	var dir = "template/sql/" + name + ".sql"

	getContent := gres.GetContent(dir)
	if len(getContent) == 0 {
		return &berror.ErrResourceNotFound
	}

	// 将内容按照 ';' 分割成多个 SQL 语句
	sqlStatements := gstr.SplitAndTrim(gconv.String(getContent), ";")
	for _, sql := range sqlStatements {
		if len(sql) == 0 {
			continue
		}
		// 执行 SQL 语句
		if _, err := g.DB().Exec(ctx, sql); err != nil {
			return berror.ErrorAddData(&berror.ErrDatabaseError, err.Error())
		}
	}
	return nil
}
