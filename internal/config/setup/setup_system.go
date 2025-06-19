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
	"bamboo-main/internal/dao"
	"bamboo-main/internal/model/do"
	"bamboo-main/internal/model/entity"
	"context"
	"github.com/XiaoLFeng/bamboo-utils/blog"
	"github.com/XiaoLFeng/bamboo-utils/butil"
	"github.com/google/uuid"
)

// ProjectSetupSystemValueInit
//
// -
func (setup *projectSetup) ProjectSetupSystemValueInit() {
	blog.BambooNotice(setup.ctx, "ProjectSetupSystemValueInit", "初始化系统值")

	// 系统基础有关数据
	checkAndInsertData(setup.ctx, "system_name", "XiaoMain")
	checkAndInsertData(setup.ctx, "system_version", "v1.0.0-alpha")
	checkAndInsertData(setup.ctx, "system_description", "XiaoMain 是一个基于 Go 语言开发的开源项目，旨在提供一个高效、灵活的 Web 应用框架。")

	// 系统用户有关数据
	checkAndInsertData(setup.ctx, "user_username", "admin")
	checkAndInsertData(setup.ctx, "user_password", butil.PasswordEncode("admin"))
	checkAndInsertData(setup.ctx, "user_email", "admin@x-lf.cn")
	checkAndInsertData(setup.ctx, "user_nickname", "筱锋")
	checkAndInsertData(setup.ctx, "user_phone", "13388888880")
	checkAndInsertData(setup.ctx, "user_avatar_type", "local")
	checkAndInsertData(setup.ctx, "user_avatar_base64", "")
	checkAndInsertData(setup.ctx, "user_avatar_url", "https://i-cdn.akass.cn/2024/05/664870a814c0d.png!wp60")

}

// checkAndInsertData
//
// 检查数据是否存在，若不存在则插入新数据；
// 若数据存在则不进行任何操作；
// 最后会返回结果值供给全局变量使用「可选」。
func checkAndInsertData(ctx context.Context, name, value string) string {
	blog.BambooDebug(ctx, "checkAndInsertData", "检查并插入数据 %s: %s", name, value)

	var systemValue *entity.System
	daoErr := dao.System.Ctx(ctx).Where(&do.System{SystemName: name}).Scan(&systemValue)
	if daoErr != nil {
		blog.BambooError(ctx, "checkAndInsertData", "查询系统值失败: %v", daoErr)
		panic(daoErr)
	}
	if systemValue == nil {
		// 插入新数据
		var newSystemValue = &entity.System{
			SystemUuid:  uuid.New().String(),
			SystemName:  name,
			SystemValue: value,
		}
		_, daoErr := dao.System.Ctx(ctx).OmitEmpty().Insert(newSystemValue)
		if daoErr != nil {
			blog.BambooError(ctx, "checkAndInsertData", "插入系统值失败: %v", daoErr)
			panic(daoErr)
		}
		return value
	} else {
		return systemValue.SystemValue
	}
}
