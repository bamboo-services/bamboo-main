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
	"bamboo-main/internal/model/response"
	"github.com/XiaoLFeng/bamboo-utils/bmodels"
	"github.com/gogf/gf/v2/frame/g"
)

type AuthLoginReq struct {
	g.Meta   `path:"/auth/login" method:"POST" tags:"认证控制器" summary:"用户登录" dc:"用户登录接口，使用用户名和密码进行身份验证，返回登录状态和用户信息"`
	User     string `json:"user" v:"required#请输入用户名、邮箱、手机号" dc:"用户名、邮箱或手机号"`
	Password string `json:"password" v:"required#请输入密码" dc:"密码"`
}

type AuthLoginRes struct {
	g.Meta `mime:"application/json; charset=utf-8" dc:"登录响应"`
	*bmodels.ResponseDTO[*response.AuthLoginResponse]
}
