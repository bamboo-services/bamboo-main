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

package v1

import (
	"github.com/XiaoLFeng/bamboo-utils/bmodels"
	"github.com/gogf/gf/v2/frame/g"
	"go/types"
)

type AuthPasswordChangeReq struct {
	g.Meta               `path:"/auth/pass/change" method:"PUT" tags:"认证控制器" summary:"修改密码" dc:"修改密码接口，允许用户更新其登录密码"`
	OriginalPassword     string `json:"original_password" v:"required#请输入原密码" dc:"原密码，用户当前的登录密码"`
	NewPassword          string `json:"new_password" v:"required#请输入新密码" dc:"新密码，用户希望设置的新登录密码"`
	NewPasswordConfirm   string `json:"new_password_confirm" v:"required|same:new_password#请输入确认密码新密码和确认密码必须一致" dc:"确认密码，用户新密码的确认字段，必须与新密码一致"`
	NeedAllDeviceRefresh bool   `json:"refresh" default:"false" v:"required#请选择是否需要刷新所有设备的登录状态" dc:"是否需要刷新所有设备的登录状态，默认为false，表示只刷新当前设备"`
}

type AuthPasswordChangeRes struct {
	g.Meta `mime:"application/json; charset=utf-8" dc:"修改密码响应"`
	*bmodels.ResponseDTO[*types.Nil]
}
