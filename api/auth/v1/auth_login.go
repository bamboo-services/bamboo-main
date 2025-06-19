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

type AuthLoginReq struct {
	g.Meta   `path:"/auth/login" method:"POST" tags:"认证控制器" summary:"用户登录" dc:"用户登录接口，使用用户名和密码进行身份验证，返回登录状态和用户信息"`
	Username string `json:"username" v:"required|regex:^[0-9A-Za-z-\\_]{6,32}#请输入用户名|用户名格式不正确，必须是6-32位的字母或数字" dc:"用户名，必填，6-32位字母或数字"`
	Password string `json:"password" v:"required#请输入密码" dc:"密码，必填"`
}

type AuthLoginRes struct {
	g.Meta `mime:"application/json; charset=utf-8" dc:"登录响应"`
	bmodels.ResponseDTO[*types.Nil]
}
